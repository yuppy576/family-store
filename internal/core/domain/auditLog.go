package domain

import (
	"encoding/json"
	"time"
)

type AuditAction string

const (
	AuditActionCreate AuditAction = "CREATE"
	AuditActionUpdate AuditAction = "UPDATE"
	AuditActionDelete AuditAction = "DELETE"
)

type AuditLog struct {
	ID           uint64
	UserID       uint64
	UserName     string
	Action       AuditAction
	ResourceType string
	ResourceID   uint64
	OldData      json.RawMessage
	NewData      json.RawMessage
	IPAddress    string
	CreatedAt    time.Time
	StoreID      uint64
}
