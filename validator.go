package validate

import (
	"fmt"
	"github.com/mbict/go-errors"
	"github.com/mbict/go-tags"
	"reflect"
	"unicode"
)

// Validator interface
type Validator interface {
	SetTag(tag string)
	WithTag(tag string) Validator
	SetValidationFunc(name string, vf ValidationFunc) error
	Validate(v interface{}) error
	Valid(val interface{}, tags string) error
}

// tag represents one of the tag items
type tag struct {
	tags.Param                // name of the validator and the arguments to send to the validator func
	Fn         ValidationFunc // validation function to call
}

// ValidateInterface describes the interface a structure can embed to enable custom validation of the structure
type ValidateInterface interface {
	Validate() errors.ErrorHash
}

// ValidationFunc is a function that receives the value of a
// field and the parameters used for the respective validation tag.
type ValidationFunc func(v interface{}, params []string) error

// validator implements the Validator interface
type validator struct {
	tagName         string                    // structure tag name being used (`validate`)
	validationFuncs map[string]ValidationFunc // validator functions map indexed by name
}

// Helper validator so users can use the
// functions directly from the package
var defaultValidator = NewValidator()

// NewValidator creates a new Validator
func NewValidator() Validator {
	return &validator{
		tagName: "validate",
		validationFuncs: map[string]ValidationFunc{
			"required":       required,
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
	}
}

// SetTag allows you to change the tag name used in structs for the default validator
func SetTag(tag string) {
	defaultValidator.SetTag(tag)
}

// WithTag creates a new Validator with the new tag name. It will leave
// the defaultValidator untouched
func WithTag(tag string) Validator {
	return defaultValidator.WithTag(tag)
}

// SetValidationFunc sets the function to be used for a given validation constraint.
// Calling this function with nil validatorFunction (vf) is the same as removing
// the constraint function from the list. The function will be added to the default
// validator
func SetValidationFunc(name string, vf ValidationFunc) error {
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

// SetTag allows you to change the tag name used in structs
func (mv *validator) SetTag(tag string) {
	mv.tagName = tag
}

// WithTag creates a new Validator based on the current validator with the new tag name.
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
	}
}

// SetValidationFunc sets the function to be used for a given validation constraint.
// Calling this function with nil validatorFunction (vf) is the same as removing
// the constraint function from the list.
func (mv *validator) SetValidationFunc(name string, vf ValidationFunc) error {
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
	st := reflect.TypeOf(v)
	if sv.Kind() == reflect.Ptr && !sv.IsNil() {
		return mv.Validate(sv.Elem().Interface())
	}

	if sv.Kind() != reflect.Struct {
		return ErrUnsupported
	}

	nfields := sv.NumField()
	m := make(errors.ErrorHash)
	for i := 0; i < nfields; i++ {
		f := sv.Field(i)

		// deal with pointers
		for f.Kind() == reflect.Ptr && !f.IsNil() {
			f = f.Elem()
		}
		tag := st.Field(i).Tag.Get(mv.tagName)

		if tag == "-" {
			continue
		}

		fname := st.Field(i).Name
		if !unicode.IsUpper(rune(fname[0])) {
			continue
		}

		var errs errors.Errors
		if tag != "" {
			err := mv.Valid(f.Interface(), tag)
			if e, ok := err.(errors.Errors); ok {
				errs = e
			} else {
				if err != nil {
					errs = errors.Errors{err}
				}
			}
		}

		if f.Kind() == reflect.Slice {
			t := f.Type().Elem()
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}

			if t.Kind() == reflect.Struct {
				for i := 0; i < f.Len(); i++ {
					e := mv.Validate(f.Index(i).Interface())
					if e, ok := e.(errors.ErrorHash); ok && len(e) > 0 {
						for j, k := range e {
							field := fmt.Sprintf("%s.%d.%s", fname, i, j)
							m[field] = k
						}
					}
				}
			}
		} else if f.Kind() == reflect.Struct {

			e := mv.Validate(f.Interface())
			if e, ok := e.(errors.ErrorHash); ok && len(e) > 0 {
				for j, k := range e {
					m[fname+"."+j] = k
				}
			}
		}

		if len(errs) > 0 {
			m[st.Field(i).Name] = errs
		}
	}

	//structure custom validator function
	i := sv.Interface()
	if validateFunc, ok := i.(ValidateInterface); ok {
		em := validateFunc.Validate()
		if em != nil && len(em) > 0 {
			for f, e := range em {
				m[f] = append(m[f], e...)
			}
		}
	}

	if len(m) > 0 {
		return m
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

	var err error
	switch v.Kind() {
	case reflect.Invalid:
		err = mv.validateVar(nil, tags)
	default:
		err = mv.validateVar(val, tags)
	}
	return err
}

// validateVar validates a single variable
func (mv *validator) validateVar(v interface{}, tag string) error {
	tags, err := mv.parseTags(tag)
	if err != nil {
		// unknown tag found.
		return err
	}
	errs := make(errors.Errors, 0, len(tags))
	for _, t := range tags {
		if err := t.Fn(v, t.Args); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// parseTags parses all individual tags found within a struct tag and
// resolve the validator function
func (mv *validator) parseTags(t string) ([]tag, error) {
	params, err := tags.Parse(t)
	if err != nil {
		return nil, ErrSyntax
	}

	tags := make([]tag, 0, len(params))
	for _, param := range params {
		validatorFunc, found := mv.validationFuncs[param.Name]
		if !found {
			return nil, ErrUnknownTag
		}

		tags = append(tags, tag{
			Param: param,
			Fn:    validatorFunc,
		})
	}
	return tags, nil
}
