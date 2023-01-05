package generator

import (
	"bytes"
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
	g, err := New(Config{
		Dir:         ".", // TODO: stub
		srcTypeName: "src",
		dstTypeName: "dst",
		// TODO: pkg with same names
		srcPkg: "generator",
		dstPkg: "generator",
	})
	require.NoError(t, err)

	b := bytes.NewBuffer(make([]byte, 0))
	_, err = g.Generate(b)
	assert.Error(t, err)
}
