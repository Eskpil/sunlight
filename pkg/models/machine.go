package models

import "time"

// MachineStatus and MachineState differs, MachineStatus tells us what status it has inside sunlight. And MachineState
// tells us the state of the machines. This state could whether it is shutdown, its rebooting. Hibernation or running
// normally. MachineStatus could be, NeedsVerification, this means an operator must verify the machines before it gets
// access to the full sunlight system. Another state might be Operational, this means the machines is completely verified
// and is ready to be incorporated with the rest of the sunlight system.
type MachineStatus string
type MachineState string

const (
	MachineStatusPending     MachineStatus = "pending"
	MachineStatusOperational MachineStatus = "operational"
	MachineStatusUnreachable MachineStatus = "unreachable"

	MachineStateShutdown    MachineStatus = "shutdown"
	MachineStateRebooting   MachineState  = "rebooting"
	MachineStateHibernating MachineState  = "hibernating"
	MachineStateOperational MachineState  = "operational"
)

type Machine struct {
	Id       string `json:"id" gorm:"primaryKey"`
	Hostname string `json:"hostname" gorm:"unique"`

	Status MachineStatus `json:"status"`
	State  MachineState  `json:"state"`

	PublicKey string `json:"public_key"`

	VerificationId string        `json:"-"`
	Verification   *Verification `json:"verification,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
