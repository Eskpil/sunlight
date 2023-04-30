package users

import (
	"fmt"
	"strings"
)

type CreateInput struct {
	Username string `json:"username"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Email string `json:"email"`

	Password string `json:"password"`
}

func (i *CreateInput) Validate() error {
	if i.FirstName == "" && i.LastName == "" {
		return fmt.Errorf("missing first name and lastname")
	}

	if i.Username == "" {
		i.Username = strings.ToLower(i.FirstName) + strings.ToLower(string(i.LastName[0]))
	}

	if i.Email == "" {
		i.Email = fmt.Sprintf("%s.%s@sunlight.local", i.FirstName, i.LastName)
	}

	return nil
}
