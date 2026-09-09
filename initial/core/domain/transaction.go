package domain

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            int64
	WalletID      uuid.UUID
	OperationType OperationType
	Amount        int64
	CreatedAt     time.Time
}