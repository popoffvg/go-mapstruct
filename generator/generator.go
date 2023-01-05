package generator

import (
	"fmt"
	"go/ast"
	"io"
	"path/filepath"

	"github.com/popoffvg/go-mapstruct"
	"github.com/popoffvg/go-mapstruct/converters"
	"github.com/popoffvg/go-mapstruct/loader"
)

type Config struct {
	srcTypeName string
	dstTypeName string
	dstPkg      string
	srcPkg      string
	Dir         string
}

type Generator struct {
	cfg Config

	loader    *loader.Loader
	templates *TemplateManager
}

type TemplateSettings struct {
	From   *TypeDefinition
	To     *TypeDefinition
	Fields []*mapstruct.FieldDefinition
}

type TypeDefinition struct {
	Pkg  string
	Name string
}

type TemplateManager struct{}

func (m *TemplateManager) Process(settings []*mapstruct.FieldSettings) ([]byte, error) {
	return nil, nil
}

func New(cfg Config) (*Generator, error) {
	if err := cfg.init(); err != nil {
		return nil, fmt.Errorf("init generator failed: %w", err)
	}

	g := Generator{
		cfg: cfg,
	}
	//TODO: stub
	srcPath := "./"
	dstPath := "./"

	l := loader.New(cfg.Dir)
	if err := l.Load(srcPath); err != nil {
		return nil, err
	}

	if err := l.Load(dstPath); err != nil {
		return nil, err
	}

	g.loader = l
	return &g, nil
}

func (g *Generator) Generate(w io.Writer) ([]mapstruct.FieldSettings, error) {
	var (
		src, dst *ast.StructType
		err      error
	)
	if src, err = g.loader.FindType(g.cfg.srcPkg, g.cfg.srcTypeName); err != nil {
		return nil, err
	}

	if dst, err = g.loader.FindType(g.cfg.dstPkg, g.cfg.dstTypeName); err != nil {
		return nil, err
	}

	converter := converters.StructConverter{
		Src:    src,
		Dst:    dst,
		Loader: g.loader,
	}

	return converter.Convert()
}

func (c *Config) init() (err error) {
	c.Dir, err = filepath.Abs(c.Dir)
	return err
}
