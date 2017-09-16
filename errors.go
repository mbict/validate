package validate

import (
	"github.com/mbict/go-errors"
)

var (
	// ErrRequired is the error returned when variable has a empty value
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

	// ErrBetween is the error returned when variable is less than
	// minimum or more than maximum specified
	ErrBetween = errors.New("not between")

	// ErrAround is the error returned when variable is more than
	// minimum and less then maximum specified
	ErrAround = errors.New("not around")

	// ErrRegexp is the error returned when the value does not
	// match the provided regular expression parameter
	ErrRegexp = errors.New("regular expression mismatch")

	// ErrAlphaDash is the error returned when the value does contains
	// other characters than alpha's or dashes
	ErrAlphaDash = errors.New("alpha dash mismatch")

	// ErrAlphaDashDot is the error returned when the value does contains
	// other characters than alpha's, dashes or dots
	ErrAlphaDashDot = errors.New("alpha dash dot mismatch")

	// ErrEmail is the error returned when the value does not match
	// a valid email pattern
	ErrEmail = errors.New("invalid email")

	// ErrURL is the error returned when the value does not match
	// a valid url pattern
	ErrURL = errors.New("invalid url")

	// ErrInclude is the error returned when the value is not found
	// in the set values that in the include list
	ErrInclude = errors.New("value not found in set")

	// ErrExclude is the error returned when the value is  found
	// in the set values to exclude
	ErrExclude = errors.New("value matches one exluded value")

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
