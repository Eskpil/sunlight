package models

import "time"

type Role struct {
	Id string `json:"id"`

	UserId  string `json:"-"`
	GroupId string `json:"-"`

	Group Group `json:"group" gorm:"foreignKey:GroupId"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
