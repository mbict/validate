package validate_test

import (
	errors "github.com/mbict/go-errors"
	validate "github.com/mbict/go-validate"
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

func (vs *ValidatorSuite) TestValidate(c *C) {
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
	err := validate.Validate(test[0])
	c.Assert(err, NotNil)

	errs, ok := err.(errors.ErrorHash)
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
	err = validate.Validate(test[1])
	c.Assert(err, NotNil)

	errs, ok = err.(errors.ErrorHash)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrMax)
	c.Assert(errs["B"], HasLen, 1)
	c.Assert(errs["B"], HasError, validate.ErrMax)
	c.Assert(errs["C"], HasLen, 1)
	c.Assert(errs["C"], HasError, validate.ErrRequired)
	c.Assert(errs["D"], HasLen, 0)
	c.Assert(errs["E"], HasLen, 1)
	c.Assert(errs["E"], HasError, validate.ErrMax)

	//should pass
	err = validate.Validate(test[2])
	c.Assert(err, IsNil)
}

func (vs *ValidatorSuite) TestValidateEmbedStruct(c *C) {
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
	err := validate.Validate(test[0])
	c.Assert(err, NotNil)

	errs, ok := err.(errors.ErrorHash)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A.A"], HasLen, 1)
	c.Assert(errs["A.A"], HasError, validate.ErrMin)
	c.Assert(errs["B"], HasLen, 1)
	c.Assert(errs["B"], HasError, validate.ErrRequired)
	c.Assert(errs["B.A"], HasLen, 0)

	//error, initialized ptr
	err = validate.Validate(test[1])
	c.Assert(err, NotNil)

	errs, ok = err.(errors.ErrorHash)
	c.Assert(ok, Equals, true)
	c.Assert(errs["B"], HasLen, 0)
	c.Assert(errs["B.A"], HasLen, 1)
	c.Assert(errs["B.A"], HasError, validate.ErrMin)

	//error, initialized ptr
	err = validate.Validate(test[1])
	c.Assert(err, NotNil)

	errs, ok = err.(errors.ErrorHash)
	c.Assert(ok, Equals, true)
	c.Assert(errs["B"], HasLen, 0)
	c.Assert(errs["B.A"], HasLen, 1)
	c.Assert(errs["B.A"], HasError, validate.ErrMin)

	//should pass
	err = validate.Validate(test[2])
	c.Assert(err, IsNil)
}

func (vs *ValidatorSuite) TestValidateSlice(c *C) {
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
	err := validate.Validate(test[0])
	c.Assert(err, NotNil)

	errs, ok := err.(errors.ErrorHash)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasLen, 2)
	c.Assert(errs["A"], HasError, validate.ErrMin)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
	c.Assert(errs["B"], HasLen, 1)
	c.Assert(errs["B"], HasError, validate.ErrRequired)

	//error, zero length slices
	err = validate.Validate(test[1])
	c.Assert(err, NotNil)

	errs, ok = err.(errors.ErrorHash)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasLen, 2)
	c.Assert(errs["A"], HasError, validate.ErrMin)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
	c.Assert(errs["B"], HasLen, 1)
	c.Assert(errs["B"], HasError, validate.ErrRequired)

	//error, filled slices
	err = validate.Validate(test[2])
	c.Assert(err, NotNil)

	errs, ok = err.(errors.ErrorHash)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrMin)
	c.Assert(errs["A.0.A"], HasLen, 1)
	c.Assert(errs["A.0.A"], HasError, validate.ErrMin)
	c.Assert(errs["B"], HasLen, 0)
	c.Assert(errs["B.0.A"], HasLen, 0)
	c.Assert(errs["B.1.A"], HasLen, 1)
	c.Assert(errs["B.1.A"], HasError, validate.ErrMin)

	//should pass
	err = validate.Validate(test[3])
	c.Assert(err, IsNil)
}

