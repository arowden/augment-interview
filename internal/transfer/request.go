package transfer

import "github.com/google/uuid"

type Request struct {
	FundID         uuid.UUID
	FromOwner      string
	ToOwner        string
	Units          int
	IdempotencyKey *uuid.UUID
}
