package entity

import (
	"fmt"
	"strings"
)

type Currency string

const (
	CurrencyUSD Currency = "USD" // Dólar dos Estados Unidos
	CurrencyEUR Currency = "EUR" // Euro da União Europeia
	CurrencyCHF Currency = "CHF" // Franco Suíço
	CurrencyGBP Currency = "GBP" // Libra Esterlina
	CurrencyBRL Currency = "BRL" // Real Brasileiro
)

// SupportedCurrencies lista todas as moedas suportadas pelo sistema
var SupportedCurrencies = []Currency{
	CurrencyUSD,
	CurrencyEUR,
	CurrencyCHF,
	CurrencyGBP,
	CurrencyBRL,
}

// SupportedCurrencyCodes retorna os códigos das moedas suportadas como slice de strings
func SupportedCurrencyCodes() []string {
	codes := make([]string, len(SupportedCurrencies))
	for i, currency := range SupportedCurrencies {
		codes[i] = string(currency)
	}
	return codes
}

// SupportedCurrencyCodesString retorna os códigos das moedas suportadas como string separada por espaços
// Usado para validação Gin binding (oneof=USD EUR CHF GBP BRL)
func SupportedCurrencyCodesString() string {
	return strings.Join(SupportedCurrencyCodes(), " ")
}

// IsValidCurrency verifica se uma string representa uma moeda válida
func IsValidCurrency(currencyCode string) bool {
	currency := Currency(currencyCode)
	for _, supported := range SupportedCurrencies {
		if currency == supported {
			return true
		}
	}
	return false
}

// ValidateCurrency valida se uma string representa uma moeda válida e retorna erro se inválida
func ValidateCurrency(currencyCode string) error {
	if !IsValidCurrency(currencyCode) {
		return fmt.Errorf("currency '%s' is not supported. Supported currencies: %s",
			currencyCode, SupportedCurrencyCodesString())
	}
	return nil
}

func (c Currency) String() string {
	return string(c)
}

// GetCurrencyName retorna o nome completo da moeda
func (c Currency) GetCurrencyName() string {
	switch c {
	case CurrencyUSD:
		return "United States Dollar"
	case CurrencyEUR:
		return "Euro"
	case CurrencyCHF:
		return "Swiss Franc"
	case CurrencyGBP:
		return "Pound Sterling"
	case CurrencyBRL:
		return "Brazilian Real"
	default:
		return "Unknown Currency"
	}
}
