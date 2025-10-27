package entity

type Currency string

const (
	CurrencyUSD Currency = "USD" // Dólar dos Estados Unidos
	CurrencyEUR Currency = "EUR" // Euro da União Europeia
	CurrencyCHF Currency = "CHF" // Franco Suíço
	CurrencyGBP Currency = "GBP" // Libra Esterlina
	CurrencyBRL Currency = "BRL" // Real Brasileiro
)

func (c Currency) String() string {
	return string(c)
}
