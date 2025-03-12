package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Status          string    `gorm:"not null" json:"status"`
	UserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TransactionType string    `gorm:"not null" json:"transaction_type"`
	Amount          int64     `gorm:"not null" json:"amount"`
	Remark          string    `gorm:"not null" json:"remark"`
	BalanceBefore   int64     `gorm:"not null" json:"balance_before"`
	BalanceAfter    int64     `gorm:"not null" json:"balance_after"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type TransactionService interface {
	ProcessTopUp(userID uuid.UUID, amount int64) (Transaction, error)
	ProcessPayment(userID uuid.UUID, amount int64, remarks string) (Transaction, error)
	ProcessTransfer(userID, target uuid.UUID, amount int64, remarks string) (Transaction, error)
	GetAllTransactions(userID uuid.UUID) ([]*Transaction, error)
}

type TransactionRepository interface {
	GetTransactionByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	CreateTransaction(ctx context.Context, transaction *Transaction) error
	GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]*Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *Transaction) error
}
