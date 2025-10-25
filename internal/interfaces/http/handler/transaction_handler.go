package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

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
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.transactionUseCase.RecordTransaction(c.Request.Context(), user.ID, request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *TransactionHandler) Update(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request dto.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.transactionUseCase.UpdateTransaction(c.Request.Context(), user.ID, c.Param("id"), request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) List(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	from, to := parseDateRange(c.Query("from"), c.Query("to"))
	response, err := h.transactionUseCase.ListTransactions(c.Request.Context(), user.ID, from, to)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) AttachReceipt(c *gin.Context) {
	user, ok := middleware.GetUserContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	opened, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer opened.Close()

	data, err := io.ReadAll(opened)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load file"})
		return
	}

	response, err := h.transactionUseCase.AttachReceipt(c.Request.Context(), user.ID, c.Param("id"), file.Filename, file.Header.Get("Content-Type"), data)
	if err != nil {
		respondError(c, err)
		return
	}

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
