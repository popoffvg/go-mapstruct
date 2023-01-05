package converters

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/popoffvg/go-mapstruct"
	"github.com/popoffvg/go-mapstruct/loader"
)

type fieldConverter struct {
	srcName string
	dstName string

	dstField *ast.Field
	dstType  types.Type

	srcField *ast.Field
	srcType  types.Type

	l *loader.Loader
}

func (c *fieldConverter) convert() (s mapstruct.FieldSettings, err error) {
	if c.dstField.Type == nil {
		return s, fmt.Errorf("destination type not found")
	}

	// TODO: group err
	// TODO: export field if not same pkg
	// TODO: not simple converter
	var settings mapstruct.FieldSettings

	settings.From = c.srcName
	settings.To = c.dstName

	switch {
	// TODO: interface case
	//TODO: struct case
	//TODO: check that dstType also struct
	//TODO: map case
	case c.isConvertableSlice(c.srcType):
		settings.TransformFn = sliceConvert(c.srcName, c.srcName, c.srcType.String())
	case isStruct(c.srcType):
		converter := StructConverter{
			Loader: c.l,
			Src:    c.srcField.Type.(*ast.StructType),
			Dst:    c.dstField.Type.(*ast.StructType),
		}
		settings.SubFieldSettings, err = converter.Convert()
		settings.TransformFn = proxyConverter(settings.SubFieldSettings)
		if err != nil {
			return s, err
		}
	default: // simple type
		settings.TransformFn = directConvert(c.srcName, c.srcName)
	}

	return settings, nil
}

func proxyConverter(fieldSettings []mapstruct.FieldSettings) string {
	// TODO: proxy converter
	return ""
}

func (c *fieldConverter) isConvertableSlice(t types.Type) bool {
	_, isSlice := t.(*types.Slice)
	return isSlice
}

func directConvert(srcName, dstName string) string {
	return fmt.Sprintf("dst.%s = src.%s", srcName, dstName)
}

func sliceConvert(srcName, dstName, typeName string) string {
	// TODO: different types
	return fmt.Sprintf(
		`
		dst.%[2]s = make(%[3]s, 0, len(src.%[1]s))
		for _, v := range src.%[1]s {
			dst.%[2]s = append(dst.%[2]s, v)
		}`,
		srcName, dstName, typeName)
}

func isStruct(t types.Type) bool {
	_, isStruct := t.(*types.Struct)
	return isStruct
}
