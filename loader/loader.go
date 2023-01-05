package loader

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/popoffvg/go-mapstruct/helpers"
	"golang.org/x/tools/go/packages"
)

var errPackageNotFound = fmt.Errorf("package not found")

type (
	Loader struct {
		dir string

		cfg *packages.Config

		fs    *token.FileSet
		pkgs  map[string]*packages.Package
		ast   map[string]*ast.Package
		types *types.Info
		files map[string]*ast.File
	}
)

func New(dir string) *Loader {
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	return &Loader{
		types: info,
		fs:    token.NewFileSet(),
		dir:   dir,
		files: make(map[string]*ast.File),
		pkgs:  make(map[string]*packages.Package),
		ast:   make(map[string]*ast.Package),
		cfg: &packages.Config{
			Dir:  dir,
			Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps,
		},
	}
}

// Load loads package by its import path
func (l *Loader) Load(path string) error {

	pkgs, err := packages.Load(l.cfg, filepath.Join(l.dir, path))
	if err != nil {
		return err
	}

	if len(pkgs) < 1 {
		return errPackageNotFound
	}

	if len(pkgs[0].Errors) > 0 {
		return pkgs[0].Errors[0]
	}

	l.pkgs[path] = pkgs[0]

	currentAST, err := l.loadAST(path)
	if err != nil {
		return err
	}

	conf := types.Config{
		IgnoreFuncBodies: true,
		FakeImportC:      true,
		Error:            nil,
		Importer:         importer.ForCompiler(l.fs, "source", nil),
	}

	var files []*ast.File
	for fPath, f := range currentAST.Files {
		if _, ok := l.files[fPath]; ok {
			continue
		}
		files = append(files, f)
	}

	_, err = conf.Check("./", l.fs, files, l.types)
	if err != nil {
		return fmt.Errorf("failed collect info: %w, files: %v", err, files)
	}

	return nil
}

func (l *Loader) GetType(expr ast.Expr) (types.TypeAndValue, error) {
	if t, ok := l.types.Types[expr]; ok {
		return t, nil
	}

	return types.TypeAndValue{}, fmt.Errorf("expression %#v not found", expr)
}

func (l *Loader) AST(path string) (*ast.Package, error) {
	if ast, ok := l.ast[path]; ok {
		return ast, nil
	}

	return nil, fmt.Errorf("%s not load", path)
}

// AST returns package's abstract syntax tree
func (l *Loader) loadAST(path string) (*ast.Package, error) {
	var (
		dir string
		p   *packages.Package
		ok  bool
	)
	if p, ok = l.pkgs[path]; ok {
		dir = Dir(p)
	}
	pkgs, err := parser.ParseDir(l.fs, dir, nil, parser.DeclarationErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}

	if ap, ok := pkgs[p.Name]; ok {
		l.ast[p.Name] = ap
		return ap, nil
	}

	return &ast.Package{Name: p.Name}, nil
}

// Dir returns absolute path of the package in a filesystem
func Dir(p *packages.Package) string {
	files := append(p.GoFiles, p.OtherFiles...)
	if len(files) < 1 {
		return p.PkgPath
	}

	return filepath.Dir(files[0])
}

func (l *Loader) FindType(pkg, name string) (*ast.StructType, error) {
	a, ok := l.ast[pkg]
	if !ok {
		return nil, fmt.Errorf(
			"ast for package:%s not found. Loaded pkgs: %#v",
			pkg,
			helpers.ToArray(l.ast),
		)
	}
	for _, f := range a.Files {
		for _, ts := range typeSpecs(f) {
			if i, ok := ts.Type.(*ast.StructType); ts.Name.Name == name {
				if ok {
					return i, nil
				}
				return nil, fmt.Errorf("type %s is not struct but was %T", name, ts.Type)
			}
		}
	}

	return nil, fmt.Errorf("type %s not found", name)
}

func typeSpecs(f *ast.File) []*ast.TypeSpec {
	var result []*ast.TypeSpec

	for _, decl := range f.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.TYPE {
			for _, spec := range gd.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					result = append(result, ts)
				}
			}
		}
	}

	return result
}
