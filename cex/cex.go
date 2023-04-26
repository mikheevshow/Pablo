package cex

const (
	Binance string = "binance"
	Bybit	string = "bybit"
	Okx		string = "okx"
)

type Creds struct {
	ApiKey string
	Secret string
}