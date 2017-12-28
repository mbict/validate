package main

import (
	"flag"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/mbict/go-validate/gen/load"
	"go/build"
	"go/types"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	options struct {
		types stringSlice
	}
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(val string) error {
	*s = strings.Split(val, ",")
	return nil
}

func init() {
	flag.Var(&options.types, "type", `Type name.
        Might be of the form '<pkg>.<type>' or just '<type>'. Where:
        - <pkg> can be either package name (e.x 'github.com/repository/project/package')
          or relative path (e.x './' or '../package').
        - <type> should be the type name.`)
	flag.Parse()
}

func main() {

	loader := load.NewLoader(
		"github.com/mbict/go-validate",
		pkgPath(),
	)

	validateInterface := loader.LookupType("github.com/mbict/go-validate/ValidateInterface").Underlying().(*types.Interface)
	spew.Dump(validateInterface)

	for _, typeName := range options.types {
		st := loader.LookupStruct(path.Join(pkgPath(), typeName))

		for i := 0; i < st.NumFields(); i++ {
			spew.Dump(st.Field(i).Type())
		}

	}
	/*

		fmt.Println(pkgPath())

		//options.types.Set("createClientRequest")

		if len(options.types) == 0 {
			log.Fatal("Must give type full name")
		}

		var errors []string

		for _, typeName := range options.types {
			log.Printf("Loading type `%s`", typeName)

			loadConfig := loader.Config{
				AllowErrors:         true,
				TypeCheckFuncBodies: func(_ string) bool { return false },
				TypeChecker: types.Config{
					DisableUnusedImportCheck:        true,
					Importer:  importer.For("source", nil),
				},
			}

			//loadConfig.Import("auth/storm/example")
			//loadConfig.Import("auth/api/client")
			loadConfig.Import("github.com/mbict/go-validate")
			loadConfig.Import(pkgPath())

			p, err := loadConfig.Load()
			if err != nil {
				log.Fatal(err)
			}

			var validateInterface *types.Interface

			pkg := p.Imported["auth/vendor/github.com/mbict/go-validate"]
			for _, scope := range pkg.Scopes {
				s := scope.Parent().Lookup("ValidateInterface")
				if s != nil {
					var ok bool
					validateInterface, ok = s.Type().Underlying().(*types.Interface)
					if ok {
						spew.Dump(validateInterface)
					}
				}
			}

			for pkgName, pkg := range p.Imported {

				log.Printf("scanning package: %s", pkgName)
				for _, scope := range pkg.Scopes {
					s := lookup(scope.Parent(), typeName)

					if s != nil {

						ss := scope.Parent().Lookup(typeName).Type()
						spew.Dump(ss)
						//					//non pointer implementation
						//					fmt.Printf("type does implement the ValidateInterface (%t)\n", types.Implements(ss, validateInterface))

						//pointer variant
						fmt.Printf("type does implement the ValidateInterface (%t)\n", types.Implements(types.NewPointer(ss), validateInterface))

						for _, t := range []types.Type{ss, types.NewPointer(ss)} {
							fmt.Printf("Method set of %s:\n", t)
							mset := types.NewMethodSet(t)
							for i := 0; i < mset.Len(); i++ {
								fmt.Println(mset.At(i))
							}
							fmt.Println()
						}

						for i := 0; i < s.NumFields(); i++ {
							spew.Dump(s.Field(i).Type())
						}


											spew.Dump(s.Field(0).Type())
											spew.Dump(s.Field(1).Type())
											spew.Dump(s.Field(2).Type())
											spew.Dump(s.Field(3).Type())
											spew.Dump(s.Field(4).Type())
											spew.Dump(s.Field(5).Type())
											spew.Dump(s.Field(5).Type().(*types.Slice).Elem())
											spew.Dump(s.Field(6).Type())
											spew.Dump(s.Field(7).Type())
											spew.Dump(s.Field(7).Type().(*types.Named).Underlying())
											//fmt.Println(s.Field(1).Type())

						return
					}

				}
			}
	*/
	//p, err := i.Import("auth/storm/example")

	//
	//conf := types.Config{Importer: i}
	//
	//fset := token.NewFileSet()
	//
	//f, err := parser.ParseDir(fset, p.Path(), nil, 0)
	//if err != nil {
	//	log.Fatal(err) // parse error
	//}
	//
	//pkg, err := conf.Check("cmd/hello", fset, []*ast.File{f}, nil)
	//if err != nil {
	//	log.Fatal(err) // type error
	//}
	//
	//fmt.Printf("Package  %q\n", pkg.Path())
	//fmt.Printf("Name:    %s\n", pkg.Name())
	//fmt.Printf("Imports: %s\n", pkg.Imports())
	//fmt.Printf("Scope:   %s\n", pkg.Scope())
	//
	//fmt.Println(p.String())
	//tp, err := load.New(typeName)
	//if err != nil {
	//	errors = append(errors, fmt.Sprintf("[%s] load type: %s", typeName, err))
	//	continue
	//}

	//log.Printf("Calculating graph")
	//g, err := graph.New(tp)
	//if err != nil {
	//	errors = append(errors, fmt.Sprintf("[%s] setting relations: %s", typeName, err))
	//	continue
	//}

	//dialects := dialect.All()

	//log.Printf("Generating code")
	//err = gen.Gen(g, dialects)
	//if err != nil {
	//	errors = append(errors, fmt.Sprintf("[%s] generate code: %s", typeName, err))
	//}
	//}
	//if len(errors) != 0 {
	//	log.Fatalf("Failed:\n%s", strings.Join(errors, "\n"))
	//} else {
	//	log.Printf("Finished successfully!")
	//}
}

func currentPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)
	fmt.Println(filepath.Abs("./"))
	fmt.Println()

	return exPath
}

func pkgPath() string {
	cwd, _ := filepath.Abs("./")
	goPath := filepath.Join(build.Default.GOPATH, "src")
	return strings.TrimPrefix(cwd, goPath+"/")
}

//func loadValidateInterface() *types.Interface {
//	const input = `
//package validate
//type Validate interface {
// Validate() error
//}
//`
//	fset := token.NewFileSet()
//	f, err := parser.ParseFile(fset, "validate.go", input, 0)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Type-check a package consisting of this file.
//	// Type information for the imported packages
//	// comes from $GOROOT/pkg/$GOOS_$GOOARCH/fmt.a.
//	conf := types.Config{Importer: importer.Default()}
//	pkg, err := conf.Check("temperature", fset, []*ast.File{f}, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Print the method sets of Celsius and *Celsius.
//	celsius := pkg.Scope().Lookup("Celsius").Type()
//	for _, t := range []types.Type{celsius, types.NewPointer(celsius)} {
//		fmt.Printf("Method set of %s:\n", t)
//		mset := types.NewMethodSet(t)
//		for i := 0; i < mset.Len(); i++ {
//			fmt.Println(mset.At(i))
//		}
//		fmt.Println()
//	}
//
//
//}
