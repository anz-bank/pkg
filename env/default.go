package env

import "os"

type defaultEnv struct{}

var _ Env = defaultEnv{}

func (defaultEnv) Clearenv() {
	os.Clearenv()
}

func (defaultEnv) Environ() []string {
	return os.Environ()
}

func (defaultEnv) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (defaultEnv) Setenv(key, value string) error {
	return os.Setenv(key, value)
}

func (defaultEnv) Unsetenv(key string) error {
	return os.Unsetenv(key)
}
