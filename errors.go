package validate

import (
	"errors"
	"fmt"
)

var (
	// ErrRequiredValue is the error returned when variable has a empty value
	ErrRequired = errors.New("required")

	// ErrMin is the error returned when variable is less than mininum
	// value specified
	ErrMin = errors.New("less than min")

	// ErrMax is the error returned when variable is more than
	// maximum specified
	ErrMax = errors.New("greater than max")

	// ErrLen is the error returned when length is not equal to
	// param specified
	ErrLen = errors.New("invalid length")

	// ErrRegexp is the error returned when the value does not
	// match the provided regular expression parameter
	ErrRegexp = errors.New("regular expression mismatch")

	// ErrUnsupported is the error error returned when a validation rule
	// is used with an unsupported variable type
	ErrUnsupported = errors.New("unsupported type")

	// ErrBadParameter is the error returned when an invalid parameter
	// is provided to a validation rule (e.g. a string where an int was
	// expected (max(foo),len=(bar))
	ErrBadParameter = errors.New("bad parameter")

	// ErrInvalidParameterCount is the error returned when there are not enough or
	// to many parameters provided to the validation rule.
	ErrInvalidParameterCount = errors.New("invalid parameter count")

	// ErrSyntax is the error who is returned when a invalid syntax is detected
	// while parsing the structure tag
	ErrSyntax = errors.New("syntax error")

	// ErrUnknownTag is the error returned when an unknown tag is found
	ErrUnknownTag = errors.New("unknown tag")

	// ErrInvalid is the error returned when variable is invalid
	// (normally a nil pointer)
	ErrInvalid = errors.New("invalid value")
)

// Errors is a slice of errors returned by the Validate function.
type Errors []error

// Errors implements the Error interface and returns all the errors
// as a comma delimited string
func (errs Errors) Error() string {
	if len(errs) > 0 {
		result := errs[0].Error()
		for _, err := range errs[1:] {
			result = result + ", " + err.Error()
		}
		return result
	}
	return ""
}

// ErrorMap is a map which contains all errors from validating a struct.
type ErrorMap map[string]Errors

// ErrorMap implements the Error interface and returns all the fields
// who has errors in a cimma delimited string
func (err ErrorMap) Error() string {
	for k, errs := range err {
		if len(errs) > 0 {
			result := fmt.Sprintf("%s: %s", k, errs[0].Error())
			for _, err := range errs[1:] {
				result = result + ", " + fmt.Sprintf("%s: %s", k, err.Error())
			}
			return result
		}
	}
	return ""
}
