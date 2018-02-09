package validate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type ValidationError error

type validationError struct {
	template string
	args     []interface{}
}

func (e validationError) Template() string {
	return e.template
}

func (e validationError) Args() []interface{} {
	return e.args
}

func (e validationError) Error() string {
	return fmt.Sprintf(e.template, e.args...)
}

func (e validationError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Error())
}

func (e *validationError) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &e.template)
}

func NewValidationError(template string, args ...interface{}) ValidationError {
	return &validationError{
		template: template,
		args:     args,
	}
}

var (
	// errOmitEmpty is the error returned when variable has a empty value
	errOmitEmpty = NewValidationError("omitempty")

	// ErrRequired is the error returned when variable has a empty value
	ErrRequired = NewValidationError("required")

	// ErrEmpty is the error returned when variable has a empty value
	ErrEmpty = NewValidationError("value is empty")

	// ErrMin is the error returned when variable is less than mininum
	// value specified
	ErrMin = NewValidationError("less than min")

	// ErrMax is the error returned when variable is more than
	// maximum specified
	ErrMax = NewValidationError("greater than max")

	// ErrLen is the error returned when length is not equal to
	// param specified
	ErrLen = NewValidationError("invalid length")

	// ErrBetween is the error returned when variable is less than
	// minimum or more than maximum specified
	ErrBetween = NewValidationError("not between")

	// ErrAround is the error returned when variable is more than
	// minimum and less then maximum specified
	ErrAround = NewValidationError("not around")

	// ErrRegexp is the error returned when the value does not
	// match the provided regular expression parameter
	ErrRegexp = NewValidationError("regular expression mismatch")

	// ErrIdentifier is the error returned when the value does not match the
	// identifier pattern
	ErrIdentifier = NewValidationError("invalid id format")

	// ErrAlpha is the error returned when the value does contains
	// other characters than alphas
	ErrAlpha = NewValidationError("alpha dash mismatch")

	// ErrAlphaNumeric is the error returned when the value does contains
	// other characters than alphas or numerics
	ErrAlphaNumeric = NewValidationError("alpha dash mismatch")

	// ErrAlphaDash is the error returned when the value does contains
	// other characters than alphas or dashes
	ErrAlphaDash = NewValidationError("alpha dash mismatch")

	// ErrAlphaDashDot is the error returned when the value does contains
	// other characters than alpha's, dashes or dots
	ErrAlphaDashDot = NewValidationError("alpha dash dot mismatch")

	// ErrEmail is the error returned when the value does not match
	// a valid email pattern
	ErrEmail = NewValidationError("invalid email")

	// ErrURL is the error returned when the value does not match
	// a valid url pattern
	ErrURL = NewValidationError("invalid url")

	// ErrInclude is the error returned when the value is not found
	// in the set values that in the include list
	ErrInclude = NewValidationError("value not found in set")

	// ErrExclude is the error returned when the value is  found
	// in the set values to exclude
	ErrExclude = NewValidationError("value matches one exluded value")

	// ErrUnsupported is the error error returned when a validation rule
	// is used with an unsupported variable type
	ErrUnsupported = NewValidationError("unsupported type")

	// ErrBadParameter is the error returned when an invalid parameter
	// is provided to a validation rule (e.g. a string where an int was
	// expected (max(foo),len=(bar))
	ErrBadParameter = NewValidationError("bad parameter")

	// ErrInvalidParameterCount is the error returned when there are not enough or
	// to many parameters provided to the validation rule.
	ErrInvalidParameterCount = NewValidationError("invalid parameter count")

	// ErrSyntax is the error who is returned when a invalid syntax is detected
	// while parsing the structure validatorTag
	ErrSyntax = NewValidationError("syntax error")

	// ErrUnknownTag is the error returned when an unknown validatorTag is found
	ErrUnknownTag = NewValidationError("unknown validatorTag")

	// ErrInvalid is the error returned when variable is invalid
	// (normally a nil pointer)
	ErrInvalid = NewValidationError("invalid value")

	// ErrNumber is the error returned when value is not a number
	ErrNumber = NewValidationError("value not a number")

	// ErrNumeric is the error returned when value is not in numeric format
	ErrNumeric = NewValidationError("value not numeric")

	// ErrUUID is the error returned when value is not in UUID format
	ErrUUID = NewValidationError("invalid UUID")

	// ErrUUID3 is the error returned when value is not in UUID format
	ErrUUID3 = NewValidationError("invalid UUID3")

	// ErrUUID4 is the error returned when value is not in UUID format
	ErrUUID4 = NewValidationError("invalid UUID4")

	// ErrUUID5 is the error returned when value is not in UUID format
	ErrUUID5 = NewValidationError("invalid UUID5")

	// ErrBase64 is the error returned when value is not a valid base64 encoded
	ErrBase64 = NewValidationError("invalid base64 encoded")

	// ErrBEnum is the error returned when value is not in a set of enum values
	ErrEnum = NewValidationError("invalid value")
)

