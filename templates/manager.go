package templates

import mapstruct "github.com/popoffvg/go-mapstruct"

type (
	Manager struct{}
)

func New() *Manager {
	return &Manager{}
}

func (m *Manager) Process(settings []mapstruct.FieldSettings) ([]byte, error) {
	return []byte("some code snippet"), nil
}
