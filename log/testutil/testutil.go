package testutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/marcelocantos/frozen"
	"github.com/stretchr/testify/require"
)

type SingleField struct {
	Name, Key string
	Val       interface{}
}

type MultipleFields struct {
	Name   string
	Fields frozen.Map
}

func GenerateSingleFieldCases() []SingleField {
	return []SingleField{
		{
			Name: "String Value",
			Key:  "random",
			Val:  "Value",
		},
		{
			Name: "Number Value",
			Key:  "int",
			Val:  3,
		},
		{
			Name: "Byte Value",
			Key:  "byte",
			Val:  'q',
		},
		{
			Name: "Empty Key",
			Key:  "",
			Val:  "Empty",
		},
		{
			Name: "Empty Value",
			Key:  "Empty",
			Val:  "",
		},
		{
			Name: "Nil Value",
			Key:  "nil",
			Val:  nil,
		},
	}
}

func GenerateMultipleFieldsCases() []MultipleFields {
	return []MultipleFields{
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
func RedirectOutput(t *testing.T, print func()) string {
	old := os.Stderr
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stderr = w

	print()

	outC := make(chan string)
	go func(tt *testing.T) {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		require.NoError(tt, err)
		outC <- buf.String()
	}(t)

	w.Close()
	os.Stderr = old
	return <-outC
}

func OutputFormattedFields(fields frozen.Map) string {
	if fields.Count() == 0 {
		return ""
	}

	keys := make([]string, fields.Count())
	index := 0
	for k := fields.Range(); k.Next(); {
		keys[index] = k.Key().(string)
		index++
	}

	sort.Strings(keys)

	output := strings.Builder{}
	output.WriteString(fmt.Sprintf("%s=%v", keys[0], fields.MustGet(keys[0])))

	if fields.Count() > 1 {
		for _, keyField := range keys[1:] {
			output.WriteString(fmt.Sprintf(" %s=%v", keyField, fields.MustGet(keyField)))
		}
	}

	return output.String()
}

func GetSortedKeys(fields frozen.Map) []string {
	keys := make([]string, fields.Count())
	index := 0
	for i := fields.Range(); i.Next(); {
		keys[index] = i.Key().(string)
		index++
	}
	sort.Strings(keys)
	return keys
}

func ConvertToGoMap(fields frozen.Map) map[string]interface{} {
	goMap := make(map[string]interface{})
	for i := fields.Range(); i.Next(); {
		goMap[i.Key().(string)] = i.Value()
	}
	return goMap
}
