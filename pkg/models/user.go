package models

import (
	"context"
	"fmt"
	"github.com/matthewhartstonge/argon2"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id       string `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`

	Password []byte `json:"-"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Keys          []Key          `json:"keys"`
	Roles         []Role         `json:"roles"`
	Verifications []Verification `json:"verifications"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Save(u).Error
}

func (u *User) EncryptPassword(password string) error {
	argon := argon2.DefaultConfig()

	encoded, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return nil
	}

	u.Password = encoded

	return nil
}

func (u *User) VerifyPassword(other string) error {
	ok, err := argon2.VerifyEncoded([]byte(other), u.Password)
	if err != nil {
		return nil
	}

	if !ok {
		return fmt.Errorf("passwords does not match")
	}

	return nil
}

// Delete deletes the user from the database using the given GORM database handle
func (u *User) Delete(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Delete(u).Error
}

// AddRole adds the given role to the user's roles and saves the user to the database using the given GORM database handle
func (u *User) AddRole(ctx context.Context, db *gorm.DB, role Role) error {
	u.Roles = append(u.Roles, role)
	return u.Save(ctx, db)
}

// RemoveRole removes the given role from the user's roles and saves the user to the database using the given GORM database handle
func (u *User) RemoveRole(ctx context.Context, db *gorm.DB, role Role) error {
	for i, r := range u.Roles {
		if r.Id == role.Id {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			return u.Save(ctx, db)
		}
	}
	return nil
}

// GetKeys retrieves all the keys belonging to the user from the database using the given GORM database handle
func (u *User) GetKeys(ctx context.Context, db *gorm.DB) ([]Key, error) {
	var keys []Key
	err := db.WithContext(ctx).Model(u).Association("Keys").Find(&keys).Error
	if err != nil {
		return nil, fmt.Errorf(err())
	}
	return keys, nil
}

// AddKey adds the given key to the user's keys and saves the user to the database using the given GORM database handle
func (u *User) AddKey(ctx context.Context, db *gorm.DB, key Key) error {
	u.Keys = append(u.Keys, key)
	return u.Save(ctx, db)
}

// RemoveKey removes the given key from the user's keys and saves the user to the database using the given GORM database handle
func (u *User) RemoveKey(ctx context.Context, db *gorm.DB, key Key) error {
	for i, k := range u.Keys {
		if k.Id == key.Id {
			u.Keys = append(u.Keys[:i], u.Keys[i+1:]...)
			return u.Save(ctx, db)
		}
	}
	return nil
}
