package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tahap2/internal/domain"
)

type TransactionRepo struct {
	DB *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{DB: db}
}

func (r *TransactionRepo) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	return r.DB.Create(transaction).Error
}

func (r *TransactionRepo) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error) {
	var trans []*domain.Transaction
	err := r.DB.Where("user_id = ?", userID).Find(&trans).Error
	if err != nil {
		return nil, err
	}
	return trans, nil
}

func (r *TransactionRepo) GetTransactionByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	var trans *domain.Transaction
	err := r.DB.Where("id = ?", id).First(&trans).Error
	if err != nil {
		return nil, err
	}
	return trans, nil
}

func (r *TransactionRepo) UpdateTransaction(ctx context.Context, trans *domain.Transaction) error {
	return r.DB.Save(trans).Error
}
