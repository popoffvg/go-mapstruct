package generator

import (
	"bytes"
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
		srcTypeName: "src",
		dstTypeName: "dst",
		// TODO: pkg with same names
		srcPkgPath: ".",
		dstPkgPath: ".",
	}
	g, err := New(cfg)
	require.NoError(t, err, fmt.Sprintf("cfg: %#v", cfg))

	b := bytes.NewBuffer(make([]byte, 0))
	_, err = g.Generate(b)
	assert.Error(t, err)
}
