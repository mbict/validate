package validate

import (
	"github.com/mbict/go-errors"
	"github.com/mbict/go-tags"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

// Validator interface
type Validator interface {
	SetTag(tag string)
	WithTag(tag string) Validator
	SetValidationFunc(name string, vf ValidatorFunc) error
	SetNameResolver(resolver NameResolverFunc)
	Validate(v interface{}) error
	Valid(val interface{}, tags string) error
}

// validatorTag represents one of the validatorTag items
type validatorTag struct {
	tags.Param               // name of the validator and the arguments to send to the validator func
	Fn         ValidatorFunc // validation function to call
}

// ValidateInterface describes the interface a structure can embed to enable custom validation of the structure
type ValidateInterface interface {
	Validate() error
}

// ValidatorFunc is a function that receives the value of a
// field and the parameters used for the respective validation validatorTag.
type ValidatorFunc func(v interface{}, params []string) error

// validator implements the Validator interface
type validator struct {
	tagName         string                   // structure validatorTag name being used (`validate`)
	validationFuncs map[string]ValidatorFunc // validator functions map indexed by name
	structRules     StructRules              // structure rules cache
	mu              sync.RWMutex             // rw mutex for structure rules cache
	nameResolver    NameResolverFunc         // func to extract the name to use for field error
}

type NameResolverFunc func(reflect.StructField) string

var DefaultNameResolver = func(field reflect.StructField) string {
	return field.Name
}

var JsonNameResolver = func(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	props := strings.SplitN(tag, ",", 2)
	if props[0] == "" {
		return field.Name
	}
	return props[0]
}

// Helper validator so users can use the
// functions directly from the package
var defaultValidator = NewValidator()

type Option func(Validator)

func NameResolverOption(resolver NameResolverFunc) Option {
	return func(v Validator) {
		v.SetNameResolver(resolver)
	}
}

func TagOption(tag string) Option {
	return func(v Validator) {
		v.SetTag(tag)
	}
}

func ValidatorOption(name string, validatorFunc ValidatorFunc) Option {
	return func(v Validator) {
		v.SetValidationFunc(name, validatorFunc)
	}
}

// NewValidator creates a new Validator
func NewValidator(options ...Option) Validator {
	v := &validator{
		tagName: "validate",
		validationFuncs: map[string]ValidatorFunc{
			"omitempty":      omitempty,
			"required":       required,
			"not_empty":      notEmpty,
			"len":            length,
			"min":            min,
			"max":            max,
			"between":        between,
			"around":         around,
			"in":             include,
			"exclude":        exclude,
			"regexp":         regex,
			"url":            url,
			"email":          email,
			"numeric":        numeric,
			"number":         number,
			"identifier":     identifier,
			"alpha_dash":     alphaDash,
			"alpha_dash_dot": alphaDashDot,
			"alpha":          alpha,
			"alphanumeric":   alphaNumeric,
			"uuid":           uuid,
			"uuid3":          uuid3,
			"uuid4":          uuid4,
			"uuid5":          uuid5,
			"base64":         base64,
		},
		structRules:  make(StructRules),
		nameResolver: DefaultNameResolver,
	}

	for _, option := range options {
		option(v)
	}

	return v
}

// SetTag allows you to change the validatorTag name used in structs for the default validator
func SetTag(tag string) {
	defaultValidator.SetTag(tag)
}

// WithTag creates a new Validator with the new validatorTag name. It will leave
// the defaultValidator untouched
func WithTag(tag string) Validator {
	return defaultValidator.WithTag(tag)
}

// SetValidationFunc sets the function to be used for a given validation constraint.
// Calling this function with nil validatorFunction (vf) is the same as removing
// the constraint function from the list. The function will be added to the default
// validator
func SetValidationFunc(name string, vf ValidatorFunc) error {
	return defaultValidator.SetValidationFunc(name, vf)
}

// Validate validates the fields of a struct based  on 'validator' tags and returns
// errors found indexed by the field name.
func Validate(v interface{}) error {
	return defaultValidator.Validate(v)
}

// Valid validates a value based on the provided tags and returns errors found or nil.
func Valid(val interface{}, tags string) error {
	return defaultValidator.Valid(val, tags)
}

// SetNameResolver allows you to change the way field names are resolved
func (mv *validator) SetNameResolver(resolver NameResolverFunc) {
	mv.nameResolver = resolver
	mv.resetCache()
}

// SetTag allows you to change the validatorTag name used in structs
func (mv *validator) SetTag(tag string) {
	mv.tagName = tag
	mv.resetCache()
}

// WithTag creates a new Validator based on the current validator with the new validatorTag name.
func (mv *validator) WithTag(tag string) Validator {
	v := mv.copy()
	v.SetTag(tag)
	return v
}

// copy creates a duplicate of the current validator and returns the new instance
func (mv *validator) copy() Validator {
	return &validator{
		tagName:         mv.tagName,
		validationFuncs: mv.validationFuncs,
		structRules:     mv.structRules,
		nameResolver:    mv.nameResolver,
	}
}

// SetValidationFunc sets the function to be used for a given validation constraint.
// Calling this function with nil validatorFunction (vf) is the same as removing
// the constraint function from the list.
func (mv *validator) SetValidationFunc(name string, vf ValidatorFunc) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if vf == nil {
		delete(mv.validationFuncs, name)
		return nil
	}
	mv.validationFuncs[name] = vf
	return nil
}

