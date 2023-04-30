package models

import "time"

type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "pending"
	VerificationStatusVerified VerificationStatus = "verified"
)

type Verification struct {
	Id string `json:"id" gorm:"primaryKey"`

	Status VerificationStatus `json:"status"`

	MachineId string `json:"-"`
	UserId    string `json:"-"`

	User *User `json:"user,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
