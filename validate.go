package validate

import (
	"fmt"
	"reflect"
)

type Validator interface {
	Validate(Errors) Errors
}

func NewValidator() *Validate {
	return &Validate{cache: newCache()}
}

type Validate struct {
	cache *cache
}

func (val *Validate) Validate(dst interface{}) (errors Errors) {
	v := reflect.Indirect(reflect.ValueOf(dst))
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			if validateErrors := val.validate(v.Index(i), fmt.Sprintf("%d.", i)); len(validateErrors) > 0 {
				errors = append(errors, validateErrors...)
			}
		}
		return
	} else if v.Kind() == reflect.Struct {
		return val.validate(v, "")
	}
	errors.Add([]string{}, "General", "Cannot validate unkown structure")
	return
}

func (val *Validate) validate(v reflect.Value, path string) Errors {
	var errors Errors
	v = reflect.Indirect(v) //v.Elem()
	t := v.Type()
	s := val.cache.get(t)
	for _, field := range s.fields {
		value := v.FieldByIndex(field.index)
		fieldPath := path + field.name

		//run validators
		for _, validator := range field.validators {
			if validateErrors := validator.validator(value, validator.params, fieldPath); len(validateErrors) > 0 {
				errors = append(errors, validateErrors...)
				break
			}
		}

		// Validate nested and embedded structs (if pointer, only do so if not nil)
		if value.Kind() == reflect.Struct || (value.Kind() == reflect.Ptr && !value.IsNil() && value.Elem().Kind() == reflect.Struct) {
			//single struct
			if validateErrors := val.validate(value, fieldPath+"."); len(validateErrors) > 0 {
				errors = append(errors, validateErrors...)
			}
		} else if field.ss == true && (value.Kind() == reflect.Slice || value.Kind() == reflect.Array) {
			//slice and arrays
			for i := 0; i < value.Len(); i++ {
				if validateErrors := val.validate(value.Index(i), fmt.Sprintf("%s.%d.", fieldPath, i)); len(validateErrors) > 0 {
					errors = append(errors, validateErrors...)
				}
			}
		}
	}

	if valdiatorFunc, ok := v.Interface().(Validator); ok {
		errors = valdiatorFunc.Validate(errors)
	}
	return errors
}
