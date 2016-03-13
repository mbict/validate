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
