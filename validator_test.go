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

func (vs *ValidatorSuite) TestValidateAll(c *C) {
	emptyStr := ""
	fullStr := "abc"
	var test = []struct {
		A int     `validate:"required;min(3);max(6)"`
		B string  `validate:"required;min(3);max(6)"`
		C *string `validate:"required"`
		D bool    `validate:"required"`
		E float64 `validate:"required;min(1.5);max(3)"`
	}{
		{
			A: 0,
			B: "",
			C: nil,
			D: false,
			E: 0,
		}, {
			A: 10,
			B: "abcdefgh",
			C: &emptyStr,
			D: true,
			E: 500.12,
		}, {
			A: 4,
			B: "abcd",
			C: &fullStr,
			D: true,
			E: 2.5,
		},
	}

	//error, empty values & nil ptr
	err := validate.ValidateAll(test[0])
	c.Assert(err, NotNil)

	errs, ok := err.(validate.Errors)

	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasLen, 2)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
	c.Assert(errs["A"], HasError, validate.ErrMin)
	c.Assert(errs["B"], HasLen, 2)
	c.Assert(errs["B"], HasError, validate.ErrRequired)
	c.Assert(errs["B"], HasError, validate.ErrMin)
	c.Assert(errs["C"], HasLen, 1)
	c.Assert(errs["C"], HasError, validate.ErrRequired)
	c.Assert(errs["D"], HasLen, 1)
	c.Assert(errs["D"], HasError, validate.ErrRequired)
	c.Assert(errs["E"], HasLen, 2)
	c.Assert(errs["E"], HasError, validate.ErrRequired)
	c.Assert(errs["E"], HasError, validate.ErrMin)

	//error, invalid values and empty string ptr
	err = validate.ValidateAll(test[1])
	c.Assert(err, NotNil)

	errs, ok = err.(validate.Errors)

	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], NotNil)
	c.Assert(errs["A"], HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrMax)
	c.Assert(errs["B"], NotNil)
	c.Assert(errs["B"], HasLen, 1)
	c.Assert(errs["B"], HasError, validate.ErrMax)
	c.Assert(errs["C"], IsNil)
	//c.Assert(errs["C"], NotNil)
	//c.Assert(errs["C"].Errors(), HasLen, 1)
	//c.Assert(errs["C"], HasError, validate.ErrRequired)
	c.Assert(errs["D"], IsNil)
	c.Assert(errs["E"], NotNil)
	c.Assert(errs["E"], HasLen, 1)
	c.Assert(errs["E"], HasError, validate.ErrMax)

	//should pass
	err = validate.ValidateAll(test[2])
	c.Assert(err, IsNil)
}

func (vs *ValidatorSuite) TestValidateAllEmbedStruct(c *C) {
	var test = []struct {
		A testSimple
		B *testSimple `validate:"required"`
	}{
		{
			B: nil,
		}, {
			A: testSimple{1},
			B: &testSimple{3},
		}, {
			A: testSimple{11},
			B: &testSimple{12},
		},
	}

	//error, nil ptr / empty values
	err := validate.ValidateAll(test[0])
	c.Assert(err, NotNil)

	errs, ok := err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A.A"], HasLen, 1)
	c.Assert(errs["A.A"], HasError, validate.ErrMin)
	c.Assert(errs["B"], HasLen, 1)
	c.Assert(errs["B"], HasError, validate.ErrRequired)
	c.Assert(errs["B.A"], IsNil)

	//error, initialized ptr
	err = validate.ValidateAll(test[1])
	c.Assert(err, NotNil)

	errs, ok = err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs["B"], IsNil)
	c.Assert(errs["B.A"], HasLen, 1)
	c.Assert(errs["B.A"], HasError, validate.ErrMin)

	//error, initialized ptr
	err = validate.ValidateAll(test[1])
	c.Assert(err, NotNil)

	errs, ok = err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs["B"], IsNil)
	c.Assert(errs["B.A"], HasLen, 1)
	c.Assert(errs["B.A"], HasError, validate.ErrMin)

	//should pass
	err = validate.ValidateAll(test[2])
	c.Assert(err, IsNil)
}

