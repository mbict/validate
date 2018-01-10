package validate

import (
	"fmt"
	"reflect"
)

type ValidateFunc func(i interface{}) error

type structRules map[reflect.Type]rules

type rules []rule

type rule struct {
	Name       string
	FieldIndex int
	IsSlice    bool
	IsStruct   bool
	Validators []validatorTag
	Subset     rules
}

func (r *rules) Validate(value reflect.Value, stopOnError bool) Errors {
	var errs Errors
	for _, rule := range *r {
		v := value.Field(rule.FieldIndex)
		if verr := rule.Validate(v, stopOnError); verr != nil {
			errs.Merge(verr)
		}
	}

	// implemented the ValidateInterface
	if validateFunc, ok := value.Interface().(ValidateInterface); ok {
		if err := validateFunc.Validate(); err != nil {
			errs.Merge(err)
		}
	}

	if errs == nil {
		return nil
	}
	return errs
}

func (r *rule) Validate(value reflect.Value, stopOnError bool) Errors {
	var errs Errors

	i := value.Interface()
	for _, validator := range r.Validators {
		if err := validator.Fn(i, validator.Args); err != nil {
			errs.Add(r.Name, err)

			if stopOnError == true {
				return errs
			}
		}
	}

	if value.Kind() == reflect.Ptr && value.IsNil() {
		return errs
	}

	value = reflect.Indirect(value)
	if r.IsSlice && r.IsStruct {
		for i := 0; i < value.Len(); i++ {
			errv := r.Subset.Validate(reflect.Indirect(value.Index(i)), stopOnError)
			if errv != nil {
				errs.MergePrefix(fmt.Sprintf("%s.%d.", r.Name, i), errv)
			}
		}
	} else if r.IsStruct {
		errv := r.Subset.Validate(value, stopOnError)
		if errv != nil {
			errs.MergePrefix(r.Name+".", errv)
		}
	}

	if errs == nil {
		return nil
	}
	return errs
}
