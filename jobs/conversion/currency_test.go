package conversion

import (
	"testing"
)

func TestParseCurrencyNumberOnly(t *testing.T) {
	val, _, _, ok := ParseCurrencyString("10")
	if ok || val != 10 {
		t.Error(ok, val)
	}
}

func TestParseCurrencyDecimal(t *testing.T) {
	val, _, _, ok := ParseCurrencyString("10.50 EUR")
	if !ok || val != 10.5 {
		t.Error(ok, val)
	}
}

func TestParseCurrencyDecimalComma(t *testing.T) {
	val, _, _, ok := ParseCurrencyString("10,50 EUR")
	if !ok || val != 10.5 {
		t.Error(ok, val)
	}
}

func TestParseCurrencyPrefix(t *testing.T) {
	val, from, _, ok := ParseCurrencyString("EUR 10")
	if !ok || val != 10 || from != "EUR" {
		t.Error(ok, val, from)
	}
}

func TestParseCurrencySuffix(t *testing.T) {
	val, from, _, ok := ParseCurrencyString("10 EUR")
	if !ok || val != 10 || from != "EUR" {
		t.Error(ok, val, from)
	}
}

func TestParseCurrencyConvert(t *testing.T) {
	val, from, to, ok := ParseCurrencyString("10 EUR in USD")
	if !ok || val != 10 || from != "EUR" || to != "USD" {
		t.Error(ok, val, from, to)
	}
}
