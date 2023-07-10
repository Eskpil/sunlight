package machines

import (
	"fmt"
)

type CreateInput struct {
	Hostname  string `json:"hostname"`
	PublicKey string `json:"public_key"`
}

func (i *CreateInput) Validate() error {
	if i.Hostname == "" {
		return fmt.Errorf("hostname is missing")
	}

	if i.PublicKey == "" {
		return fmt.Errorf("public key is missing")
	}

	return nil
}
