package validate

import (
	"fmt"
	"reflect"
)

type ValidateFunc func(i interface{}) error

type StructRules map[reflect.Type]Rules

type Rules []Rule

type Rule struct {
	Name       string
	FieldIndex int
	IsSlice    bool
	IsStruct   bool
	Validators []validatorTag
	Subset     Rules
}

func (r *Rules) Validate(value reflect.Value) Errors {
	var errs Errors
	for _, rule := range *r {
		v := value.Field(rule.FieldIndex)
		if verr := rule.Validate(v); verr != nil {
			errs.Merge(verr)
		}
	}

	// implemented the ValidateInterface
	if validateFunc, ok := value.Interface().(ValidateInterface); ok {
		if err := validateFunc.Validate(); err != nil {
			errs.Merge(err)
		}
	}

	return errs
}

func (r *Rule) Validate(value reflect.Value) Errors {
	var errs Errors

	i := value.Interface()
	for _, validator := range r.Validators {
		if err := validator.Fn(i, validator.Args); err != nil {
			errs.Add(r.Name, err)
		}
	}

	if value.Kind() == reflect.Ptr && value.IsNil() {
		return errs
	}

	value = reflect.Indirect(value)
	if r.IsSlice && r.IsStruct {
		for i := 0; i < value.Len(); i++ {
			errv := r.Subset.Validate(reflect.Indirect(value.Index(i)))
			if errv != nil {
				errs.MergePrefix(fmt.Sprintf("%s.%d.", r.Name, i), errv)
			}
		}
	} else if r.IsStruct {
		errv := r.Subset.Validate(value)
		if errv != nil {
			errs.MergePrefix(r.Name+".", errv)
		}
	}

	return errs
}
