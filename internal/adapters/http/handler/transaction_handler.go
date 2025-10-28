package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/vasconcellos/financial-control/internal/adapters/http/middleware"
	"github.com/vasconcellos/financial-control/internal/domain/dto"
	domainErrors "github.com/vasconcellos/financial-control/internal/domain/errors"
	"github.com/vasconcellos/financial-control/internal/usecase"
)

type TransactionHandler struct {
	transactionUseCase *usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUseCase *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{transactionUseCase: transactionUseCase}
}

// Create
// @Summary Create a new transaction
// @Description Cria uma nova transação financeira e publica evento para processamento assíncrono de budgets
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTransactionRequest true "Dados da transação"
// @Success 201 {object} dto.TransactionResponse "Transação criada com sucesso"
// @Failure 400 {object} ErrorResponse "Dados inválidos"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Failure 500 {object} ErrorResponse "Erro interno"
// @Router /transactions [post]
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

	log.Info("recording transaction",
		zap.String("user_id", user.ID),
		zap.String("account_id", request.AccountID),
		zap.String("category_id", request.CategoryID),
		zap.Float64("amount", request.Amount),
		zap.String("currency", request.Currency))

	response, err := h.transactionUseCase.RecordTransaction(c.Request.Context(), user.ID, request)
	if err != nil {
		log.Error("failed to record transaction",
			zap.Error(err),
			zap.String("account_id", request.AccountID),
			zap.String("category_id", request.CategoryID),
			zap.Float64("amount", request.Amount),
			zap.String("user_id", user.ID))

		respondError(c, err)
		return
	}

	log.Info("transaction recorded", zap.String("transaction_id", response.ID))
	c.JSON(http.StatusCreated, response)
}

// Update
// @Summary Update a transaction
// @Description Atualiza uma transação existente
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da transação"
// @Param request body dto.UpdateTransactionRequest true "Dados atualizados"
// @Success 200 {object} dto.TransactionResponse "Transação atualizada"
// @Failure 400 {object} ErrorResponse "Dados inválidos"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Failure 404 {object} ErrorResponse "Transação não encontrada"
// @Router /transactions/{id} [patch]
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

// List
// @Summary List transactions
// @Description Lista transações do usuário com filtros opcionais de data e paginação
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "Data inicial (ISO 8601)"
// @Param to query string false "Data final (ISO 8601)"
// @Param limit query int false "Número máximo de resultados (default: 100, max: 200)"
// @Param offset query int false "Número de resultados para pular (default: 0)"
// @Success 200 {array} dto.TransactionResponse "Lista de transações"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Failure 500 {object} ErrorResponse "Erro interno"
// @Router /transactions [get]
func (h *TransactionHandler) List(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	user, ok := middleware.GetUserContext(c)
	if !ok {
		log.Warn("unauthorized transaction list attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	from, to := parseDateRange(c.Query("from"), c.Query("to"))
	limit, offset, err := parsePagination(c.Query("limit"), c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info("listing transactions", zap.String("user_id", user.ID), zap.Time("from", from), zap.Time("to", to), zap.Int64("limit", limit), zap.Int64("offset", offset))
	response, err := h.transactionUseCase.ListTransactions(c.Request.Context(), user.ID, from, to, limit, offset)
	if err != nil {
		log.Error("failed to list transactions", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// AttachReceipt
// @Summary Attach receipt to transaction
// @Description Envia e criptografa um recibo para a transação. Arquivo é armazenado no S3 com AES-256
// @Tags transactions
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID da transação"
// @Param file formData file true "Arquivo do recibo (PDF, PNG, JPG - max 5MB)"
// @Success 200 {object} dto.TransactionResponse "Recibo anexado com sucesso"
// @Failure 400 {object} ErrorResponse "Arquivo inválido ou muito grande"
// @Failure 401 {object} ErrorResponse "Não autenticado"
// @Failure 404 {object} ErrorResponse "Transação não encontrada"
// @Failure 500 {object} ErrorResponse "Erro interno"
// @Router /transactions/{id}/receipt [post]
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

	const maxReceiptSizeBytes = usecase.MaxReceiptSizeBytes
	if file.Size > 0 && file.Size > maxReceiptSizeBytes {
		log.Warn("receipt exceeds size limit", zap.Int64("size", file.Size), zap.Int64("limit_bytes", maxReceiptSizeBytes))
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

	response, err := h.transactionUseCase.AttachReceipt(
		c.Request.Context(),
		user.ID,
		transactionID,
		file.Filename,
		file.Header.Get("Content-Type"),
		opened,
	)
	if err != nil {
		if err == domainErrors.ErrPayloadTooLarge {
			log.Warn("receipt exceeds size limit during processing", zap.String("transaction_id", transactionID), zap.Int64("limit_bytes", maxReceiptSizeBytes))
			c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 5MB)"})
			return
		}
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