func (vs *ValidatorSuite) TestValidateIgnoreNonExportedVars(c *C) {
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
	err := validate.Validate(test[0])
	c.Assert(err, NotNil)

	errs, ok := err.(errors.ErrorHash)
	c.Assert(ok, Equals, true)
	c.Assert(errs["A"], HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
	c.Assert(errs["b"], HasLen, 0)

	//error, initialized ptr
	err = validate.Validate(test[1])
	c.Assert(err, IsNil)
}

func (vs *ValidatorSuite) TestValidateInputNilValue(c *C) {
	//nil pointer
	err := validate.Validate(nil)
	c.Assert(err, Equals, validate.ErrUnsupported)

	//typed nil pointer
	err = validate.Validate((*testSimple)(nil))
	c.Assert(err, Equals, validate.ErrUnsupported)
}

func (vs *ValidatorSuite) TestValidErrorUnkownTag(c *C) {
	err := validate.Valid(1, "min(10);nonexistingvalidator(1,2,3)")

	c.Assert(err, NotNil)
	c.Assert(err, Equals, validate.ErrUnknownTag)
}

func (vs *ValidatorSuite) TestValidErrorSyntax(c *C) {
	err := validate.Valid(1, "min(10)|bv")

	c.Assert(err, NotNil)
	c.Assert(err, Equals, validate.ErrSyntax)
}

func (vs *ValidatorSuite) TestValidIgnoreTagReturnNil(c *C) {
	err := validate.Valid(1, "-")

	c.Assert(err, IsNil)
}

func (vs *ValidatorSuite) TestValidateIgnoreTag(c *C) {
	test := struct {
		A testSimple `validate:"-"`
	}{}

	err := validate.Validate(test)

	c.Assert(err, IsNil)
}

type testStructValidateFuncInterface struct {
	A            int `validate:"min(2)"`
	B            int `validate:"min(2)"`
	C            int `validate:"min(10)"`
	validateFunc func() errors.ErrorHash
}

func (s testStructValidateFuncInterface) Validate() errors.ErrorHash {
	return s.validateFunc()
}

func (vs *ValidatorSuite) TestStructValidateInterface(c *C) {
	customErr1 := errors.New("custom 1")
	customErr2 := errors.New("cutsom 2")

	test := testStructValidateFuncInterface{
		A: 1,
		B: 2,
		C: 3,
		validateFunc: func() errors.ErrorHash {
			return errors.ErrorHash{
				"A": errors.Errors{customErr1, customErr2},
				"B": errors.Errors{customErr1},
				"D": errors.Errors{customErr2},
			}
		},
	}

	errs := validate.Validate(test).(errors.ErrorHash)

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

	errs = validate.Validate(&test).(errors.ErrorHash)

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

	err := validate.Valid(m, "required")
	c.Assert(err, NotNil)
	errs, ok := err.(errors.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrRequired)

	err = validate.Valid(m, "min(1)")
	c.Assert(err, NotNil)
	errs, ok = err.(errors.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMin)

	m = map[string]string{"A": "a", "B": "a"}
	err = validate.Valid(m, "max(1)")
	c.Assert(err, NotNil)
	errs, ok = err.(errors.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMax)

	err = validate.Valid(m, "min(2), max(5)")
	c.Assert(err, IsNil)

	m = map[string]string{
		"1": "a",
		"2": "b",
		"3": "c",
		"4": "d",
		"5": "e",
	}
	err = validate.Valid(m, "len(4),min(6),max(1),required")
	c.Assert(err, NotNil)
	errs, ok = err.(errors.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrLen)
	c.Assert(errs, HasError, validate.ErrMin)
	c.Assert(errs, HasError, validate.ErrMax)
	c.Assert(errs, Not(HasError), validate.ErrRequired)

}

func (vs *ValidatorSuite) TestValidFloat(c *C) {
	err := validate.Valid(12.34, "required")
	c.Assert(err, IsNil)

	err = validate.Valid(0.0, "required")
	c.Assert(err, NotNil)
	errs, ok := err.(errors.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrRequired)
}

func (vs *ValidatorSuite) TestValidInt(c *C) {
	i := 123
	err := validate.Valid(i, "required")
	c.Assert(err, IsNil)

	err = validate.Valid(i, "min(1)")
	c.Assert(err, IsNil)

	err = validate.Valid(i, "min(124), max(122)")
	c.Assert(err, NotNil)
	errs, ok := err.(errors.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMin)
	c.Assert(errs, HasError, validate.ErrMax)

	err = validate.Valid(i, "max(10)")
	c.Assert(err, NotNil)
	errs, ok = err.(errors.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrMax)
}

func (vs *ValidatorSuite) TestValidString(c *C) {
	s := "test1234"
	err := validate.Valid(s, "len(8)")
	c.Assert(err, IsNil)

	err = validate.Valid(s, "len(0)")
	c.Assert(err, NotNil)
	errs, ok := err.(errors.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrLen)

	err = validate.Valid(s, `regexp("^[tes]{4}.*")`)
	c.Assert(err, IsNil)

	err = validate.Valid(s, `regexp("^.*[0-9]{5}$")`)
	c.Assert(errs, NotNil)

	err = validate.Valid("", `required,len(3),max(1)`)
	c.Assert(err, NotNil)
	errs, ok = err.(errors.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 2)
	c.Assert(errs, HasError, validate.ErrRequired)
	c.Assert(errs, HasError, validate.ErrLen)
	c.Assert(errs, Not(HasError), validate.ErrMax)
}

func (vs *ValidatorSuite) TestValidPtr(c *C) {
	s := "test1234"
	err := validate.Valid(&s, "len(8)")
	c.Assert(err, IsNil)

	err = validate.Valid(&s, "len(0)")
	c.Assert(err, NotNil)
	errs, ok := err.(errors.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrLen)
}

func (vs *ValidatorSuite) TestValidateWithCustomValidator(c *C) {

	validator := validate.NewValidator()

	err := validator.SetValidationFunc("equals", func(val interface{}, params []string) error {
		v := val.(string)
		if v != params[0] {
			return validate.ErrInvalid
		}
		return nil
	})
	c.Assert(err, IsNil)

	err = validator.Valid("bar", `equals("bar")`)
	c.Assert(err, IsNil)

	err = validator.Valid("foo", `equals("bar")`)
	c.Assert(err, NotNil)
	errs, ok := err.(errors.Errors)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasError, validate.ErrInvalid)
}

