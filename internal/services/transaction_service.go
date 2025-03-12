package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tahap2/internal/domain"
	"tahap2/internal/workers"
)

type TransactionService struct {
	userRepo        domain.UserRepository
	transactionRepo domain.TransactionRepository
	eventBus        *workers.EventBus
}

func NewTransactionService(userRepo domain.UserRepository, transactionRepo domain.TransactionRepository, eventBus *workers.EventBus) *TransactionService {
	return &TransactionService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		eventBus:        eventBus,
	}
}

func (s *TransactionService) ProcessTopUp(userID uuid.UUID, amount int64) (domain.Transaction, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("user not found")
		}
		return domain.Transaction{}, err
	}

	balBefore := user.Balance
	user.Balance += amount
	balAfter := user.Balance
	err = s.userRepo.UpdateUser(context.TODO(), user)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("error updating user: %w", err)
	}

	newTransaction := domain.Transaction{
		Status:          "success", // direct success status since it doesn't run on background
		UserID:          userID,
		TransactionType: "DEBIT",
		Amount:          amount,
		Remark:          "topup",
		BalanceBefore:   balBefore,
		BalanceAfter:    balAfter,
	}
	err = s.transactionRepo.CreateTransaction(context.TODO(), &newTransaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	return newTransaction, nil
}

func (s *TransactionService) ProcessPayment(userID uuid.UUID, amount int64, remarks string) (domain.Transaction, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("user not found")
		}
		return domain.Transaction{}, err
	}
	if amount > user.Balance {
		return domain.Transaction{}, errors.New("insufficient balance")
	}

	balBefore := user.Balance
	user.Balance -= amount
	balAfter := user.Balance
	err = s.userRepo.UpdateUser(context.TODO(), user)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("error updating user: %w", err)
	}

	newTransaction := domain.Transaction{
		Status:          "success", // direct success status since it doesn't run on background
		UserID:          userID,
		TransactionType: "CREDIT",
		Amount:          amount,
		Remark:          remarks,
		BalanceBefore:   balBefore,
		BalanceAfter:    balAfter,
	}
	err = s.transactionRepo.CreateTransaction(context.TODO(), &newTransaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	return newTransaction, nil
}

func (s *TransactionService) ProcessTransfer(userID, target uuid.UUID, amount int64, remarks string) (domain.Transaction, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("user not found")
		}
		return domain.Transaction{}, err
	}
	if amount > user.Balance {
		return domain.Transaction{}, errors.New("insufficient balance")
	}

	_, err = s.userRepo.GetUserByID(target)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("user target not found")
		}
		return domain.Transaction{}, err
	}

	balBefore := user.Balance
	user.Balance -= amount
	balAfter := user.Balance

	newTransaction := domain.Transaction{
		Status:          "pending", // status will be updated in task queue
		UserID:          userID,
		TransactionType: "CREDIT",
		Amount:          amount,
		Remark:          remarks,
		BalanceBefore:   balBefore,
		BalanceAfter:    balAfter,
	}
	err = s.transactionRepo.CreateTransaction(context.TODO(), &newTransaction)
	if err != nil {
		return domain.Transaction{}, err
	}

	// publish transfer transaction to be process in background.
	go s.eventBus.Publish(workers.EventTypeTransfer, workers.TransferParam{
		TransferInfo: newTransaction,
		TargetID:     target,
	})

	return newTransaction, nil
}

func (s *TransactionService) GetAllTransactions(userID uuid.UUID) ([]*domain.Transaction, error) {
	transactions, err := s.transactionRepo.GetTransactionsByUserID(context.TODO(), userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
