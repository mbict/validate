package validate_test

import (
	"encoding/json"
	"github.com/mbict/go-validate"
	. "gopkg.in/check.v1"
)

type ErrorsSuite struct{}

var _ = Suite(&ErrorsSuite{})

func (vs *ErrorsSuite) TestError(c *C) {
	var errors validate.Errors

	errors.Add("foo", validate.ErrRequired)
	errors.Add("bar", validate.ErrEmpty, validate.ErrEmail)

	c.Assert(errors.Error(), Equals, "bar: [value is empty, invalid email], foo: [required]")
}

func (vs *ErrorsSuite) TestJsonSerialize(c *C) {
	var errors validate.Errors

	errors.Add("foo", validate.ErrRequired)
	errors.Add("bar", validate.ErrEmpty, validate.ErrEmail)

	res, err := json.Marshal(errors)

	c.Assert(err, IsNil)
	c.Assert(string(res), Equals, `{"bar":["value is empty","invalid email"],"foo":["required"]}`)

}