func (vs *ValidatorSuite) TestValidateUnsetValidator(c *C) {

	validator := validate.NewValidator()

	err := validator.SetValidationFunc("required", nil)
	c.Assert(err, IsNil)

	err = validator.Valid("foo", `required`)
	c.Assert(err, Equals, validate.ErrUnknownTag)
}

func (vs *ValidatorSuite) TestValidateSetValidatorWithEmptyNameShouldError(c *C) {
	validator := validate.NewValidator()

	err := validator.SetValidationFunc("", func(_ interface{}, _ []string) error {
		return nil
	})

	c.Assert(err, NotNil)
	c.Assert(err, ErrorMatches, "name cannot be empty")
}

func (vs *ValidatorSuite) TestValidateWithTag(c *C) {
	test := struct {
		A string `testvalidate:"required" validate:"required,min(10)"`
	}{}

	validator := validate.WithTag(`testvalidate`)

	err := validator.Validate(test)
	c.Assert(err, NotNil)

	errs, ok := err.(errors.ErrorHash)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs["A"], HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
}

func (vs *ValidatorSuite) TestValidateSetTag(c *C) {
	test := struct {
		A string `testvalidate:"required" validate:"required,min(10)"`
	}{}

	validator := validate.NewValidator()
	validator.SetTag(`testvalidate`)

	err := validator.Validate(test)
	c.Assert(err, NotNil)

	errs, ok := err.(errors.ErrorHash)
	c.Assert(ok, Equals, true)
	c.Assert(errs, HasLen, 1)
	c.Assert(errs["A"], HasLen, 1)
	c.Assert(errs["A"], HasError, validate.ErrRequired)
}

type testSimpleValidateEmbed struct {
	A int `validate:"min(10)"`
}

func (t *testSimpleValidateEmbed) Validate() errors.ErrorHash {
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
	slice, ok = params[0].(errors.Errors)
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
