package transfer

import (
	"time"

	"github.com/google/uuid"
)

// Transfer represents an immutable record of units transferred between owners.
type Transfer struct {
	ID             uuid.UUID
	FundID         uuid.UUID
	FromOwner      string
	ToOwner        string
	Units          int
	IdempotencyKey *uuid.UUID
	TransferredAt  time.Time
}
