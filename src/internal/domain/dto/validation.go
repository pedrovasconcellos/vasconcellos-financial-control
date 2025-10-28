package dto

import (
	"github.com/vasconcellos/financial-control/src/internal/domain/entity"
)

// CurrencyValidationTag retorna a tag de validação Gin para moedas suportadas
func CurrencyValidationTag() string {
	return "oneof=" + entity.SupportedCurrencyCodesString()
}

// RequiredCurrencyValidationTag retorna a tag de validação Gin para moedas obrigatórias
func RequiredCurrencyValidationTag() string {
	return "required," + CurrencyValidationTag()
}

// OptionalCurrencyValidationTag retorna a tag de validação Gin para moedas opcionais
func OptionalCurrencyValidationTag() string {
	return "omitempty," + CurrencyValidationTag()
}
