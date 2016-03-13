package validate_test

import (
	"github.com/mbict/go-validate"
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type ValidatorSuite struct{}

var _ = Suite(&ValidatorSuite{})

type testSimple struct {
	A int `validate:"min(10)"`
}

type testStruct struct {
	A int `validate:"required"`
	B string
	C float64 `validate:"required,min(1)"`
	D *string `validate:"required"`
}

type testModel struct {
	A   int    `validate:"required"`
	B   string `validate:"len(8),min(6),max(4)"`
	Sub testStruct
	D   *testSimple  `validate:"required"`
	E   []testStruct `validate:"min(3),max(6)"` //<-- add a poniter variant too
}

func (ms *ValidatorSuite) TestValidate(c *C) {
	ptrEmptyStr := ""
	ptrString := "abc"
	t := testModel{
		A: 0,
		B: "abcdefg",
		Sub: testStruct{
			A: 0,
			B: "x",
			C: 0.2,
		},
		D: &testSimple{2},
		E: []testStruct{
			testStruct{
				A: 0,
				B: "x",
				C: 10.0,
				D: &ptrString,
			},
			testStruct{
				A: 5,
				B: "y",
				C: 0.0,
				D: &ptrEmptyStr,
			},
		},
	}

	err := validate.Validate(t)
	c.Assert(err, NotNil)

	errs, ok := err.(validate.ErrorMap)

	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
	c.Assert(errs["B"], HasLen, 2)
	c.Assert(errs["B"], HasError, validate.ErrLen)
	c.Assert(errs["B"], HasError, validate.ErrMax)
	c.Assert(errs["Sub.A"], HasLen, 1)
	c.Assert(errs["Sub.A"], HasError, validate.ErrRequired)
	c.Assert(errs["Sub.B"], HasLen, 0)
	c.Assert(errs["Sub.C"], HasLen, 1)
	c.Assert(errs["Sub.C"], HasError, validate.ErrMin)
	c.Assert(errs["Sub.D"], HasLen, 1)
	c.Assert(errs["Sub.D"], HasError, validate.ErrRequired)
	c.Assert(errs["D"], HasLen, 0)
	c.Assert(errs["D.A"], HasLen, 1)
	c.Assert(errs["D.A"], HasError, validate.ErrMin)
	c.Assert(errs["E"], HasLen, 1)
	c.Assert(errs["E"], HasError, validate.ErrMin)

	c.Assert(errs["E.0.A"], HasLen, 1)
	c.Assert(errs["E.0.A"], HasError, validate.ErrRequired)
	c.Assert(errs["E.0.B"], HasLen, 0)
	c.Assert(errs["E.0.C"], HasLen, 0)
	c.Assert(errs["E.0.D"], HasLen, 0)

	c.Assert(errs["E.1.A"], HasLen, 0)
	c.Assert(errs["E.1.B"], HasLen, 0)
	c.Assert(errs["E.1.C"], HasLen, 2)
	c.Assert(errs["E.1.C"], HasError, validate.ErrRequired)
	c.Assert(errs["E.1.C"], HasError, validate.ErrMin)
	c.Assert(errs["E.1.D"], HasLen, 1)
	c.Assert(errs["E.1.D"], HasError, validate.ErrRequired)

}

/*

func (ms *ValidatorSuite) TestValidSlice(c *C) {
	s := make([]int, 0, 10)
	err := validator.Valid(s, "required")
	c.Assert(err, NotNil)
	errs, ok := err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrRequired)

	for i := 0; i < 10; i++ {
		s = append(s, i)
	}

	err = validator.Valid(s, "min=11,max=5,len=9,required")
	c.Assert(err, NotNil)
	errs, ok = err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMin)
	c.Assert(errs, HasError, validate.ErrMax)
	c.Assert(errs, HasError, validate.ErrLen)
	c.Assert(errs, Not(HasError), validate.ErrRequired)
}

func (ms *ValidatorSuite) TestValidMap(c *C) {
	m := make(map[string]string)
	err := validator.Valid(m, "required")
	c.Assert(err, NotNil)
	errs, ok := err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrRequired)

	err = validator.Valid(m, "min=1")
	c.Assert(err, NotNil)
	errs, ok = err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMin)

	m = map[string]string{"A": "a", "B": "a"}
	err = validator.Valid(m, "max=1")
	c.Assert(err, NotNil)
	errs, ok = err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMax)

	err = validator.Valid(m, "min=2, max=5")
	c.Assert(err, IsNil)

	m = map[string]string{
		"1": "a",
		"2": "b",
		"3": "c",
		"4": "d",
		"5": "e",
	}
	err = validator.Valid(m, "len=4,min=6,max=1,required")
	c.Assert(err, NotNil)
	errs, ok = err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrLen)
	c.Assert(errs, HasError, validate.ErrMin)
	c.Assert(errs, HasError, validate.ErrMax)
	c.Assert(errs, Not(HasError), validate.ErrRequired)

}

func (ms *ValidatorSuite) TestValidFloat(c *C) {
	err := validator.Valid(12.34, "required")
	c.Assert(err, IsNil)

	err = validator.Valid(0.0, "required")
	c.Assert(err, NotNil)
	errs, ok := err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrRequired)
}

func (ms *ValidatorSuite) TestValidInt(c *C) {
	i := 123
	err := validator.Valid(i, "required")
	c.Assert(err, IsNil)

	err = validator.Valid(i, "min=1")
	c.Assert(err, IsNil)

	err = validator.Valid(i, "min=124, max=122")
	c.Assert(err, NotNil)
	errs, ok := err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMin)
	c.Assert(errs, HasError, validate.ErrMax)

	err = validator.Valid(i, "max=10")
	c.Assert(err, NotNil)
	errs, ok = err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMax)
}

func (ms *ValidatorSuite) TestValidString(c *C) {
	s := "test1234"
	err := validator.Valid(s, "len=8")
	c.Assert(err, IsNil)

	err = validator.Valid(s, "len=0")
	c.Assert(err, NotNil)
	errs, ok := err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrLen)

	err = validator.Valid(s, "regexp=^[tes]{4}.*")
	c.Assert(err, IsNil)

	err = validator.Valid(s, "regexp=^.*[0-9]{5}$")
	c.Assert(errs, NotNil)

	err = validator.Valid("", "required,len=3,max=1")
	c.Assert(err, NotNil)
	errs, ok = err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs, HasError, validate.ErrRequired)
	c.Assert(errs, HasError, validate.ErrLen)
	c.Assert(errs, Not(HasError), validate.ErrMax)
}

func (ms *ValidatorSuite) TestValidateStructVar(c *C) {
	// just verifies that a the given val is a struct
	validator.SetValidationFunc("struct", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Struct {
			return nil
		}
		return validate.ErrUnsupported
	})

	type test struct {
		A int
	}
	err := validator.Valid(test{}, "struct")
	c.Assert(err, IsNil)

	type test2 struct {
		B int
	}
	type test1 struct {
		A test2 `validate:"struct"`
	}

	err = validator.Validate(test1{})
	c.Assert(err, IsNil)

	type test4 struct {
		B int `validate:"foo"`
	}
	type test3 struct {
		A test4
	}
	err = validator.Validate(test3{})
	errs, ok := err.(validate.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A.B"], HasError, validate.ErrUnknownTag)
}

func (ms *ValidatorSuite) TestValidatePointerVar(c *C) {
	// just verifies that a the given val is a struct
	validator.SetValidationFunc("struct", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.Kind() == reflect.Struct {
			return nil
		}
		return validate.ErrUnsupported
	})
	validator.SetValidationFunc("nil", func(val interface{}, _ string) error {
		v := reflect.ValueOf(val)
		if v.IsNil() {
			return nil
		}
		return validate.ErrUnsupported
	})

	type test struct {
		A int
	}
	err := validator.Valid(&test{}, "struct")
	c.Assert(err, IsNil)

	type test2 struct {
		B int
	}
	type test1 struct {
		A *test2 `validate:"struct"`
	}

	err = validator.Validate(&test1{&test2{}})
	c.Assert(err, IsNil)

	type test4 struct {
		B int `validate:"foo"`
	}
	type test3 struct {
		A test4
	}
	err = validator.Validate(&test3{})
	errs, ok := err.(validate.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A.B"], HasError, validate.ErrUnknownTag)

	err = validator.Valid((*test)(nil), "nil")
	c.Assert(err, IsNil)

	type test5 struct {
		A *test2 `validate:"nil"`
	}
	err = validator.Validate(&test5{})
	c.Assert(err, IsNil)

	type test6 struct {
		A *test2 `validate:"required"`
	}
	err = validator.Validate(&test6{})
	errs, ok = err.(validate.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasError, validate.ErrRequired)

	err = validator.Validate(&test6{&test2{}})
	c.Assert(err, IsNil)
}

func (ms *ValidatorSuite) TestValidateOmittedStructVar(c *C) {
	type test2 struct {
		B int `validate:"min=1"`
	}
	type test1 struct {
		A test2 `validate:"-"`
	}

	t := test1{}
	err := validator.Validate(t)
	c.Assert(err, IsNil)

	errs := validator.Valid(test2{}, "-")
	c.Assert(errs, IsNil)
}

func (ms *ValidatorSuite) TestUnknownTag(c *C) {
	type test struct {
		A int `validate:"foo"`
	}
	t := test{}
	err := validator.Validate(t)
	c.Assert(err, NotNil)
	errs, ok := err.(validate.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrUnknownTag)
}

func (ms *ValidatorSuite) TestUnsupported(c *C) {
	type test struct {
		A int     `validate:"regexp=a.*b"`
		B float64 `validate:"regexp=.*"`
	}
	t := test{}
	err := validator.Validate(t)
	c.Assert(err, NotNil)
	errs, ok := err.(validate.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs["A"], HasError, validate.ErrUnsupported)
	c.Assert(errs["B"], HasError, validate.ErrUnsupported)
}

func (ms *ValidatorSuite) TestBadParameter(c *C) {
	type test struct {
		A string `validate:"min="`
		B string `validate:"len=="`
		C string `validate:"max=foo"`
	}
	t := test{}
	err := validator.Validate(t)
	c.Assert(err, NotNil)
	errs, ok := err.(validate.ErrorMap)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 3)
	c.Assert(errs["A"], HasError, validate.ErrBadParameter)
	c.Assert(errs["B"], HasError, validate.ErrBadParameter)
	c.Assert(errs["C"], HasError, validate.ErrBadParameter)
}
*/

type hasErrorChecker struct {
	*CheckerInfo
}

func (c *hasErrorChecker) Check(params []interface{}, names []string) (bool, string) {
	var (
		ok    bool
		slice []error
		value error
	)
	slice, ok = params[0].(validate.Errors)
	if !ok {
		return false, "First parameter is not an Errorarray"
	}
	value, ok = params[1].(error)
	if !ok {
		return false, "Second parameter is not an error"
	}

	for _, v := range slice {
		if v == value {
			return true, ""
		}
	}
	return false, ""
}

func (c *hasErrorChecker) Info() *CheckerInfo {
	return c.CheckerInfo
}

var HasError = &hasErrorChecker{&CheckerInfo{Name: "HasError", Params: []string{"HasError", "expected to contain"}}}
