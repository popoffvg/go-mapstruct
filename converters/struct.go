package converters

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/hashicorp/go-multierror"
	"github.com/popoffvg/go-mapstruct"
	"github.com/popoffvg/go-mapstruct/loader"
)

type StructConverter struct {
	Src, Dst *ast.StructType
	DstIndex map[string]*ast.Field
	Loader   *loader.Loader
}

func (c *StructConverter) Convert() ([]mapstruct.FieldSettings, error) {
	var err error

	// TODO: error level
	var (
		result  []mapstruct.FieldSettings
		dstF    *ast.Field
		fs      mapstruct.FieldSettings
		ok      bool
		convErr error
	)
	c.DstIndex = make(map[string]*ast.Field, 0)
	for _, f := range c.Dst.Fields.List {
		c.DstIndex[f.Names[0].Name] = f
	}

	fc := fieldConverter{}

	for _, f := range c.Src.Fields.List {
		srcName := f.Names[0].Name
		// TODO: remap fields
		// TODO: why massive of names
		if dstF, ok = c.DstIndex[srcName]; !ok { // check that field has receiver
			// TODO: milti level struct
			err = multierror.Append(err, fmt.Errorf("not found destination field for %s", srcName))
			continue
		}

		fc.srcType, convErr = extractType(c.Loader.GetType(f.Type))
		if convErr != nil {
			err = multierror.Append(err, fmt.Errorf("get type of field %s failed: %w", srcName, convErr))
			continue
		}

		fc.dstType, convErr = extractType(c.Loader.GetType(dstF.Type))
		if convErr != nil {
			// TODO: field remap
			err = multierror.Append(err, fmt.Errorf("get type of field %s failed: %w", srcName, convErr))
			continue
		}

		// TODO: field remap
		fc.srcName = srcName
		fc.dstName = srcName
		fc.srcField = f
		fc.dstField = dstF
		fs, convErr = fc.convert()
		if convErr != nil {
			err = multierror.Append(err, convErr)
			continue
		}
		result = append(result, fs)
	}
	return result, err
}

func extractType(typeAndValue types.TypeAndValue, err error) (types.Type, error) {
	if err != nil {
		return nil, err
	}
	return typeAndValue.Type, nil
}
