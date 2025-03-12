package workers

import (
	"context"
	"log"
	"tahap2/internal/domain"
)

type TransactionWorker struct {
	eventBus        *EventBus
	userRepository  domain.UserRepository
	transactionRepo domain.TransactionRepository
}

func NewTransactionWorker(eventBus *EventBus, userRepo domain.UserRepository, transactionRepo domain.TransactionRepository) *TransactionWorker {
	return &TransactionWorker{eventBus, userRepo, transactionRepo}
}

// StartWorker Starts listening for transaction events
func (w *TransactionWorker) StartWorker() {
	eventChan := w.eventBus.Subscribe(EventTypeTransfer)

	go func() {
		for event := range eventChan {
			trans, ok := event.(TransferParam)
			if !ok {
				log.Println("Invalid event received")
				continue
			}
			w.processTransfer(trans)
		}
	}()
}

func (w *TransactionWorker) processTransfer(trans TransferParam) {
	sender, err := w.userRepository.GetUserByID(trans.TransferInfo.UserID)
	if err != nil {
		log.Printf("User not found: %v", err)
		return
	}

	target, err := w.userRepository.GetUserByID(trans.TargetID)
	if err != nil {
		log.Printf("Target user not found: %v", err)
		return
	}

	sender.Balance -= trans.TransferInfo.Amount
	err = w.userRepository.UpdateUser(context.TODO(), sender)
	if err != nil {
		log.Printf("User update error: %v", err)
		return
	}

	target.Balance += trans.TransferInfo.Amount
	err = w.userRepository.UpdateUser(context.TODO(), target)
	if err != nil {
		log.Printf("Target user update error: %v", err)
		return
	}

	transInfo, err := w.transactionRepo.GetTransactionByID(context.Background(), trans.TransferInfo.ID)
	if err != nil {
		log.Printf("Failed to get transaction info: %v", err)
		return
	}
	transInfo.Status = "success"
	err = w.transactionRepo.UpdateTransaction(context.TODO(), transInfo)
	if err != nil {
		log.Printf("Failed to update transaction info: %v", err)
		return
	}

	log.Printf("Transaction info %s updated!", transInfo.ID)
}