func (vs *ValidatorSuite) TestValidateAllSlice(c *C) {
	var test = []struct {
		A []testSimple  `validate:"required,min(3),max(6)"`
		B []*testSimple `validate:"required,max(6)"`
	}{
		{
			A: nil,
			B: nil,
		}, {
			A: []testSimple{},
			B: []*testSimple{},
		}, {
			A: []testSimple{{1}, {11}},
			B: []*testSimple{&testSimple{11}, &testSimple{3}},
		}, {
			A: []testSimple{{11}, {12}, {13}},
			B: []*testSimple{&testSimple{11}},
		},
	}

	//error, nil slices
	err := validate.ValidateAll(test[0])
	c.Assert(err, NotNil)

	errs, ok := err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasLen, 2)
	c.Assert(errs["A"], HasError, validate.ErrMin)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
	c.Assert(errs["B"], HasLen, 1)
	c.Assert(errs["B"], HasError, validate.ErrRequired)

	//error, zero length slices
	err = validate.ValidateAll(test[1])
	c.Assert(err, NotNil)

	errs, ok = err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasLen, 2)
	c.Assert(errs["A"], HasError, validate.ErrMin)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
	c.Assert(errs["B"], HasLen, 1)
	c.Assert(errs["B"], HasError, validate.ErrRequired)

	//error, filled slices
	err = validate.ValidateAll(test[2])
	c.Assert(err, NotNil)

	errs, ok = err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrMin)
	c.Assert(errs["A.0.A"], HasLen, 1)
	c.Assert(errs["A.0.A"], HasError, validate.ErrMin)
	c.Assert(errs["B"], IsNil)
	c.Assert(errs["B.0.A"], IsNil)
	c.Assert(errs["B.1.A"], HasLen, 1)
	c.Assert(errs["B.1.A"], HasError, validate.ErrMin)

	//should pass
	err = validate.ValidateAll(test[3])
	c.Assert(err, IsNil)
}

func (vs *ValidatorSuite) TestValidateAllIgnoreNonExportedVars(c *C) {
	var test = []struct {
		A int `validate:"required"`
		b int `validate:"required"`
	}{
		{
			A: 0,
			b: 0,
		}, {
			A: 1,
			b: 1,
		},
	}

	//error, nil ptr / empty values
	err := validate.ValidateAll(test[0])
	c.Assert(err, NotNil)

	errs, ok := err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
	c.Assert(errs["b"], IsNil)

	//error, initialized ptr
	err = validate.ValidateAll(test[1])
	c.Assert(err, IsNil)
}

func (vs *ValidatorSuite) TestValidateAllInputNilValue(c *C) {
	err := validate.ValidateAll(nil)

	c.Assert(err, Equals, validate.ErrUnsupported)
}

func (vs *ValidatorSuite) TestValidateAllInputTypedNilValue(c *C) {
	err := validate.ValidateAll((*testSimple)(nil))

	c.Assert(err, Equals, validate.ErrUnsupported)
}

func (vs *ValidatorSuite) TestValidErrorUnkownTag(c *C) {
	err := validate.ValidAll(1, "min(10);nonexistingvalidator(1,2,3)")

	c.Assert(err, NotNil)
	c.Assert(err, Equals, validate.ErrUnknownTag)
}

func (vs *ValidatorSuite) TestValidErrorSyntax(c *C) {
	err := validate.ValidAll(1, "min(10)|bv")

	c.Assert(err, NotNil)
	c.Assert(err, Equals, validate.ErrSyntax)
}

func (vs *ValidatorSuite) TestValidIgnoreTagReturnNil(c *C) {
	err := validate.ValidAll(1, "-")

	c.Assert(err, IsNil)
}

func (vs *ValidatorSuite) TestValidateAllIgnoreTag(c *C) {
	test := struct {
		A testSimple `validate:"-"`
	}{}

	err := validate.ValidateAll(test)

	c.Assert(err, IsNil)
}

