package validate

import (
	"errors"
	"mime/multipart"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

var invalidPath = errors.New("schema: invalid path")

// newCache returns a new cache.
func newCache() *cache {
	c := cache{
		m:          make(map[reflect.Type]*structInfo),
		validators: make(map[string]ValidatorFunc),
		tag:        "validate",
	}
	for k, v := range validators {
		c.validators[k] = v
	}
	return &c
}

// cache caches meta-data about a struct.
type cache struct {
	l          sync.RWMutex
	m          map[reflect.Type]*structInfo
	validators map[string]ValidatorFunc
	tag        string
}

// get returns a cached structInfo, creating it if necessary.
func (c *cache) get(t reflect.Type) *structInfo {
	c.l.RLock()
	info := c.m[t]
	c.l.RUnlock()
	if info == nil {
		info = c.create(t, nil, nil)
		c.l.Lock()
		c.m[t] = info
		c.l.Unlock()
	}
	return info
}

var errorsType = reflect.TypeOf(Errors{})

// creat creates a structInfo with meta-data about a struct.
func (c *cache) create(t reflect.Type, info *structInfo, index []int) *structInfo {
	if info == nil {
		info = &structInfo{fields: []*fieldInfo{}}
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			ft := field.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				c.create(ft, info, append(index, field.Index...))
			}
		} else {
			c.createField(field, info, index)
		}
	}
	return info
}

var multipartFileheaderType = reflect.TypeOf(multipart.FileHeader{})

// createField creates a fieldInfo for the given field.
func (c *cache) createField(field reflect.StructField, info *structInfo, index []int) {
	//extract the field validators
	validators := c.fieldValidators(field, c.tag)

	// Check if the type is supported and don't cache it if not.
	// First let's get the basic type.
	isSlice, isStruct, isFileheader := false, false, false
	ft := field.Type
	if ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	if isSlice = ft.Kind() == reflect.Slice; isSlice {
		ft = ft.Elem()
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
	}

	//for multipart files we dont threat them as slice structures
	isFileheader = ft == multipartFileheaderType
	//_, isScanner := reflect.New(ft).Interface().(sql.Scanner)

	info.fields = append(info.fields, &fieldInfo{
		typ:   field.Type,
		name:  field.Name,
		index: append(index, field.Index...),
		ss:    isSlice && isStruct && !isFileheader,
		//scanner:    isScanner,
		validators: validators,
	})
}

// ----------------------------------------------------------------------------

type structInfo struct {
	fields          []*fieldInfo
	hasValidateFunc bool
	scanner         bool
}

//func (i *structInfo) get(name string) *fieldInfo {
//	for _, field := range i.fields {
//		if strings.EqualFold(field.name, alias) {
//			return field
//		}
//	}
//	return nil
//}

type fieldInfo struct {
	typ        reflect.Type
	name       string // field name in the struct.
	index      []int
	ss         bool // true if this is a slice of structs.
	validators []validatorInfo
}

type validatorInfo struct {
	validator ValidatorFunc //validator function
	params    []string      //parsed tag params
}

// ----------------------------------------------------------------------------

var (
	reValidatorTags   = regexp.MustCompile("[0-9A-Za-z_]+(\\([^\\)]*\\))?")
	reValidatorParams = regexp.MustCompile("[\t\n\v\f\r ,\\(\\)]+")
)

// fieldAlias parses a field tag to get the validators.
func (c *cache) fieldValidators(field reflect.StructField, tagName string) (vi []validatorInfo) {
	if tag := field.Tag.Get(tagName); tag != "" {
		for _, field := range reValidatorTags.FindAllString(tag, -1) {
			valParams := reValidatorParams.Split(strings.Trim(field, " ()"), -1)
			val := strings.ToLower(valParams[0])

			if validator, ok := c.validators[val]; ok {
				vi = append(vi, validatorInfo{
					validator: validator,
					params:    valParams[1:],
				})
			} else {
				//error unkown validator
				panic("Unkown validator " + tag)
			}
		}
	}
	return
}
