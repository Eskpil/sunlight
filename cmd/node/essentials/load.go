package essentials

import (
	"errors"
	"os"
)

func Load() error {
	if err := os.Mkdir(os.ExpandEnv("$SUNLIGHT_VAR_DIR/csr"), 0700); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	if err := os.Mkdir(os.ExpandEnv("$SUNLIGHT_VAR_DIR/certs"), 0700); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	if err := os.Mkdir(os.ExpandEnv("$SUNLIGHT_VAR_DIR/keys"), 0700); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	return nil
}