type testStructValidateFuncInterface struct {
	A            int `validate:"min(2)"`
	B            int `validate:"min(2)"`
	C            int `validate:"min(10)"`
	validateFunc func() error
}

func (s testStructValidateFuncInterface) Validate() error {
	return s.validateFunc()
}

func (vs *ValidatorSuite) TestStructValidateInterface(c *C) {
	customErr1 := validate.NewValidationError("custom 1")
	customErr2 := validate.NewValidationError("cutsom 2")

	test := testStructValidateFuncInterface{
		A: 1,
		B: 2,
		C: 3,
		validateFunc: func() error {
			return validate.Errors{
				"A": validate.ErrorList{customErr1, customErr2},
				"B": validate.ErrorList{customErr1},
				"D": validate.ErrorList{customErr2},
			}
		},
	}

	errs := validate.ValidateAll(test).(validate.Errors)
	c.Assert(errs, HasLen, 4)
	c.Assert(errs["A"], HasLen, 3)
	c.Assert(errs["A"], HasError, validate.ErrMin)
	c.Assert(errs["A"], HasError, customErr1)
	c.Assert(errs["A"], HasError, customErr2)
	c.Assert(errs["B"], HasLen, 1)
	c.Assert(errs["B"], HasError, customErr1)
	c.Assert(errs["C"], HasLen, 1)
	c.Assert(errs["C"], HasError, validate.ErrMin)
	c.Assert(errs["D"], HasLen, 1)
	c.Assert(errs["D"], HasError, customErr2)

	errs = validate.ValidateAll(&test).(validate.Errors)
	c.Assert(errs, HasLen, 4)
	c.Assert(errs["A"], HasLen, 3)
	c.Assert(errs["A"], HasError, validate.ErrMin)
	c.Assert(errs["A"], HasError, customErr1)
	c.Assert(errs["A"], HasError, customErr2)
	c.Assert(errs["B"], HasLen, 1)
	c.Assert(errs["B"], HasError, customErr1)
	c.Assert(errs["C"], HasLen, 1)
	c.Assert(errs["C"], HasError, validate.ErrMin)
	c.Assert(errs["D"], HasLen, 1)
	c.Assert(errs["D"], HasError, customErr2)
}

func (vs *ValidatorSuite) TestValidMap(c *C) {
	m := make(map[string]string)

	err := validate.ValidAll(m, "required")
	c.Assert(err, NotNil)
	errs, ok := err.(validate.ErrorList)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrRequired)

	err = validate.ValidAll(m, "min(1)")
	c.Assert(err, NotNil)
	errs, ok = err.(validate.ErrorList)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMin)

	m = map[string]string{"A": "a", "B": "a"}
	err = validate.ValidAll(m, "max(1)")
	c.Assert(err, NotNil)
	errs, ok = err.(validate.ErrorList)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMax)

	err = validate.ValidAll(m, "min(2), max(5)")
	c.Assert(err, IsNil)

	m = map[string]string{
		"1": "a",
		"2": "b",
		"3": "c",
		"4": "d",
		"5": "e",
	}
	err = validate.ValidAll(m, "len(4),min(6),max(1),required")
	c.Assert(err, NotNil)
	errs, ok = err.(validate.ErrorList)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrLen)
	c.Assert(errs, HasError, validate.ErrMin)
	c.Assert(errs, HasError, validate.ErrMax)
	c.Assert(errs, Not(HasError), validate.ErrRequired)

}

func (vs *ValidatorSuite) TestValidFloat(c *C) {
	err := validate.ValidAll(12.34, "required")
	c.Assert(err, IsNil)

	err = validate.ValidAll(0.0, "required")
	c.Assert(err, NotNil)
	errs, ok := err.(validate.ErrorList)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrRequired)
}

func (vs *ValidatorSuite) TestValidInt(c *C) {
	i := 123
	err := validate.ValidAll(i, "required")
	c.Assert(err, IsNil)

	err = validate.ValidAll(i, "min(1)")
	c.Assert(err, IsNil)

	err = validate.ValidAll(i, "min(124), max(122)")
	c.Assert(err, NotNil)
	errs, ok := err.(validate.ErrorList)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMin)
	c.Assert(errs, HasError, validate.ErrMax)

	err = validate.ValidAll(i, "max(10)")
	c.Assert(err, NotNil)
	errs, ok = err.(validate.ErrorList)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMax)
}

