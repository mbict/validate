package load

import (
	"go/importer"
	"go/types"
	"golang.org/x/tools/go/loader"
	"log"
	"path"
	"strings"
)

type Loader struct {
	program *loader.Program
}

func NewLoader(packages ...string) *Loader {
	loadConfig := loader.Config{
		AllowErrors:         true,
		TypeCheckFuncBodies: func(_ string) bool { return false },
		TypeChecker: types.Config{
			DisableUnusedImportCheck: true,
			Importer:/*importer.Default()*/ importer.For("source", nil),
		},
	}

	for _, pkg := range packages {
		loadConfig.Import(pkg)
	}

	p, err := loadConfig.Load()
	if err != nil {
		log.Fatal(err)
	}

	return &Loader{
		program: p,
	}
}

func (l *Loader) LookupType(typeName string) types.Type {
	pkgPath, typeName := path.Split(typeName)
	for pkgName, pkg := range l.program.Imported {
		if strings.HasSuffix(pkgName, path.Clean(pkgPath)) {
			log.Printf("scanning package: %s", pkgName)
			for _, scope := range pkg.Scopes {
				if s := lookup(scope.Parent(), typeName); s != nil {
					log.Printf("found type %s in package: %s", typeName, pkgName)
					return s
				}
			}
		}
	}
	return nil
}

func (l *Loader) LookupInterface(typeName string) types.Type {
	if s := l.lookupCallback(typeName, toInterface); s != nil {
		return s.(*types.Interface)
	}
	return nil
}

func (l *Loader) LookupStruct(typeName string) *types.Struct {
	if s := l.lookupCallback(typeName, toStruct); s != nil {
		return s.(*types.Struct)
	}
	return nil
}

func (l *Loader) lookupCallback(typeName string, toFunc func(types.Type) types.Type) types.Type {
	pkgPath, typeName := path.Split(typeName)
	for pkgName, pkg := range l.program.Imported {
		if strings.HasSuffix(pkgName, path.Clean(pkgPath)) {
			log.Printf("scanning package: %s", pkgName)
			for _, scope := range pkg.Scopes {
				if s := lookupFunc(scope.Parent(), typeName, toFunc); s != nil {
					log.Printf("found type %s in package: %s", typeName, pkgName)
					return s
				}
			}
		}
	}
	return nil
}

func lookup(s *types.Scope, name string) types.Type {
	if o := s.Lookup(name); o != nil {
		return o.Type()
	}

	for i := 0; i < s.NumChildren(); i++ {
		if s := lookup(s.Child(i), name); s != nil {
			return s
		}
	}
	return nil
}

func toInterface(t types.Type) types.Type {
	if s, ok := t.Underlying().(*types.Interface); ok {
		return s
	}
	return nil
}

func toStruct(t types.Type) types.Type {
	if s, ok := t.Underlying().(*types.Struct); ok {
		return s
	}
	return nil
}

func lookupFunc(s *types.Scope, name string, toFunc func(types.Type) types.Type) types.Type {
	if o := s.Lookup(name); o != nil {
		if u := toFunc(o.Type()); u != nil {
			return u
		}
	}

	for i := 0; i < s.NumChildren(); i++ {
		s := lookupFunc(s.Child(i), name, toFunc)
		if s != nil {
			return s
		}
	}
	return nil
}
