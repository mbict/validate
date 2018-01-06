package validate

import (
	"reflect"
	"strings"
	"unicode"
)

type NameResolverFunc func(reflect.StructField) string

var DefaultNameResolver = func(field reflect.StructField) string {
	return field.Name
}

var JsonNameResolver = FallbackNameResolver(TagNameResolver("json"))

var JsonNameSnakeCaseResolver = SnakeCaseResolver(JsonNameResolver)

func FallbackNameResolver(resolvers ...NameResolverFunc) NameResolverFunc {
	return func(field reflect.StructField) string {
		for _, resolver := range resolvers {
			if val := resolver(field); val != "" {
				return val
			}
		}
		return DefaultNameResolver(field)
	}
}

func TagNameResolver(tag string) NameResolverFunc {
	return func(field reflect.StructField) string {
		tag := field.Tag.Get(tag)
		props := strings.SplitN(tag, ",", 2)
		return props[0]
	}
}

func SnakeCaseResolver(resolver NameResolverFunc) NameResolverFunc {
	return func(field reflect.StructField) string {
		return toSnakeCase(resolver(field))
	}
}

func toSnakeCase(in string) string {
	runes := []rune(in)

	var out []rune
	for i := 0; i < len(runes); i++ {
		if i > 0 && (unicode.IsUpper(runes[i]) || unicode.IsNumber(runes[i])) && ((i+1 < len(runes) && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