// Validate validates the fields of a struct based on 'validator' tags and returns
// errors found indexed by the field name.
// The returned error is of the type errors.ErrorHash.
func (mv *validator) Validate(v interface{}) error {
	sv := reflect.ValueOf(v)

	//nil pointer not type found
	if sv.Kind() == reflect.Invalid {
		return ErrUnsupported
	}

	if sv.Kind() == reflect.Ptr && !sv.IsNil() {
		return mv.Validate(sv.Elem().Interface())
	}

	if sv.Kind() != reflect.Struct {
		return ErrUnsupported
	}

	mv.mu.RLock()
	rules, ok := mv.structRules[sv.Type()]
	mv.mu.RUnlock()
	if !ok {
		mv.mu.Lock()
		defer mv.mu.Unlock()
		r, err := mv.parseStruct(sv.Type())
		if err != nil {
			return err
		}
		mv.structRules[sv.Type()] = r
		rules = r
	}

	if errs := rules.Validate(sv); len(errs) > 0 {
		return errs
	}
	return nil
}

// Valid validates a value based on the provided tags and returns errors found or nil.
func (mv *validator) Valid(val interface{}, tags string) error {
	if tags == "-" {
		return nil
	}
	v := reflect.ValueOf(val)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		return mv.Valid(v.Elem().Interface(), tags)
	}

	if v.Kind() == reflect.Invalid {
		return mv.validateVar(nil, tags)
	}
	return mv.validateVar(val, tags)
}

func (mv *validator) resetCache() {
	mv.mu.Lock()
	defer mv.mu.Unlock()
	mv.structRules = make(StructRules, 0)
}

// validateVar validates a single variable
func (mv *validator) validateVar(v interface{}, tag string) error {
	tags, err := mv.parseTags(tag)
	if err != nil {
		// unknown validatorTag found.
		return err
	}
	var errs ErrorList
	for _, t := range tags {
		if err := t.Fn(v, t.Args); err != nil {
			if err == errOmitEmpty {
				return nil
			}
			errs = append(errs, err)
		}
	}
	return errs
}

// parseTags parses all individual tags found within a struct validatorTag and
// resolve the validator function
func (mv *validator) parseTags(t string) ([]validatorTag, error) {
	params, err := tags.Parse(t)
	if err != nil {
		return nil, ErrSyntax
	}

	tags := make([]validatorTag, 0, len(params))
	for _, param := range params {
		validatorFunc, found := mv.validationFuncs[param.Name]
		if !found {
			return nil, ErrUnknownTag
		}

		tags = append(tags, validatorTag{
			Param: param,
			Fn:    validatorFunc,
		})
	}
	return tags, nil
}

// parseStruct will extract all the validation rules from the given structure
func (mv *validator) parseStruct(t reflect.Type) (Rules, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, ErrUnsupported
	}

	rules := Rules{}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		tag := sf.Tag.Get(mv.tagName)

		if tag == "-" {
			continue
		}

		fieldName := mv.nameResolver(sf)
		if !unicode.IsUpper(rune(fieldName[0])) {
			continue
		}

		rule := Rule{
			Name:       fieldName,
			FieldIndex: i,
			IsSlice:    false,
			IsStruct:   false,
		}

		if tag != "" {
			//extract the validator properties
			validatorTags, err := mv.parseTags(tag)
			if err != nil {
				// unknown validatorTag found.
				return nil, err
			}

			//Add validatorTags to rules
			rule.Validators = validatorTags

		}

		st := sf.Type
		if st.Kind() == reflect.Slice {
			rule.IsSlice = true
			st = st.Elem()
		}

		if st.Kind() == reflect.Ptr {
			st = st.Elem()
		}

		if st.Kind() == reflect.Struct {
			subset, ok := mv.structRules[st]
			if !ok {
				var err error
				subset, err = mv.parseStruct(st)
				if err != nil {
					return nil, err
				}
				mv.structRules[st] = subset
			}
			rule.IsStruct = true
			rule.Subset = subset
		}

		rules = append(rules, rule)
	}

	return rules, nil
}
