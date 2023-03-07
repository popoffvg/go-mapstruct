package generator

import (
	"fmt"
	"go/ast"
	"path/filepath"

	"github.com/popoffvg/go-mapstruct"
	"github.com/popoffvg/go-mapstruct/converters"
	"github.com/popoffvg/go-mapstruct/loader"
	"github.com/popoffvg/go-mapstruct/templates"
)

type Config struct {
	SrcTypeName string
	DstTypeName string
	DstPkg      string
	SrcPkg      string
	Dir         string
}

type Generator struct {
	cfg Config

	loader    *loader.Loader
	templates *templates.Manager
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

func New(cfg Config) (*Generator, error) {
	if err := cfg.init(); err != nil {
		return nil, fmt.Errorf("init generator failed: %w", err)
	}

	g := Generator{
		cfg: cfg,
	}

	l, err := loader.New(cfg.Dir)
	if err != nil {
		return nil, err
	}
	if err := l.Load(cfg.SrcPkg); err != nil {
		return nil, err
	}

	if err := l.Load(cfg.DstPkg); err != nil {
		return nil, err
	}

	g.loader = l
	return &g, nil
}

func (g *Generator) Run() ([]mapstruct.FieldSettings, error) {
	var (
		src, dst *ast.StructType
		err      error
	)
	if src, err = g.loader.FindType(g.cfg.SrcPkg, g.cfg.SrcTypeName); err != nil {
		return nil, err
	}

	if dst, err = g.loader.FindType(g.cfg.DstPkg, g.cfg.DstTypeName); err != nil {
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
	if err != nil {
		return err
	}
	if c.SrcPkg == "" {
		c.SrcPkg = "./"
	}
	if c.DstPkg == "" {
		c.DstPkg = "./"
	}
	return nil
}