type ErrorList []error

func (e ErrorList) Error() string {
	errs := make([]string, len(e))
	for i, err := range e {
		errs[i] = err.Error()
	}
	return strings.Join(errs, ", ")
}

func (e ErrorList) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')

	//stringify errors
	for i, es := range e {
		if i != 0 {
			buf.WriteByte(',')
		}
		val, err := json.Marshal(es.Error())
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}
	buf.WriteByte(']')

	return buf.Bytes(), nil
}

type Errors map[string][]error

func (e Errors) Error() string {
	names := make([]string, 0, len(e))
	for name := range e {
		names = append(names, name)
	}
	sort.Strings(names)

	stringErrors := make([]string, len(e))
	for ei, name := range names {
		errs := make([]string, len(e[name]))
		for ie, err := range e[name] {
			errs[ie] = err.Error()
		}
		stringErrors[ei] = fmt.Sprintf("%s: [%s]", name, strings.Join(errs, ", "))
	}
	return strings.Join(stringErrors, ", ")
}

func (e *Errors) Add(field string, errors ...error) {
	if *e == nil {
		*e = make(Errors, 1)
	}

	(*e)[field] = append((*e)[field], errors...)
}

func (e *Errors) Merge(errors error) {
	switch verr := errors.(type) {
	case Errors:
		for field, errs := range verr {
			e.Add(field, errs...)
		}
	default:
		e.Add("_", errors)
	}
}

func (e *Errors) MergePrefix(prefix string, errors error) {
	switch verr := errors.(type) {
	case Errors:
		for field, errs := range verr {
			e.Add(prefix+field, errs...)
		}
		//case Error:
		//	e.Add(prefix+verr.Field(), verr.Errors()...)
	default:
		e.Add(prefix+"_", errors)
	}
}

func (e Errors) MarshalJSON() ([]byte, error) {
	//order alphabetic field/namess
	names := make([]string, 0, len(e))
	for name := range e {
		names = append(names, name)
	}
	sort.Strings(names)

	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, name := range names {
		if i != 0 {
			buf.WriteByte(',')
		}

		//name/field who contains errors
		key, err := json.Marshal(name)
		if err != nil {
			return nil, err
		}
		buf.Write(key)

		buf.WriteString(":[")

		//stringify errors
		for esi, es := range e[name] {
			if esi != 0 {
				buf.WriteByte(',')
			}
			val, err := json.Marshal(es.Error())
			if err != nil {
				return nil, err
			}
			buf.Write(val)
		}
		buf.WriteByte(']')

	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (e *Errors) UnmarshalJSON(data []byte) error {
	errs := map[string][]validationError{}
	err := json.Unmarshal(data, &errs)
	if err != nil {
		return err
	}

	*e = make(map[string][]error, len(errs))
	for k, v := range errs {
		(*e)[k] = make([]error, len(v))
		for i, err := range v {
			(*e)[k][i] = error(err)
		}
	}
	return nil
}
