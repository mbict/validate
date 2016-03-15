package validate_test

import (
	"github.com/mbict/go-validate"
	. "gopkg.in/check.v1"
)

type ErrorSuite struct{}

var _ = Suite(&ErrorSuite{})

func (ms *ErrorSuite) TestErrorsToString(c *C) {
	errs := validate.Errors{
		validate.ErrBadParameter,
		validate.ErrLen,
		validate.ErrMax,
	}

	c.Assert(errs, ErrorMatches, "bad parameter, invalid length, greater than max")
}

func (ms *ErrorSuite) TestNoErrorsToEmptyString(c *C) {
	errs := validate.Errors{}

	c.Assert(errs, ErrorMatches, "")
}

func (ms *ErrorSuite) TestErrorMapToString(c *C) {
	errs := make(validate.ErrorMap)
	errs["A"] = validate.Errors{validate.ErrBadParameter, validate.ErrLen}
	errs["B"] = validate.Errors{validate.ErrMax}

	c.Assert(errs, ErrorMatches, ".*B:\\[greater than max\\].*")
	c.Assert(errs, ErrorMatches, ".*A:\\[bad parameter, invalid length\\].*")
}

func (ms *ErrorSuite) TestErrorMapWithNoErrorsToEmptyString(c *C) {
	errs := make(validate.ErrorMap)

	c.Assert(errs, ErrorMatches, "")
}

func (ms *ErrorSuite) TestHasErrors(c *C) {
	errs := make(validate.ErrorMap)
	errs["A"] = validate.Errors{validate.ErrBadParameter, validate.ErrLen}

	c.Assert(validate.HasErrors(errs, "B"), Equals, false)
	c.Assert(validate.HasErrors(errs, "A"), Equals, true)
}

func (ms *ErrorSuite) TestHasError(c *C) {
	errs := make(validate.ErrorMap)
	errs["A"] = validate.Errors{validate.ErrLen}

	c.Assert(validate.HasError(errs, "B", "invalid length"), Equals, false)
	c.Assert(validate.HasError(errs, "A", "bad parameter"), Equals, false)
	c.Assert(validate.HasError(errs, "A", "invalid length"), Equals, true)
}
