package loader

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/popoffvg/go-mapstruct/helpers"
	"golang.org/x/tools/go/packages"
)

var errPackageLoadFailed = fmt.Errorf("packages load failed")

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

func New(dir string) (_ *Loader, err error) {
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		errPackageLoadFailed = err
		return &Loader{}, fmt.Errorf("failed get absolute path for %s: %w", dir, err)
	}

	l := &Loader{
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

	err = l.loadRecursively()
	if err != nil {
		return l, err
	}
	return l, nil
}

// Load loads package by its import path
func (l *Loader) Load(path string) (err error) {
	if len(l.pkgs) < 1 {
		return errPackageLoadFailed
	}

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
	if p, ok = l.pkgs[path]; !ok {
		return nil, fmt.Errorf("not found pkg: %s", path)
	}
	dir = Dir(p)

	pkgs, err := parser.ParseDir(l.fs, dir, nil, parser.DeclarationErrors|parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("load ast failed for %s(package) in %s (directory): %w", path, dir, err)
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
			"ast for package:%s not found. Loaded pkgs: %s",
			pkg,
			pkgNames(helpers.ToArray(l.ast)),
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

func pkgNames(packages []*ast.Package) string {
	var r []string
	for _, p := range packages {
		r = append(r, p.Name)
	}
	return strings.Join(r, ", ")
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

func (l *Loader) loadRecursively() error {
	return filepath.WalkDir(l.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed load packages: %w", err)
		}
		if !d.IsDir() {
			return nil
		}

		if strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		pkgs, err := packages.Load(l.cfg, path)
		if err != nil {
			return fmt.Errorf("failed load packages: %w", err)
		}

		if len(pkgs) != 1 || pkgs[0].Name == "" {
			// dir with only subdirectories or without *.go files
			return nil
		}

		l.pkgs[pkgs[0].Name] = pkgs[0]

		return nil
	})
}
