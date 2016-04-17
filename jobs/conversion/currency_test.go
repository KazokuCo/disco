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

func TestParseCurrencyConvertPrefix(t *testing.T) {
	val, from, to, ok := ParseCurrencyString("EUR 10 in USD")
	if !ok || val != 10 || from != "EUR" || to != "USD" {
		t.Error(ok, val, from, to)
	}
}

func TestConvertCurrencyBaseToBase(t *testing.T) {
	rates := CurrencyRates{Base: "EUR", Rates: map[string]float64{"USD": 2}}
	v, err := rates.Convert(10, "EUR", "EUR")
	if err != nil {
		t.Error(err)
	}
	if v != 10 {
		t.Fail()
	}
}

func TestConvertCurrencyBaseToOther(t *testing.T) {
	rates := CurrencyRates{Base: "EUR", Rates: map[string]float64{"USD": 2}}
	v, err := rates.Convert(10, "EUR", "USD")
	if err != nil {
		t.Error(err)
	}
	if v != 20 {
		t.Fail()
	}
}

func TestConvertCurrencyOtherToBase(t *testing.T) {
	rates := CurrencyRates{Base: "EUR", Rates: map[string]float64{"USD": 2}}
	v, err := rates.Convert(10, "USD", "EUR")
	if err != nil {
		t.Error(err)
	}
	if v != 5 {
		t.Fail()
	}
}

func TestConvertCurrencyFromUnknown(t *testing.T) {
	rates := CurrencyRates{Base: "EUR", Rates: map[string]float64{"USD": 2}}
	_, err := rates.Convert(10, "butts", "EUR")
	if err == nil {
		t.Fail()
	}
}

func TestConvertCurrencyToUnknown(t *testing.T) {
	rates := CurrencyRates{Base: "EUR", Rates: map[string]float64{"USD": 2}}
	_, err := rates.Convert(10, "EUR", "butts")
	if err == nil {
		t.Fail()
	}
}

func TestConvertCurrencyFromAndToUnknown(t *testing.T) {
	rates := CurrencyRates{Base: "EUR", Rates: map[string]float64{"USD": 2}}
	_, err := rates.Convert(10, "butts", "florps")
	if err == nil {
		t.Fail()
	}
}

func TestResolveCurrency(t *testing.T) {
	if ResolveCurrency("yen") != "JPY" {
		t.Fail()
	}
}

func TestResolveCurrencyCaseInsensitive(t *testing.T) {
	if ResolveCurrency("YEN") != "JPY" {
		t.Fail()
	}
}

func TestResolveCurrencyUnaliased(t *testing.T) {
	if ResolveCurrency("USD") != "USD" {
		t.Fail()
	}
}

func TestResolveCurrencyUnaliasedCaseInsensitive(t *testing.T) {
	if ResolveCurrency("usd") != "USD" {
		t.Fail()
	}
}
