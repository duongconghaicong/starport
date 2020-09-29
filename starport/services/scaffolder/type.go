package scaffolder

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/tendermint/starport/starport/pkg/cosmosver"
	"github.com/tendermint/starport/starport/pkg/gomodulepath"
	"github.com/tendermint/starport/starport/templates/typed"
)

// AddType adds a new type stype to scaffolded app by using optional type fields.
func (s *Scaffolder) AddType(stype string, fields ...string) error {
	path, err := gomodulepath.ParseFile(s.path)
	if err != nil {
		return err
	}
	ok, err := isTypeCreated(s.path, path.Package, stype)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("%s type is already added.", stype)
	}
	var tfields []typed.Field
	for _, f := range fields {
		fs := strings.Split(f, ":")
		name := fs[0]
		var datatype string
		acceptedTypes := map[string]bool{
			"string": true,
			"bool":   true,
			"int":    true,
			"float":  true,
		}
		if len(fs) == 2 && acceptedTypes[fs[1]] {
			datatype = fs[1]
		} else {
			datatype = "string"
		}
		tfields = append(tfields, typed.Field{
			Name:     name,
			Datatype: datatype,
		})
	}
	version, err := s.version()
	if err != nil {
		return err
	}

	var (
		g    *genny.Generator
		opts = &typed.Options{
			ModulePath: path.RawPath,
			AppName:    path.Package,
			TypeName:   stype,
			Fields:     tfields,
		}
	)
	if version == cosmosver.Launchpad {
		g, err = typed.NewLaunchpad(opts)
	} else {
		g, err = typed.NewStargate(opts)
	}
	if err != nil {
		return err
	}
	run := genny.WetRunner(context.Background())
	run.With(g)
	if err := run.Run(); err != nil {
		return err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	return s.protoc(pwd, version)
}

func isTypeCreated(appPath, appName, typeName string) (isCreated bool, err error) {
	abspath, err := filepath.Abs(filepath.Join(appPath, "x", appName, "types"))
	if err != nil {
		return false, err
	}
	fset := token.NewFileSet()
	all, err := parser.ParseDir(fset, abspath, func(os.FileInfo) bool { return true }, parser.ParseComments)
	if err != nil {
		return false, err
	}
	for _, pkg := range all {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(x ast.Node) bool {
				typeSpec, ok := x.(*ast.TypeSpec)
				if !ok {
					return true
				}
				if _, ok := typeSpec.Type.(*ast.StructType); !ok {
					return true
				}
				if "Msg"+strings.Title(typeName) != typeSpec.Name.Name {
					return true
				}
				isCreated = true
				return false
			})
		}
	}
	return
}