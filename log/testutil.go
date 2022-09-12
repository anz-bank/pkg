package log

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/arr-ai/frozen"
	"github.com/stretchr/testify/require"
)

type multipleFields struct {
	Name                 string
	Fields, GlobalFields frozen.Map
}

func generateMultipleFieldsCases() []multipleFields {
	return []multipleFields{
		{
			Name: "Multiple types of Values",
			Fields: frozen.NewMap(
				frozen.KV("byte", '1'),
				frozen.KV("int", 123),
				frozen.KV("string", "this is an unnecessarily long sentence"),
			),
		},
		{
			Name:   "Empty Key",
			Fields: frozen.NewMap(frozen.KV("", "stuff")),
		},
		{
			Name:   "Nil Value",
			Fields: frozen.NewMap(frozen.KV("Nil", nil)),
		},
	}
}

// Adapted from https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
func redirectOutput(t *testing.T, print func()) string {
	old := os.Stderr
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stderr = w

	print()

	outC := make(chan string)
	go func(t *testing.T) {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		require.NoError(t, err)
		outC <- buf.String()
	}(t)

	w.Close()
	os.Stderr = old
	return <-outC
}

func convertToGoMap(fields frozen.Map) map[string]interface{} {
	goMap := make(map[string]interface{})
	for i := fields.Range(); i.Next(); {
		goMap[i.Key().(string)] = i.Value()
	}
	return goMap
}
