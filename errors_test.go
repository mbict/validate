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

func (vs *ErrorsSuite) TestJsonMarshal(c *C) {
	var errors validate.Errors

	errors.Add("foo", validate.ErrRequired)
	errors.Add("bar", validate.ErrEmpty, validate.ErrEmail)

	res, err := json.Marshal(errors)

	c.Assert(err, IsNil)
	c.Assert(string(res), Equals, `{"bar":["value is empty","invalid email"],"foo":["required"]}`)
}

func (vs *ErrorsSuite) TestJsonUnmarshal(c *C) {
	var errors validate.Errors

	err := json.Unmarshal([]byte(`{"bar":["value is empty","invalid email"],"foo":["required"]}`), &errors)

	c.Assert(err, IsNil)

	v, ok := errors["bar"]
	c.Assert(ok, Equals, true)
	c.Assert(v, HasLen, 2)
	c.Assert(v[0].Error(), Equals, "value is empty")
	c.Assert(v[1].Error(), Equals, "invalid email")

	v, ok = errors["foo"]
	c.Assert(ok, Equals, true)
	c.Assert(v, HasLen, 1)
	c.Assert(v[0].Error(), Equals, "required")
}
