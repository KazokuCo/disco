package conversion

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	currencyNamePattern string         = `[a-zA-Z$€£₤₽¥円]*`
	currencyRegex       *regexp.Regexp = regexp.MustCompile(`(?i)` +
		`(` + currencyNamePattern + `) *([\d\.\,]+) *(` + currencyNamePattern + `)` +
		`(?: *(in|to|as)? *` +
		`(` + currencyNamePattern + `)` +
		`)?`)
	currencyAliases map[string]string = map[string]string{
		"$":    "USD",
		"A$":   "AUD",
		"AU$":  "AUD",
		"AUS$": "AUD",
		"C$":   "CAD",
		"CA$":  "CAD",
		"CAN$": "CAD",
		"HK$":  "HKD",
		"NZ$":  "NZD",
		"S$":   "SGD",
		"US$":  "USD",
		"R$":   "BRL",
		"€":    "EUR",
		"KR":   "SEK",
		"DKR":  "DKK",
		"NKR":  "NOK",
		"£":    "GBP",
		"₤":    "GBP",
		"₽":    "RUB",
		"¥":    "JPY",
		"円":    "JPY",
		"YEN":  "JPY",
	}
)

type CurrencyRates struct {
	Base  string
	Rates map[string]float64
}

func (r CurrencyRates) Convert(val float64, from, to string) (float64, error) {
	if from != r.Base {
		rate, ok := r.Rates[from]
		if !ok {
			return val, errors.New("No from rate found")
		}
		val /= rate
	}
	if to != r.Base {
		rate, ok := r.Rates[to]
		if !ok {
			return val, errors.New("No to rate found")
		}
		val *= rate
	}
	return val, nil
}

func ParseCurrencyString(s string) (val float64, from, to string, ok bool) {
	matches := currencyRegex.FindAllStringSubmatch(s, -1)
	if matches != nil {
		return ParseCurrency(matches[0])
	}
	return 0, "", "", false
}

func ParseCurrency(m []string) (val float64, from, to string, ok bool) {
	val, err := strconv.ParseFloat(strings.Replace(m[2], ",", ".", -1), 64)
	if err != nil {
		return 0, "", "", false
	}

	// Expressions like `$10 in EUR` get misparsed (in=suffix)
	prefix := m[1]
	suffix := m[3]
	if suffix == "in" || suffix == "to" || suffix == "as" {
		suffix = ""
	}

	// Prefer suffix to prefix
	if suffix != "" {
		from = suffix
	} else if prefix != "" {
		from = prefix
	}

	to = m[5]

	// A value is meaningless without a unit
	if val != 0 && from == "" {
		return val, from, to, false
	}
	return val, from, to, true
}

func ResolveCurrency(c string) string {
	c = strings.ToUpper(c)
	if c2, ok := currencyAliases[c]; ok {
		return c2
	}
	return c
}

func FetchRates() (rates CurrencyRates, err error) {
	res, err := http.Get("https://api.fixer.io/latest")
	if err != nil {
		return rates, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return rates, errors.New(fmt.Sprintf("Status != 200: %v", res.StatusCode))
	}

	body, _ := ioutil.ReadAll(res.Body)
	if err = json.Unmarshal(body, &rates); err != nil {
		return rates, err
	}

	return rates, nil
}

func (j *Job) HandleCurrency(s *discordgo.Session, msg *discordgo.Message, matches [][]string) {
	val, from, to, ok := ParseCurrency(matches[0])
	if !ok {
		return
	}
	from = ResolveCurrency(from)
	to = ResolveCurrency(to)
	log.WithFields(log.Fields{"val": val, "from": from, "to": to}).Info("Looks like currency")

	rates, err := FetchRates()
	if err != nil {
		log.WithError(err).Error("Conversion: Couldn't fetch currency rates")
		return
	}

	if to != "" {
		val2, err := rates.Convert(val, from, to)
		if err != nil {
			log.WithError(err).Error("Conversion: Couldn't convert currency")
			return
		}
		fromText := fmt.Sprintf("%.2f %s", val, from)
		toText := fmt.Sprintf("%.2f %s", val2, to)
		text := fmt.Sprintf(j.Lines.Currency, fromText, toText)
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("<@%s> %s", msg.Author.ID, text))
	} else {
		usd, err := rates.Convert(val, from, "USD")
		eur, _ := rates.Convert(val, from, "EUR")
		gbp, _ := rates.Convert(val, from, "GBP")
		jpy, _ := rates.Convert(val, from, "JPY")
		if err != nil {
			log.WithError(err).Error("Conversion: Couldn't convert currency")
			return
		}
		fromText := fmt.Sprintf("%.2f %s", val, from)
		text := fmt.Sprintf(j.Lines.CurrencyMulti, fromText, fmt.Sprintf("%.2f", usd), fmt.Sprintf("%.2f", eur), fmt.Sprintf("%.2f", gbp), fmt.Sprintf("%.2f", jpy))
		s.ChannelMessageSend(msg.ChannelID, text)
	}
}
