package mod

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

func ModInit(modName string) error {
	err := runGo(context.Background(), ioutil.Discard, "mod", "init", modName)
	if err != nil {
		return errors.New(fmt.Sprintf("go mod init failed: %s", err.Error()))
	}

	return nil
}
