package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ValidatorFunc func(reflect.Value, []string, string) Errors

const (
	RequiredError     = "RequiredError"
	AlphaDashError    = "AlphaDashError"
	AlphaDashDotError = "AlphaDashDotError"
	MinError          = "MinError"
	MaxError          = "MaxError"
	EmailError        = "EmailError"
	UrlError          = "UrlError"
	RangeError        = "RangeError"
	InError           = "InError"
	NotInError        = "NotInError"
	IncludeError      = "IncludeError"
	ExcludeError      = "ExcludeError"
)

var (
	alphaDashPattern    = regexp.MustCompile("[^\\d\\w-_]")
	alphaDashDotPattern = regexp.MustCompile("[^\\d\\w-_\\.]")
	emailPattern        = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")
	urlPattern          = regexp.MustCompile(`(http|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
)

var validators = map[string]ValidatorFunc{
	"required":     requiredValidator,
	"alphadash":    alphaDashValidator,
	"alphadashdot": alphaDashDotValidator,
	"min":          minValidator,
	"max":          maxValidator,
	"email":        emailValidator,
	"url":          urlValidator,
	"range":        rangeValidator,
	"include":      includeValidator,
	"exclude":      excludeValidator,
	"in":           inValidator,
	"notin":        notInValidator,
}

func requiredValidator(v reflect.Value, params []string, path string) (errors Errors) {
	zero := reflect.Zero(v.Type()).Interface()
	if reflect.DeepEqual(zero, v.Interface()) {
		errors.Add([]string{path}, RequiredError, "Required")
	}
	return
}

func alphaDashValidator(v reflect.Value, params []string, path string) (errors Errors) {
	if alphaDashPattern.MatchString(fmt.Sprintf("%v", v.Interface())) {
		errors.Add([]string{path}, AlphaDashError, "AlphaDash")
	}
	return
}

func alphaDashDotValidator(v reflect.Value, params []string, path string) (errors Errors) {
	if alphaDashDotPattern.MatchString(fmt.Sprintf("%v", v.Interface())) {
		errors.Add([]string{path}, AlphaDashDotError, "AlphaDashDot")
	}
	return
}

func minValidator(v reflect.Value, params []string, path string) (errors Errors) {
	if len(params) == 0 {
		return
	}
	min, _ := strconv.Atoi(params[0])
	if v.Kind() == reflect.String && utf8.RuneCountInString(v.String()) < min {
		errors.Add([]string{path}, MinError, "MinSize")
	} else if v.Kind() == reflect.Slice && v.Len() < min {
		errors.Add([]string{path}, MinError, "MinSize")
	}
	return
}

func maxValidator(v reflect.Value, params []string, path string) (errors Errors) {
	if len(params) == 0 {
		return
	}
	max, _ := strconv.Atoi(params[0])
	if v.Kind() == reflect.String && utf8.RuneCountInString(v.String()) > max {
		errors.Add([]string{path}, MaxError, "MaxSize")
	} else if v.Kind() == reflect.Slice && v.Len() > max {
		errors.Add([]string{path}, MaxError, "MaxSize")
	}
	return
}

func emailValidator(v reflect.Value, params []string, path string) (errors Errors) {
	str := fmt.Sprintf("%v", v.Interface())
	if len(str) > 0 && !emailPattern.MatchString(str) {
		errors.Add([]string{path}, EmailError, "Email")
	}
	return
}

func urlValidator(v reflect.Value, params []string, path string) (errors Errors) {
	str := fmt.Sprintf("%v", v.Interface())
	if len(str) > 0 && !urlPattern.MatchString(str) {
		errors.Add([]string{path}, UrlError, "Url")
	}
	return
}

func rangeValidator(v reflect.Value, params []string, path string) (errors Errors) {
	if len(params) < 2 {
		return
	}
	val, _ := strconv.ParseInt(fmt.Sprintf("%v", v.Interface()), 10, 32)
	a, _ := strconv.ParseInt(params[0], 10, 32)
	b, _ := strconv.ParseInt(params[1], 10, 32)
	if val < a || val > b {
		errors.Add([]string{path}, RangeError, "Range")
	}
	return
}

func inValidator(v reflect.Value, params []string, path string) (errors Errors) {
	if !in(fmt.Sprintf("%v", v.Interface()), params) {
		errors.Add([]string{path}, InError, "In")
	}
	return
}

func notInValidator(v reflect.Value, params []string, path string) (errors Errors) {
	if in(fmt.Sprintf("%v", v.Interface()), params) {
		errors.Add([]string{path}, NotInError, "NotIn")
	}
	return
}

func in(val string, arr []string) bool {
	isIn := false
	for _, v := range arr {
		if v == val {
			isIn = true
			break
		}
	}
	return isIn
}

func includeValidator(v reflect.Value, params []string, path string) (errors Errors) {
	if len(params) >= 2 {
		return
	}

	if !strings.Contains(fmt.Sprintf("%v", v.Interface()), params[0]) {
		errors.Add([]string{path}, IncludeError, "Include")
	}
	return
}

func excludeValidator(v reflect.Value, params []string, path string) (errors Errors) {
	if len(params) >= 2 {
		return
	}

	if strings.Contains(fmt.Sprintf("%v", v.Interface()), params[0]) {
		errors.Add([]string{path}, ExcludeError, "Exclude")
	}
	return
}
