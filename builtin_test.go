package validate_test

import (
	"github.com/mbict/go-validate"
	. "gopkg.in/check.v1"
)

type BuiltinSuite struct{}

var _ = Suite(&BuiltinSuite{})

func (s *BuiltinSuite) TestRequiredFail(c *C) {
	c.Check(validate.Valid(string(""), "required"), ErrorMatches, "required")
	c.Check(validate.Valid([]int{}, "required"), ErrorMatches, "required")
	c.Check(validate.Valid(map[int]int{}, "required"), ErrorMatches, "required")
	c.Check(validate.Valid(int(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(int8(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(int16(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(int32(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(int64(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(uint(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(uint8(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(uint16(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(uint32(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(uint64(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(uintptr(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(float32(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(float64(0), "required"), ErrorMatches, "required")
	c.Check(validate.Valid(bool(false), "required"), ErrorMatches, "required")
}

func (s *BuiltinSuite) TestRequiredPass(c *C) {

	c.Check(validate.Valid(string("abc"), "required"), IsNil)
	c.Check(validate.Valid([]int{1, 2, 3}, "required"), IsNil)
	c.Check(validate.Valid(map[int]int{1: 1, 2: 2, 3: 3}, "required"), IsNil)
	c.Check(validate.Valid(int(-1), "required"), IsNil)
	c.Check(validate.Valid(int8(-1), "required"), IsNil)
	c.Check(validate.Valid(int16(-1), "required"), IsNil)
	c.Check(validate.Valid(int32(-1), "required"), IsNil)
	c.Check(validate.Valid(int64(-1), "required"), IsNil)
	c.Check(validate.Valid(int(1), "required"), IsNil)
	c.Check(validate.Valid(int8(1), "required"), IsNil)
	c.Check(validate.Valid(int16(1), "required"), IsNil)
	c.Check(validate.Valid(int32(1), "required"), IsNil)
	c.Check(validate.Valid(int64(1), "required"), IsNil)
	c.Check(validate.Valid(uint(1), "required"), IsNil)
	c.Check(validate.Valid(uint8(1), "required"), IsNil)
	c.Check(validate.Valid(uint16(1), "required"), IsNil)
	c.Check(validate.Valid(uint32(1), "required"), IsNil)
	c.Check(validate.Valid(uint64(1), "required"), IsNil)
	c.Check(validate.Valid(uintptr(1), "required"), IsNil)
	c.Check(validate.Valid(float32(1.1), "required"), IsNil)
	c.Check(validate.Valid(float64(1.1), "required"), IsNil)
	c.Check(validate.Valid(bool(true), "required"), IsNil)
}
