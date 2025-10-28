package entity

import (
	"testing"
)

func TestCurrencyValidation(t *testing.T) {
	tests := []struct {
		name         string
		currencyCode string
		expected     bool
	}{
		{"USD válida", "USD", true},
		{"EUR válida", "EUR", true},
		{"CHF válida", "CHF", true},
		{"GBP válida", "GBP", true},
		{"BRL válida", "BRL", true},
		{"CNY inválida", "CNY", false},
		{"JPY inválida", "JPY", false},
		{"CAD inválida", "CAD", false},
		{"AUD inválida", "AUD", false},
		{"String vazia", "", false},
		{"Moeda inexistente", "XYZ", false},
		{"Case sensitive", "usd", false},
		{"Case sensitive", "Usd", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidCurrency(tt.currencyCode)
			if result != tt.expected {
				t.Errorf("IsValidCurrency(%s) = %v, expected %v", tt.currencyCode, result, tt.expected)
			}
		})
	}
}

func TestValidateCurrency(t *testing.T) {
	tests := []struct {
		name         string
		currencyCode string
		shouldError  bool
	}{
		{"USD válida", "USD", false},
		{"EUR válida", "EUR", false},
		{"CHF válida", "CHF", false},
		{"GBP válida", "GBP", false},
		{"BRL válida", "BRL", false},
		{"CNY inválida", "CNY", true},
		{"JPY inválida", "JPY", true},
		{"CAD inválida", "CAD", true},
		{"AUD inválida", "AUD", true},
		{"String vazia", "", true},
		{"Moeda inexistente", "XYZ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCurrency(tt.currencyCode)
			hasError := err != nil

			if hasError != tt.shouldError {
				t.Errorf("ValidateCurrency(%s) error = %v, expected error = %v",
					tt.currencyCode, hasError, tt.shouldError)
			}

			if tt.shouldError && err != nil {
				// Verifica se a mensagem de erro contém informações úteis
				if !contains(err.Error(), tt.currencyCode) {
					t.Errorf("Error message should contain currency code '%s', got: %s",
						tt.currencyCode, err.Error())
				}
			}
		})
	}
}

func TestSupportedCurrencyCodes(t *testing.T) {
	codes := SupportedCurrencyCodes()
	expectedCodes := []string{"USD", "EUR", "CHF", "GBP", "BRL"}

	if len(codes) != len(expectedCodes) {
		t.Errorf("Expected %d supported currencies, got %d", len(expectedCodes), len(codes))
	}

	for _, expected := range expectedCodes {
		found := false
		for _, code := range codes {
			if code == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected currency code %s not found in supported currencies", expected)
		}
	}
}

func TestSupportedCurrencyCodesString(t *testing.T) {
	codesString := SupportedCurrencyCodesString()
	expected := "USD EUR CHF GBP BRL"

	if codesString != expected {
		t.Errorf("Expected '%s', got '%s'", expected, codesString)
	}
}

func TestGetCurrencyName(t *testing.T) {
	tests := []struct {
		currency Currency
		expected string
	}{
		{CurrencyUSD, "United States Dollar"},
		{CurrencyEUR, "Euro"},
		{CurrencyCHF, "Swiss Franc"},
		{CurrencyGBP, "Pound Sterling"},
		{CurrencyBRL, "Brazilian Real"},
		{Currency("XYZ"), "Unknown Currency"},
	}

	for _, tt := range tests {
		t.Run(string(tt.currency), func(t *testing.T) {
			result := tt.currency.GetCurrencyName()
			if result != tt.expected {
				t.Errorf("GetCurrencyName() = %s, expected %s", result, tt.expected)
			}
		})
	}
}

// Função auxiliar para verificar se uma string contém outra
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