func (vs *ValidatorSuite) TestValidString(c *C) {
	s := "test1234"
	err := validate.ValidAll(s, "len(8)")
	c.Assert(err, IsNil)

	err = validate.ValidAll(s, "len(0)")
	c.Assert(err, NotNil)
	errs, ok := err.(validate.ErrorList)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrLen)

	err = validate.ValidAll(s, `regexp("^[tes]{4}.*")`)
	c.Assert(err, IsNil)

	err = validate.ValidAll(s, `regexp("^.*[0-9]{5}$")`)
	c.Assert(errs, NotNil)

	err = validate.ValidAll("", `required,len(3),max(1)`)
	c.Assert(err, NotNil)
	errs, ok = err.(validate.ErrorList)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs, HasError, validate.ErrRequired)
	c.Assert(errs, HasError, validate.ErrLen)
	c.Assert(errs, Not(HasError), validate.ErrMax)
}

func (vs *ValidatorSuite) TestValidPtr(c *C) {
	s := "test1234"
	err := validate.ValidAll(&s, "len(8)")
	c.Assert(err, IsNil)

	err = validate.ValidAll(&s, "len(0)")
	c.Assert(err, NotNil)
	errs, ok := err.(validate.ErrorList)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrLen)
}

func (vs *ValidatorSuite) TestValidateAllWithCustomValidator(c *C) {

	validator := validate.NewValidator()

	err := validator.SetValidationFunc("equals", func(val interface{}, params []string) error {
		v := val.(string)
		if v != params[0] {
			return validate.ErrInvalid
		}
		return nil
	})
	c.Assert(err, IsNil)

	err = validator.ValidAll("bar", `equals("bar")`)
	c.Assert(err, IsNil)

	err = validator.ValidAll("foo", `equals("bar")`)
	c.Assert(err, NotNil)
	errs, ok := err.(validate.ErrorList)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrInvalid)
}

func (vs *ValidatorSuite) TestValidateAllUnsetValidator(c *C) {

	validator := validate.NewValidator()

	err := validator.SetValidationFunc("required", nil)
	c.Assert(err, IsNil)

	err = validator.ValidAll("foo", `required`)
	c.Assert(err, Equals, validate.ErrUnknownTag)
}

func (vs *ValidatorSuite) TestValidateAllSetValidatorWithEmptyNameShouldError(c *C) {
	validator := validate.NewValidator()

	err := validator.SetValidationFunc("", func(_ interface{}, _ []string) error {
		return nil
	})

	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "name cannot be empty")
}

func (vs *ValidatorSuite) TestValidateAllWithTag(c *C) {
	test := struct {
		A string `testvalidate:"required" validate:"required,min(10)"`
	}{}

	validator := validate.WithTag(`testvalidate`)

	err := validator.ValidateAll(test)
	c.Assert(err, NotNil)

	errs, ok := err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs["A"], HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
}

func (vs *ValidatorSuite) TestValidateAllSetTag(c *C) {
	test := struct {
		A string `testvalidate:"required" validate:"required,min(10)"`
	}{}

	validator := validate.NewValidator()
	validator.SetTag(`testvalidate`)

	err := validator.ValidateAll(test)
	c.Assert(err, NotNil)

	errs, ok := err.(validate.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs["A"], HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
}

type testSimpleValidateEmbed struct {
	A int `validate:"min(10)"`
}

func (t *testSimpleValidateEmbed) Validate() validate.Errors {
	return nil
}

type hasErrorChecker struct {
	*CheckerInfo
}

func (c *hasErrorChecker) Check(params []interface{}, names []string) (bool, string) {
	var (
		ok    bool
		slice []error
		value error
	)

	if e, ok := params[0].([]error); ok {
		slice = e
	} else if e, ok := params[0].(validate.ErrorList); ok {
		slice = e
	} else {
		return false, "First parameter is not an Errors"
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
