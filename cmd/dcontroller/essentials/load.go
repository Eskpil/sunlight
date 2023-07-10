package essentials

import (
	"errors"
	"os"
)

func Load() error {
	if err := os.Mkdir(os.ExpandEnv("$SUNLIGHT_PKI_DIR/"), 0700); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	if err := os.Mkdir(os.ExpandEnv("$SUNLIGHT_PKI_DIR/certs"), 0700); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	if err := os.Mkdir(os.ExpandEnv("$SUNLIGHT_PKI_DIR/keys"), 0700); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}

	return nil
}
