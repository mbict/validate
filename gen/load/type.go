package load

import "go/types"

type ValidateStruct struct {
	Name string

	pkg *types.Package
	st  *types.Struct
}

type Field struct {
	Name string
	Type types.Type
}
