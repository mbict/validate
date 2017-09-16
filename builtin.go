package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// nonzero tests whether a variable value non-zero
// as defined by the golang spec.
func required(v interface{}, params []string) error {
	st := reflect.ValueOf(v)
	valid := true
	switch st.Kind() {
	case reflect.String:
		valid = len(st.String()) != 0
	case reflect.Ptr, reflect.Interface:
		valid = !st.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		valid = st.Len() != 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid = st.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid = st.Uint() != 0
	case reflect.Float32, reflect.Float64:
		valid = st.Float() != 0
	case reflect.Bool:
		valid = st.Bool()
	case reflect.Invalid:
		valid = false // always invalid
	case reflect.Struct:
		valid = true // always valid since only nil pointers are empty
	default:
		return ErrUnsupported
	}

	if !valid {
		return ErrRequired
	}
	return nil
}

// length tests whether a variable's length is equal to a given
// value. For strings it tests the number of characters whereas
// for maps and slices it tests the number of items.
func length(v interface{}, params []string) error {
	if len(params) != 1 {
		return ErrInvalidParameterCount
	}

	st := reflect.ValueOf(v)
	valid := true
	switch st.Kind() {
	case reflect.String:
		p, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}
		valid = int64(len(st.String())) == p
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}
		valid = int64(st.Len()) == p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}
		valid = st.Int() == p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(params[0])
		if err != nil {
			return ErrBadParameter
		}
		valid = st.Uint() == p
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(params[0])
		if err != nil {
			return ErrBadParameter
		}
		valid = st.Float() == p
	default:
		return ErrUnsupported
	}
	if !valid {
		return ErrLen
	}
	return nil
}

// min tests whether a variable value is larger or equal to a given
// number. For number types, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func min(v interface{}, params []string) error {
	if len(params) != 1 {
		return ErrInvalidParameterCount
	}

	st := reflect.ValueOf(v)
	invalid := false
	switch st.Kind() {
	case reflect.String:
		p, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}
		invalid = int64(len(st.String())) < p
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}
		invalid = int64(st.Len()) < p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}
		invalid = st.Int() < p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(params[0])
		if err != nil {
			return ErrBadParameter
		}
		invalid = st.Uint() < p
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(params[0])
		if err != nil {
			return ErrBadParameter
		}
		invalid = st.Float() < p
	default:
		return ErrUnsupported
	}
	if invalid {
		return ErrMin
	}
	return nil
}

// max tests whether a variable value is lesser than a given
// value. For numbers, it's a simple lesser-than test; for
// strings it tests the number of characters whereas for maps
// and slices it tests the number of items.
func max(v interface{}, params []string) error {

	if len(params) != 1 {
		return ErrInvalidParameterCount
	}

	st := reflect.ValueOf(v)
	var invalid bool
	switch st.Kind() {
	case reflect.String:
		p, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}
		invalid = int64(len(st.String())) > p
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}
		invalid = int64(st.Len()) > p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}
		invalid = st.Int() > p
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := asUint(params[0])
		if err != nil {
			return ErrBadParameter
		}
		invalid = st.Uint() > p
	case reflect.Float32, reflect.Float64:
		p, err := asFloat(params[0])
		if err != nil {
			return ErrBadParameter
		}
		invalid = st.Float() > p
	default:
		return ErrUnsupported
	}
	if invalid {
		return ErrMax
	}
	return nil
}

// regex is the builtin validation function that checks
// whether the string variable matches a regular expression
func regex(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if len(params) != 1 {
		return ErrInvalidParameterCount
	}

	re, err := regexp.Compile(params[0])
	if err != nil {
		return ErrBadParameter
	}

	if !re.MatchString(s) {
		return ErrRegexp
	}
	return nil
}

func between(v interface{}, params []string) error {
	if len(params) != 2 {
		return ErrInvalidParameterCount
	}

	st := reflect.ValueOf(v)
	var invalid bool
	switch st.Kind() {
	case reflect.String:
		a, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}

		b, err := asInt(params[1])
		if err != nil {
			return ErrBadParameter
		}

		len := int64(len(st.String()))
		if a > b {
			invalid = len < b || len > a
		} else { //inverse
			invalid = len < a || len > b
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		a, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}

		b, err := asInt(params[1])
		if err != nil {
			return ErrBadParameter
		}

		len := int64(st.Len())
		if a > b {
			invalid = len < b || len > a
		} else { //inverse
			invalid = len < a || len > b
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		a, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}

		b, err := asInt(params[1])
		if err != nil {
			return ErrBadParameter
		}

		val := st.Int()
		if a > b {
			invalid = val < b || val > a
		} else { //inverse
			invalid = val < a || val > b
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		a, err := asUint(params[0])
		if err != nil {
			return ErrBadParameter
		}

		b, err := asUint(params[1])
		if err != nil {
			return ErrBadParameter
		}

		val := st.Uint()
		if a > b {
			invalid = val < b || val > a
		} else { //inverse
			invalid = val < a || val > b
		}
	case reflect.Float32, reflect.Float64:
		a, err := asFloat(params[0])
		if err != nil {
			return ErrBadParameter
		}

		b, err := asFloat(params[1])
		if err != nil {
			return ErrBadParameter
		}

		val := st.Float()
		if a > b {
			invalid = val < b || val > a
		} else { //inverse
			invalid = val < a || val > b
		}

	default:
		return ErrUnsupported
	}
	if invalid {
		return ErrBetween
	}
	return nil
}

