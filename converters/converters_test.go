package converters

import (
	"testing"

	"github.com/popoffvg/go-mapstruct"
	"github.com/popoffvg/go-mapstruct/loader"
	"github.com/stretchr/testify/suite"
)

type FieldSuite struct {
	suite.Suite

	l *loader.Loader
}

type ( // TestConvertSimpleType
	SimpleSrc struct {
		StringField        string
		IntField           int
		Int32Field         int32
		Int64Field         int64
		StringFieldPointer *string
	}

	SimpleDSt struct {
		StringField        string
		IntField           int
		Int32Field         int32
		Int64Field         int64
		StringFieldPointer *string
	}
)

func TestConverters(t *testing.T) {
	suite.Run(t, new(FieldSuite))
}

func (s *FieldSuite) SetupSuite() {
	s.l = loader.New("./")
	s.Require().NoError(s.l.Load("./"))
}

func (s *FieldSuite) TestConvertSimpleType() {

	src, err := s.l.FindType("converters", "SimpleSrc")
	s.Require().NoError(err)
	dst, err := s.l.FindType("converters", "SimpleDSt")
	s.Require().NoError(err)
	c := StructConverter{
		Loader: s.l,
		Src:    src,
		Dst:    dst,
	}

	result, err := c.Convert()
	s.Assert().NoError(err)
	s.Assert().Equal([]mapstruct.FieldSettings{
		{
			From:             "StringField",
			To:               "StringField",
			TransformFn:      "dst.StringField = src.StringField",
			SubFieldSettings: nil,
		},
		{
			From:             "IntField",
			To:               "IntField",
			TransformFn:      "dst.IntField = src.IntField",
			SubFieldSettings: nil,
		},
		{
			From:             "Int32Field",
			To:               "Int32Field",
			TransformFn:      "dst.Int32Field = src.Int32Field",
			SubFieldSettings: nil,
		},
		{
			From:             "Int64Field",
			To:               "Int64Field",
			TransformFn:      "dst.Int64Field = src.Int64Field",
			SubFieldSettings: nil,
		},
		{
			From:             "StringFieldPointer",
			To:               "StringFieldPointer",
			TransformFn:      "dst.StringFieldPointer = src.StringFieldPointer",
			SubFieldSettings: nil,
		},
	}, result)
}

type ( // TestConvertTypeToPointer
	SimplePointerSrc struct {
		StringField string
		IntField    *int
	}

	SimplePointerDSt struct {
		StringField *string
		IntField    int
	}
)

func (s *FieldSuite) TestConvertTypeToPointer() {
	s.T().Skip("not implemented")
	src, err := s.l.FindType("converters", "SimplePointerSrc")
	s.Require().NoError(err)
	dst, err := s.l.FindType("converters", "SimplePointerDSt")
	s.Require().NoError(err)
	c := StructConverter{
		Loader: s.l,
		Src:    src,
		Dst:    dst,
	}

	result, err := c.Convert()
	s.Assert().NoError(err)

	s.Assert().Equal([]mapstruct.FieldSettings{
		{
			From: "StringField",
			To:   "StringField",
			TransformFn: `StringField := src.StringField
			dst.StringField = &StringField`,
			SubFieldSettings: nil,
		},
		{
			From: "IntField",
			To:   "IntField",
			TransformFn: `if src.IntField != nil {
				dst.IntField = *src.IntField
			}`,
			SubFieldSettings: nil,
		},
	}, result)
}

type ( // TestConvertSliceSimpleType
	SliceSimpleSrc struct {
		Field []string
	}

	SliceSimpleDst struct {
		Field []string
	}
)

func (s *FieldSuite) TestConvertSliceSimpleType() {
	src, err := s.l.FindType("converters", "SliceSimpleSrc")
	s.Require().NoError(err)
	dst, err := s.l.FindType("converters", "SliceSimpleDst")
	s.Require().NoError(err)
	c := StructConverter{
		Loader: s.l,
		Src:    src,
		Dst:    dst,
	}

	result, err := c.Convert()
	s.Assert().NoError(err)
	s.Assert().Equal([]mapstruct.FieldSettings{
		{
			From: "Field",
			To:   "Field",
			TransformFn: `
		dst.Field = make([]string, 0, len(src.Field))
		for _, v := range src.Field {
			dst.Field = append(dst.Field, v)
		}`,
			SubFieldSettings: nil,
		},
	}, result)
}
func (s *FieldSuite) TestConvertAliasTypes() {
}

func (s *FieldSuite) TestConvertMap() {
}

func (s *FieldSuite) TestConvertStruct() {
}

func (s *FieldSuite) TestConvertSliceStructType() {
}

func (s *FieldSuite) TestConvertConvertableInterface() {
}

func (s *FieldSuite) TestConvertConvertableInterfaceWithError() {
}
