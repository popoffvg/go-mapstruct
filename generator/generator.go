package generator

import (
	"fmt"
	"go/ast"
	"io"
	"path/filepath"

	"github.com/popoffvg/go-mapstruct"
	"github.com/popoffvg/go-mapstruct/converters"
	"github.com/popoffvg/go-mapstruct/loader"
	"github.com/popoffvg/go-mapstruct/templates"
)

type Config struct {
	srcTypeName string
	dstTypeName string
	dstPkgPath  string
	srcPkgPath  string
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

	l := loader.New(cfg.Dir)
	if err := l.Load(cfg.srcPkgPath); err != nil {
		return nil, err
	}

	if err := l.Load(cfg.dstPkgPath); err != nil {
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
	if src, err = g.loader.FindType(g.cfg.srcPkgPath, g.cfg.srcTypeName); err != nil {
		return nil, err
	}

	if dst, err = g.loader.FindType(g.cfg.dstPkgPath, g.cfg.dstTypeName); err != nil {
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
	if c.srcPkgPath == "" {
		c.srcPkgPath = "./"
	}
	if c.dstPkgPath == "" {
		c.dstPkgPath = "./"
	}
	return nil
}
