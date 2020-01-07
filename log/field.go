package log

import "github.com/marcelocantos/frozen"

// Fields is an interface that represents a collection of key values
type Fields interface {
	GetKeyValues() []frozen.KeyValue
}

// SingleField represents a single key value
type SingleField struct {
	key   string
	value interface{}
}

// MultipleFields is a type for map of string to interface, meant for type checking
// you can use FromMap to create it or directly use MultipleFields{}
type MultipleFields map[string]interface{}

// NewField creates a Fields interface representing a single key value
func NewField(key string, value interface{}) Fields {
	return &SingleField{key: key, value: value}
}

func (sf *SingleField) GetKeyValues() []frozen.KeyValue {
	return []frozen.KeyValue{frozen.KV(sf.key, sf.value)}
}

func FromMap(m map[string]interface{}) Fields {
	return MultipleFields(m)
}

func (mf MultipleFields) GetKeyValues() []frozen.KeyValue {
	i := 0
	keyVals := make([]frozen.KeyValue, len(mf))
	for k, v := range mf {
		keyVals[i] = frozen.KV(k, v)
		i++
	}
	return keyVals
}
