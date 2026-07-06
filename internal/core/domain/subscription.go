package domain

import "time"

type SubscriptionPlan string

const (
	PlanTrial     SubscriptionPlan = "TRIAL"
	PlanPersonal  SubscriptionPlan = "PERSONAL"
	PlanProfessional SubscriptionPlan = "PROFESSIONAL"
	PlanLifetime  SubscriptionPlan = "LIFETIME"
)

type Subscription struct {
	ID        uint64            `json:"id"`
	StoreID   uint64            `json:"store_id"`
	Plan      SubscriptionPlan  `json:"plan"`
	Status    SubscriptionStatus `json:"status"`
	StartDate time.Time         `json:"start_date"`
	EndDate   *time.Time        `json:"end_date"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type SubscriptionStatus string

const (
	SubActive   SubscriptionStatus = "ACTIVE"
	SubTrial    SubscriptionStatus = "TRIAL"
	SubExpired  SubscriptionStatus = "EXPIRED"
	SubFrozen   SubscriptionStatus = "FROZEN"
)