func around(v interface{}, params []string) error {
	if len(params) != 2 {
		return ErrInvalidParameterCount
	}

	st := reflect.ValueOf(v)
	var invalid bool
	switch st.Kind() {
	case reflect.String:
		a, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}

		b, err := asInt(params[1])
		if err != nil {
			return ErrBadParameter
		}

		len := int64(len(st.String()))
		if a < b {
			invalid = len < b && len > a
		} else { //inverse
			invalid = len < a && len > b
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		a, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}

		b, err := asInt(params[1])
		if err != nil {
			return ErrBadParameter
		}

		len := int64(st.Len())
		if a < b {
			invalid = len < b && len > a
		} else { //inverse
			invalid = len < a && len > b
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		a, err := asInt(params[0])
		if err != nil {
			return ErrBadParameter
		}

		b, err := asInt(params[1])
		if err != nil {
			return ErrBadParameter
		}

		val := st.Int()
		if a < b {
			invalid = val < b && val > a
		} else { //inverse
			invalid = val < a && val > b
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		a, err := asUint(params[0])
		if err != nil {
			return ErrBadParameter
		}

		b, err := asUint(params[1])
		if err != nil {
			return ErrBadParameter
		}

		val := st.Uint()
		if a < b {
			invalid = val < b && val > a
		} else { //inverse
			invalid = val < a && val > b
		}
	case reflect.Float32, reflect.Float64:
		a, err := asFloat(params[0])
		if err != nil {
			return ErrBadParameter
		}

		b, err := asFloat(params[1])
		if err != nil {
			return ErrBadParameter
		}

		val := st.Float()
		if a < b {
			invalid = val < b && val > a
		} else { //inverse
			invalid = val < a && val > b
		}

	default:
		return ErrUnsupported
	}
	if invalid {
		return ErrAround
	}
	return nil
}

func include(i interface{}, params []string) error {
	if len(params) < 1 {
		return ErrInvalidParameterCount
	}

	str := fmt.Sprintf("%v", i)
	for _, in := range params {
		if strings.Compare(str, in) == 0 {
			//we found a match
			return nil
		}
	}

	//if no match is found we error out
	return ErrInclude
}

func exclude(i interface{}, params []string) error {
	if len(params) < 1 {
		return ErrInvalidParameterCount
	}

	str := fmt.Sprintf("%v", i)
	for _, in := range params {
		if strings.Compare(str, in) == 0 {
			//found a match, we error out
			return ErrExclude
		}
	}

	//if no match is found all is ok
	return nil
}

func number(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if numberRegex.MatchString(s) {
		return ErrNumber
	}
	return nil
}

func numeric(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if numericRegex.MatchString(s) {
		return ErrNumeric
	}
	return nil
}

func alpha(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if alphaRegex.MatchString(s) {
		return ErrAlpha
	}
	return nil
}

func alphaNumeric(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if alphaNumericRegex.MatchString(s) {
		return ErrAlphaNumeric
	}
	return nil
}

func alphaDash(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if alphaDashRegex.MatchString(s) {
		return ErrAlphaDash
	}
	return nil
}

func alphaDashDot(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if alphaDashDotRegex.MatchString(s) {
		return ErrAlphaDashDot
	}
	return nil
}

func email(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if !emailRegex.MatchString(s) {
		return ErrEmail
	}
	return nil
}

func url(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if !urlRegex.MatchString(s) {
		return ErrURL
	}
	return nil
}

func uuid(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if !uUIDRegex.MatchString(s) {
		return ErrUUID
	}
	return nil
}

func uuid3(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if !uUID3Regex.MatchString(s) {
		return ErrUUID3
	}
	return nil
}

func uuid4(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if !uUID4Regex.MatchString(s) {
		return ErrUUID4
	}
	return nil
}

func uuid5(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if !uUID5Regex.MatchString(s) {
		return ErrUUID5
	}
	return nil
}

func base64(v interface{}, params []string) error {
	s, ok := v.(string)
	if !ok {
		return ErrUnsupported
	}

	if base64Regex.MatchString(s) {
		return ErrBase64
	}
	return nil
}

// asInt returns the parameter as a int64
// or panics if it can't convert
func asInt(param string) (int64, error) {
	i, err := strconv.ParseInt(param, 0, 64)
	if err != nil {
		return 0, ErrBadParameter
	}
	return i, nil
}

// asUint returns the parameter as a uint64
// or panics if it can't convert
func asUint(param string) (uint64, error) {
	i, err := strconv.ParseUint(param, 0, 64)
	if err != nil {
		return 0, ErrBadParameter
	}
	return i, nil
}

// asFloat returns the parameter as a float64
// or panics if it can't convert
func asFloat(param string) (float64, error) {
	i, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return 0.0, ErrBadParameter
	}
	return i, nil
}
