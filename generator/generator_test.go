package generator

import (
	"fmt"
	. "net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

// TODO: embeded struct
type src struct {
	SimpleField string
	Interface   Listener
}

type dst struct {
	SimpleField string
}

func TestSimpleStruct(t *testing.T) {
	cfg := Config{
		Dir:         ".", // TODO: stub
		SrcTypeName: "src",
		DstTypeName: "dst",
		// TODO: pkg with same names
		SrcPkg: "generator",
		DstPkg: "generator",
	}
	g, err := New(cfg)
	require.NoError(t, err, fmt.Sprintf("cfg: %#v", cfg))

	_, err = g.Run()
	assert.Error(t, err)
}
