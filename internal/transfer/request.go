package transfer

import "github.com/google/uuid"

// Request represents a transfer request with validation-ready fields.
type Request struct {
	FundID         uuid.UUID
	FromOwner      string
	ToOwner        string
	Units          int
	IdempotencyKey *uuid.UUID
}
