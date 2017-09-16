[![wercker status](https://app.wercker.com/status/3a36cd7798b739f402718a0ba24334c4/s "wercker status")](https://app.wercker.com/project/bykey/3a36cd7798b739f402718a0ba24334c4)
[![Build Status](https://travis-ci.org/mbict/go-validate.png?branch=master)](https://travis-ci.org/mbict/go-validate)
[![GoDoc](https://godoc.org/github.com/mbict/go-validate?status.png)](http://godoc.org/github.com/mbict/go-validate)
[![GoCover](http://gocover.io/_badge/github.com/mbict/go-validate)](http://gocover.io/github.com/mbict/go-validate)
[![GoReportCard](http://goreportcard.com/badge/mbict/go-validate)](http://goreportcard.com/report/mbict/go-validate)

Validate
========

Validate provides a simple way to validate the contents of variables and structures

Installation
============

Use go get for installing.

	go get github.com/mbict/go-validate

And then import the package into your own code.

	import (
		"github.com/mbict/go-validate"
	)

Usage
=====

A simple example would.

	type YourStruct struct {
		Username string `validate:"min(3);max(40);regexp(\"^[a-zA-Z]*$\")"`
		Page string     `validate:"required"`
		Age int         `validate:"min(21)"`
		Tags []string   `validate:"min(1);max(4)"`
	}

	test := YourStruct{Username: "something", Age: 20}
	if errs := validate.Validate(test); errs != nil {
		// values not valid, deal with errors here
		
		errsMap := errs.(validate.ErrorMap)
		
		errsMap["Page"].Error()
	}

Builtin validators
==================

There are a few common set of builin validators included in the package.

	len
		For numeric numbers, max will simply make sure that the
		value is equal to the parameter given. For strings, it
		checks that the string length is exactly that number of
		characters. For slices,	arrays, and maps, validates the
		number of items. Usage: len(10)
	
	max
		For numeric numbers, max will simply make sure that the
		value is lesser or equal to the parameter given. For strings,
		it checks that the string length is at most that number of
		characters. For slices,	arrays, and maps, validates the
		number of items. Usage: max(10)
	
	min
		For numeric numbers, min will simply make sure that the value
		is greater or equal to the parameter given. For strings, it
		checks that the string length is at least that number of
		characters. For slices, arrays, and maps, validates the
		number of items. Usage: min(10)
	
	required
		This validates that the value is not zero. The appropriate
		zero value is given by the Go spec (e.g. for int it's 0, for
		string it's "", for pointers is nil, etc.) For structs, it
		will not check to see if the struct itself has all zero
		values, instead use a pointer or put nonzero on the struct's
		keys that you care about. Usage: required
	
	regexp
		Only valid for string types, it will validate that the
		value matches the regular expression provided as parameter.
		Usage: regexp("^a.*b$")

    email
		Only valid for string types, it will validate that the
		value matches a valid email pattern. Usage: email
		
	url
        Only valid for string types, it will validate that the
        value matches a valid url pattern. Usage: url
        
    between
        tbd
    
    around
        tbd
        
    include
        tbd
        
    exclude
        tbd

    number
        tbd
        
    numeric
        tbd
        
    alpha
        tbd
        
    alphanumeric
        tbd

    alpha_dash
        tbd
        
    alpha_dash_dot
        tbd
        
    uuid
        tbd
        
    uuid3
        tbd
        
    uuid4
        tbd
        
    uuid5
        tbd
        
    base64
        tbd
        
Custom validators

It is possible to define your own custom validators by using SetValidationFunc.
You needs to create a validation function that is follows the definition 
of ValidatorFunc.

	// Custom validator
	func notValue(v interface{}, params []string) error {
		st := reflect.ValueOf(v)
		if st.Kind() != reflect.String {
			return errors.New("notValue only validates strings")
		}
		
		if len(params) < 1 {
		return errors.New("notValue not enough params")
		
		if st.String() == param[0] {
			return errors.New("value cannot be "+param[0])
		}
		return nil
	}

Then you need to add it to the list of validators and give it a "tag"
name.

	validate.SetValidationFunc("notValue", notValue)

From this point on you can use the notValue validation tag in your 
structure tags

	type T struct {
		A string  `validate:"required,notValue(abc)"`
	}
	t := T{"abc"}
	if errs := validate.Validate(t); errs != nil {
		fmt.Printf("Field A error: %s\n", errs["A"][0])
	}

You can also have multiple sets of validator rules with SetTag().

	type T struct {
		A int `foo:"required" bar:"min(10)"`
	}
	t := T{5}
	validate.SetTag("foo")
	validate.Validate(t) // valid as it's required
	validate.SetTag("bar")
	validate.Validate(t) // invalid as it's less than 10

SetTag is probably better used with multiple validators.

	fooValidator := validate.NewValidator()
	fooValidator.SetTag("foo")
	barValidator := validate.NewValidator()
	barValidator.SetTag("bar")
	fooValidator.Validate(t)
	barValidator.Validate(t)

This keeps the default validator's tag clean.


Structure custom validation
===========================
Your structure maybe needs a custom validation that cannot be solved with the builtin or custom validator.
Therefor if the structure implements the ValidateInterface the method validate is called if all validations are done

```go
var CustomErr := errors.New("custom error")

type YourStruct struct {
	...
}

func (o YourStruct) Validate() validate.ErrorMap {
	// in here you can do some custom validation
	// when you got a validation error, init and add a entry to the error map.
	return validate.ErrorMap{"A": validate.Errors{CustomErr,}}
}

```

Dependencies
============
go-validate requires go-tags for parsing the structure tags into something useful