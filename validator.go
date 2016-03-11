package validator

import (
	"errors"
	"reflect"
	"strings"
	"unicode"
)

// tag represents one of the tag items
type tag struct {
	Name  string         // name of the tag
	Fn    ValidationFunc // validation function to call
	Param string         // parameter to send to the validation function
}

// ValidationFunc is a function that receives the value of a
// field and a parameter used for the respective validation tag.
type ValidationFunc func(v interface{}, param string) error

// Validator implements a validator
type Validator struct {
	// Tag name being used.
	tagName string

	// validationFuncs is a map of ValidationFuncs indexed
	// by their name.
	validationFuncs map[string]ValidationFunc
}

// Helper validator so users can use the
// functions directly from the package
var defaultValidator = NewValidator()

// NewValidator creates a new Validator
func NewValidator() *Validator {
	return &Validator{
		tagName: "validate",
		validationFuncs: map[string]ValidationFunc{
			"nonzero": nonzero,
			"len":     length,
			"min":     min,
			"max":     max,
			"regexp":  regex,
		},
	}
}

// SetTag allows you to change the tag name used in structs
func SetTag(tag string) {
	defaultValidator.SetTag(tag)
}

// WithTag creates a new Validator with the new tag name. It is
// useful to chain-call with Validate so we don't change the tag
// name permanently: validator.WithTag("foo").Validate(t)
func WithTag(tag string) *Validator {
	return defaultValidator.WithTag(tag)
}

// SetValidationFunc sets the function to be used for a given
// validation constraint. Calling this function with nil vf
// is the same as removing the constraint function from the list.
func SetValidationFunc(name string, vf ValidationFunc) error {
	return defaultValidator.SetValidationFunc(name, vf)
}

// Validate validates the fields of a struct based
// on 'validator' tags and returns errors found indexed
// by the field name.
func Validate(v interface{}) error {
	return defaultValidator.Validate(v)
}

// Valid validates a value based on the provided
// tags and returns errors found or nil.
func Valid(val interface{}, tags string) error {
	return defaultValidator.Valid(val, tags)
}

// SetTag allows you to change the tag name used in structs
func (mv *Validator) SetTag(tag string) {
	mv.tagName = tag
}

// WithTag creates a new Validator based on the current validator with
// the new tag name.
func (mv *Validator) WithTag(tag string) *Validator {
	v := mv.copy()
	v.SetTag(tag)
	return v
}

// Copy creates a duplicate of the current validator and returns
// the new instance
func (mv *Validator) copy() *Validator {
	return &Validator{
		tagName:         mv.tagName,
		validationFuncs: mv.validationFuncs,
	}
}

// SetValidationFunc sets the function to be used for a given
// validation constraint. Calling this function with nil vf
// is the same as removing the constraint function from the list.
func (mv *Validator) SetValidationFunc(name string, vf ValidationFunc) error {
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

// Validate validates the fields of a struct based
// on 'validator' tags and returns errors found indexed
// by the field name.
func (mv *Validator) Validate(v interface{}) error {
	sv := reflect.ValueOf(v)
	st := reflect.TypeOf(v)
	if sv.Kind() == reflect.Ptr && !sv.IsNil() {
		return mv.Validate(sv.Elem().Interface())
	}
	if sv.Kind() != reflect.Struct {
		return ErrUnsupported
	}

	nfields := sv.NumField()
	m := make(ErrorMap)
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
		var errs Errors

		if tag != "" {
			err := mv.Valid(f.Interface(), tag)
			if errors, ok := err.(Errors); ok {
				errs = errors
			} else {
				if err != nil {
					errs = Errors{err}
				}
			}
		}
		if f.Kind() == reflect.Struct {
			if !unicode.IsUpper(rune(fname[0])) {
				continue
			}
			e := mv.Validate(f.Interface())
			if e, ok := e.(ErrorMap); ok && len(e) > 0 {
				for j, k := range e {
					m[fname+"."+j] = k
				}
			}
		}
		if len(errs) > 0 {
			m[st.Field(i).Name] = errs
		}
	}
	if len(m) > 0 {
		return m
	}
	return nil
}

// Valid validates a value based on the provided
// tags and returns errors found or nil.
func (mv *Validator) Valid(val interface{}, tags string) error {
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

// validateVar validates one single variable
func (mv *Validator) validateVar(v interface{}, tag string) error {
	tags, err := mv.parseTags(tag)
	if err != nil {
		// unknown tag found.
		return err
	}
	errs := make(Errors, 0, len(tags))
	for _, t := range tags {
		if err := t.Fn(v, t.Param); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// parseTags parses all individual tags found within a struct tag.
func (mv *Validator) parseTags(t string) ([]tag, error) {
	tl := strings.Split(t, ",")
	tags := make([]tag, 0, len(tl))
	for _, i := range tl {
		tg := tag{}
		v := strings.SplitN(i, "=", 2)
		tg.Name = strings.Trim(v[0], " ")
		if tg.Name == "" {
			return []tag{}, ErrUnknownTag
		}
		if len(v) > 1 {
			tg.Param = strings.Trim(v[1], " ")
		}
		var found bool
		if tg.Fn, found = mv.validationFuncs[tg.Name]; !found {
			return []tag{}, ErrUnknownTag
		}
		tags = append(tags, tg)

	}
	return tags, nil
}
