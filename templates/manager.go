package templates

import "github.com/popoffvg/go-mapstruct"

type (
	Manager struct{}
)

func (m *Manager) Process(settings []*mapstruct.FieldSettings) ([]byte, error) {
	return nil, nil
}
