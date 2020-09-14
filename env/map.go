package env

import "sync"

// Map implements an Env backed on a map[string]string.
type Map struct {
	mux sync.RWMutex
	env map[string]string
}

func NewMap(initialEnv map[string]string) *Map {
	var env map[string]string
	if initialEnv != nil {
		env = make(map[string]string, len(initialEnv))
		for key, value := range initialEnv {
			env[key] = value
		}
	} else {
		env = map[string]string{}
	}
	return &Map{env: env}
}

func (m *Map) Clearenv() {
	m.mux.Lock()
	defer m.mux.Unlock()
	for key := range m.env {
		delete(m.env, key)
	}
}

func (m *Map) Environ() []string {
	m.mux.RLock()
	defer m.mux.RUnlock()
	result := make([]string, 0, len(m.env))
	for key, value := range m.env {
		result = append(result, key+"="+value)
	}
	return result
}

func (m *Map) LookupEnv(key string) (string, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	value, has := m.env[key]
	return value, has
}

func (m *Map) Setenv(key, value string) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.env[key] = value
	return nil
}

func (m *Map) Unsetenv(key string) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.env, key)
	return nil
}
