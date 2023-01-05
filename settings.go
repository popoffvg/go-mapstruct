package mapstruct

type FieldSettings struct {
	From        string
	To          string
	TransformFn string

	SubFieldSettings []FieldSettings
}

type Option func(*FieldSettings)

var (
	Field = func(field string) Option {
		return func(s *FieldSettings) {
			s.From = field
		}
	}

	To = func(field string) Option {
		return func(s *FieldSettings) {
			s.To = field
		}
	}

	With = func(fn any) Option {
		return func(s *FieldSettings) {
			// TODO: transform func
			s.TransformFn = ""
		}
	}

	Map = func(settings FieldSettings) Option {
		return func(s *FieldSettings) {
			s.SubFieldSettings = append(s.SubFieldSettings, settings)
		}
	}
)

func NewSetting(from, to string, fnBody string) *FieldSettings {
	return &FieldSettings{
		From:        from,
		To:          to,
		TransformFn: fnBody,
	}
}
