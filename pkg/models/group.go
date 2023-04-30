package models

import "time"

type Group struct {
	Id    string `json:"id"`
	Users []Role `json:"users"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
