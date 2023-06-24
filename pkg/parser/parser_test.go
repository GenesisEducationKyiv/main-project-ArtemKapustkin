package parser

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestGetExchangeRate(t *testing.T) {
	if err := godotenv.Load("../../.env.test"); err != nil {
		t.Fatal("Failed to load .env.test file")
	}
	binanceParser := NewBinanceCryptoParser(os.Getenv("BASE_URL"))
	rate, err := binanceParser.GetExchangeRate("BTC", "UAH")
	if err != nil {
		t.Errorf("failure occurs while parsing exchange rate: %v", err)
	}

	if rate < 0 {
		t.Errorf("exchange rate is negative value: %v", rate)
	}
}

func TestGetExchangeRateFault(t *testing.T) {
	if err := godotenv.Load("../../.env.test"); err != nil {
		t.Fatal("Failed to load .env.test file")
	}
	binanceParser := NewBinanceCryptoParser(os.Getenv("BASE_URL"))
	rate, err := binanceParser.GetExchangeRate("BTC", "BTC")
	if err == errors.New("invalid syntax") {
		t.Errorf("invalid syntax failure occurs while parsing exchange rate: %v", err)
	}

	if rate < 0 {
		t.Errorf("exchange rate is negative value: %v", rate)
	}
}
