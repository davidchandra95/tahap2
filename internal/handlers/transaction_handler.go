package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"tahap2/internal/domain"
	"tahap2/internal/middlewares"
	"time"
)

type TransactionHandler struct {
	transService domain.TransactionService
}

func NewTransactionHandler(transService domain.TransactionService) *TransactionHandler {
	return &TransactionHandler{transService: transService}
}

func (h *TransactionHandler) TopupHandler(c echo.Context) error {
	userID := c.Get(middlewares.UserIDKey).(uuid.UUID)

	var req struct {
		Amount int64 `json:"amount"`
	}

	if err := c.Bind(&req); err != nil || req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid amount"})
	}

	transInfo, err := h.transService.ProcessTopUp(userID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": fmt.Sprintf("topup failed. : %s", err.Error())})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"result": toTransactionResponse(transInfo),
	})
}

func (h *TransactionHandler) PaymentHandler(c echo.Context) error {
	userID := c.Get(middlewares.UserIDKey).(uuid.UUID)

	var req struct {
		Amount  int64  `json:"amount"`
		Remarks string `json:"remarks"`
	}

	if err := c.Bind(&req); err != nil || req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid amount"})
	}

	transInfo, err := h.transService.ProcessPayment(userID, req.Amount, req.Remarks)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": fmt.Sprintf("payment failed. : %s", err.Error())})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"result": toTransactionResponse(transInfo),
	})
}

func (h *TransactionHandler) TransferHandler(c echo.Context) error {
	userID := c.Get(middlewares.UserIDKey).(uuid.UUID)

	var req struct {
		TargetUserID uuid.UUID `json:"target_user"`
		Amount       int64     `json:"amount"`
		Remarks      string    `json:"remarks"`
	}

	if err := c.Bind(&req); err != nil || req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "invalid amount"})
	}

	transInfo, err := h.transService.ProcessTransfer(userID, req.TargetUserID, req.Amount, req.Remarks)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": fmt.Sprintf("transfer failed. : %s", err.Error())})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"result": toTransactionResponse(transInfo),
	})
}

func (h *TransactionHandler) GetAllTransactions(c echo.Context) error {
	userID := c.Get(middlewares.UserIDKey).(uuid.UUID)

	transactions, err := h.transService.GetAllTransactions(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": fmt.Sprintf("get transactions failed. : %s", err.Error())})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status": "success",
		"result": toTransactionsDetailResponse(transactions),
	})
}

func toTransactionsDetailResponse(trans []*domain.Transaction) []*TransactionDetailsResponse {
	result := make([]*TransactionDetailsResponse, len(trans))
	for i, tran := range trans {
		result[i] = &TransactionDetailsResponse{
			TransactionID:   tran.ID.String(),
			UserID:          tran.UserID.String(),
			TransactionType: tran.TransactionType,
			Amount:          tran.Amount,
			Remarks:         tran.Remark,
			BalanceBefore:   tran.BalanceBefore,
			BalanceAfter:    tran.BalanceAfter,
			Status:          tran.Status,
			CreatedAt:       tran.CreatedAt.Format(time.DateTime),
		}
	}

	return result
}

func toTransactionResponse(src domain.Transaction) TransactionResponse {
	return TransactionResponse{
		TransactionID: src.ID.String(),
		Amount:        src.Amount,
		BalanceBefore: src.BalanceBefore,
		BalanceAfter:  src.BalanceAfter,
		Status:        src.Status,
		CreatedAt:     src.CreatedAt.Format(time.DateTime),
	}
}

type TransactionDetailsResponse struct {
	TransactionID   string `json:"transaction_id"`
	UserID          string `json:"user_id"`
	TransactionType string `json:"transaction_type"`
	Amount          int64  `json:"amount"`
	Remarks         string `json:"remarks"`
	BalanceBefore   int64  `json:"balance_before"`
	BalanceAfter    int64  `json:"balance_after"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
}

type TransactionResponse struct {
	TransactionID string `json:"transaction_id"`
	Amount        int64  `json:"amount"`
	BalanceBefore int64  `json:"balance_before"`
	BalanceAfter  int64  `json:"balance_after"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}
