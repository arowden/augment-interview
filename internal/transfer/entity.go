package transfer

import (
	"time"

	"github.com/google/uuid"
)

type Transfer struct {
	ID             uuid.UUID
	FundID         uuid.UUID
	FromOwner      string
	ToOwner        string
	Units          int
	IdempotencyKey *uuid.UUID
	TransferredAt  time.Time
}
