package mapstruct

type (
	FieldDefinition struct {
		// TODO: add "to" and "from" fields
		Name    string
		Type    string
		Package string

		isArray  bool
		isStruct bool
	}
)

func (d *FieldDefinition) IsArray() bool {
	return d.isArray
}

func (d *FieldDefinition) IsStruct() bool {
	return d.isStruct
}

func (d *FieldDefinition) Convert(srcName, dstName string) string {
	return ""
}

func (d *FieldDefinition) ConvertArray(srcName, dstName, iteratorName string) string {
	return ""
}
