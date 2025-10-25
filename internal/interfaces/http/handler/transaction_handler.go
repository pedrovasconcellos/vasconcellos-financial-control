package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vasconcellos/finance-control/internal/domain/dto"
	"github.com/vasconcellos/finance-control/internal/interfaces/http/middleware"
	"github.com/vasconcellos/finance-control/internal/usecase"
)

type TransactionHandler struct {
	transactionUseCase *usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{transactionUseCase: transactionUseCase}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized transaction creation attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn("invalid transaction payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("recording transaction", zap.String("user_id", user.ID), zap.String("account_id", request.AccountID), zap.String("category_id", request.CategoryID))
	response, err := h.transactionUseCase.RecordTransaction(c.Request.Context(), user.ID, request)
	if err != nil {
		log.Error("failed to record transaction", zap.Error(err))
		respondError(c, err)
		return
	}

	log.Info("transaction recorded", zap.String("transaction_id", response.ID))
	c.JSON(http.StatusCreated, response)
}

func (h *TransactionHandler) Update(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized transaction update attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn("invalid transaction update payload", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactionID := c.Param("id")
	log.Info("updating transaction", zap.String("transaction_id", transactionID), zap.String("user_id", user.ID))
	response, err := h.transactionUseCase.UpdateTransaction(c.Request.Context(), user.ID, c.Param("id"), request)
	if err != nil {
		log.Error("failed to update transaction", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) List(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized transaction list attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	from, to := parseDateRange(c.Query("from"), c.Query("to"))
	log.Info("listing transactions", zap.String("user_id", user.ID), zap.Time("from", from), zap.Time("to", to))
	response, err := h.transactionUseCase.ListTransactions(c.Request.Context(), user.ID, from, to)
	if err != nil {
		log.Error("failed to list transactions", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) AttachReceipt(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized receipt upload attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	transactionID := c.Param("id")
	log.Info("attaching receipt", zap.String("transaction_id", transactionID), zap.String("user_id", user.ID))
	file, err := c.FormFile("file")
	if err != nil {
		log.Warn("receipt payload missing", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	const maxReceiptSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxReceiptSize {
		log.Warn("receipt exceeds size limit", zap.Int64("size", file.Size), zap.Int("limit_bytes", maxReceiptSize))
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 5MB)"})
		return
	}

	opened, err := file.Open()
	if err != nil {
		log.Error("failed to open receipt", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer opened.Close()

	limitedReader := io.LimitReader(opened, int64(maxReceiptSize)+1)
	response, err := h.transactionUseCase.AttachReceipt(c.Request.Context(), user.ID, transactionID, file.Filename, file.Header.Get("Content-Type"), limitedReader)
	if err != nil {
		log.Error("failed to attach receipt", zap.Error(err))
		respondError(c, err)
		return
	}

	log.Info("receipt attached", zap.String("transaction_id", transactionID))
	c.JSON(http.StatusOK, response)
}

func parseDateRange(fromRaw, toRaw string) (time.Time, time.Time) {
	const layout = time.RFC3339
	now := time.Now().UTC()
	if fromRaw == "" {
		fromRaw = now.AddDate(0, 0, -30).Format(layout)
	}
	if toRaw == "" {
		toRaw = now.Format(layout)
	}

	from, err := time.Parse(layout, fromRaw)
	if err != nil {
		from = now.AddDate(0, 0, -30)
	}
	to, err := time.Parse(layout, toRaw)
	if err != nil {
		to = now
	}
	return from, to
}
