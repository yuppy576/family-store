package domain

import "time"

type StoreStatus string

const (
	StoreActive   StoreStatus = "ACTIVE"
	StoreTrial    StoreStatus = "TRIAL"
	StoreFrozen   StoreStatus = "FROZEN"
	StoreExpired  StoreStatus = "EXPIRED"
)

type Store struct {
	ID        uint64     `json:"id"`
	Name      string     `json:"name"`
	Domain    string     `json:"domain"`
	Status    StoreStatus `json:"status"`
	TrialEnd  *time.Time `json:"trial_end"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}