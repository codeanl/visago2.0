package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const exchangeRateSource = "Frankfurter"

type exchangeCountryItem struct {
	CountryCode    string `json:"countryCode"`
	CountryName    string `json:"countryName"`
	EnglishName    string `json:"englishName"`
	CurrencyCode   string `json:"currencyCode"`
	CurrencyName   string `json:"currencyName"`
	CurrencySymbol string `json:"currencySymbol"`
	Supported      bool   `json:"supported"`
}

type exchangeQuoteItem struct {
	Date            string  `json:"date"`
	FromCurrency    string  `json:"fromCurrency"`
	ToCurrency      string  `json:"toCurrency"`
	Amount          float64 `json:"amount"`
	ConvertedAmount float64 `json:"convertedAmount"`
	Rate            float64 `json:"rate"`
	Source          string  `json:"source"`
}

type exchangeTrendPoint struct {
	Date string  `json:"date"`
	Rate float64 `json:"rate"`
}

type exchangeTrendItem struct {
	FromCurrency string               `json:"fromCurrency"`
	ToCurrency   string               `json:"toCurrency"`
	Source       string               `json:"source"`
	Points       []exchangeTrendPoint `json:"points"`
}

type frankfurterQuoteResponse struct {
	Date    string  `json:"date"`
	Base    string  `json:"base"`
	Quote   string  `json:"quote"`
	Rate    float64 `json:"rate"`
	Message string  `json:"message"`
}

//go:embed data/exchange_countries.json
var exchangeCountriesJSON []byte

var (
	exchangeCountriesOnce       sync.Once
	exchangeCountriesCache      []exchangeCountryItem
	exchangeSupportedCurrencies map[string]struct{}
	exchangeCountriesErr        error
)

func (s *appServer) handleExchangeCountries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	items, err := loadExchangeCountries()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	supportedOnly := parseTruthy(r.URL.Query().Get("supported"))
	filtered := make([]exchangeCountryItem, 0, len(items))
	for _, item := range items {
		if supportedOnly && !item.Supported {
			continue
		}
		if q != "" && !matchExchangeCountry(item, q) {
			continue
		}
		filtered = append(filtered, item)
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: filtered})
}

func (s *appServer) handleExchangeQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	from := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("from")))
	to := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("to")))
	if from == "" || to == "" {
		writeError(w, http.StatusBadRequest, errors.New("from and to are required"))
		return
	}
	amount := 1.0
	if raw := strings.TrimSpace(r.URL.Query().Get("amount")); raw != "" {
		parsed, err := strconv.ParseFloat(raw, 64)
		if err != nil || parsed < 0 {
			writeError(w, http.StatusBadRequest, errors.New("amount must be a valid number"))
			return
		}
		amount = parsed
	}
	if from == to {
		writeJSON(w, http.StatusOK, apiResponse{
			Message: "ok",
			Data: exchangeQuoteItem{
				Date:            time.Now().Format("2006-01-02"),
				FromCurrency:    from,
				ToCurrency:      to,
				Amount:          amount,
				ConvertedAmount: amount,
				Rate:            1,
				Source:          exchangeRateSource,
			},
		})
		return
	}
	if !isExchangeCurrencySupported(from) {
		writeError(w, http.StatusBadRequest, fmt.Errorf("暂不支持 %s 的实时汇率查询", from))
		return
	}
	if !isExchangeCurrencySupported(to) {
		writeError(w, http.StatusBadRequest, fmt.Errorf("暂不支持 %s 的实时汇率查询", to))
		return
	}
	quote, err := fetchExchangeQuote(r.Context(), from, to)
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{
		Message: "ok",
		Data: exchangeQuoteItem{
			Date:            quote.Date,
			FromCurrency:    from,
			ToCurrency:      to,
			Amount:          amount,
			ConvertedAmount: amount * quote.Rate,
			Rate:            quote.Rate,
			Source:          exchangeRateSource,
		},
	})
}

func (s *appServer) handleExchangeTrend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	from := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("from")))
	to := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("to")))
	if from == "" || to == "" {
		writeError(w, http.StatusBadRequest, errors.New("from and to are required"))
		return
	}
	days := 7
	if raw := strings.TrimSpace(r.URL.Query().Get("days")); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 2 || parsed > 30 {
			writeError(w, http.StatusBadRequest, errors.New("days must be between 2 and 30"))
			return
		}
		days = parsed
	}
	if from == to {
		points := make([]exchangeTrendPoint, 0, days)
		now := time.Now()
		for i := days - 1; i >= 0; i-- {
			points = append(points, exchangeTrendPoint{
				Date: now.AddDate(0, 0, -i).Format("2006-01-02"),
				Rate: 1,
			})
		}
		writeJSON(w, http.StatusOK, apiResponse{
			Message: "ok",
			Data: exchangeTrendItem{
				FromCurrency: from,
				ToCurrency:   to,
				Source:       exchangeRateSource,
				Points:       points,
			},
		})
		return
	}
	if !isExchangeCurrencySupported(from) {
		writeError(w, http.StatusBadRequest, fmt.Errorf("暂不支持 %s 的历史汇率查询", from))
		return
	}
	if !isExchangeCurrencySupported(to) {
		writeError(w, http.StatusBadRequest, fmt.Errorf("暂不支持 %s 的历史汇率查询", to))
		return
	}
	points, err := fetchExchangeTrend(r.Context(), from, to, days)
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{
		Message: "ok",
		Data: exchangeTrendItem{
			FromCurrency: from,
			ToCurrency:   to,
			Source:       exchangeRateSource,
			Points:       points,
		},
	})
}

func loadExchangeCountries() ([]exchangeCountryItem, error) {
	exchangeCountriesOnce.Do(func() {
		if err := json.Unmarshal(exchangeCountriesJSON, &exchangeCountriesCache); err != nil {
			exchangeCountriesErr = fmt.Errorf("decode exchange countries: %w", err)
			return
		}
		exchangeSupportedCurrencies = make(map[string]struct{}, len(exchangeCountriesCache))
		for _, item := range exchangeCountriesCache {
			if item.Supported {
				exchangeSupportedCurrencies[strings.ToUpper(strings.TrimSpace(item.CurrencyCode))] = struct{}{}
			}
		}
	})
	if exchangeCountriesErr != nil {
		return nil, exchangeCountriesErr
	}
	out := make([]exchangeCountryItem, len(exchangeCountriesCache))
	copy(out, exchangeCountriesCache)
	return out, nil
}

func isExchangeCurrencySupported(code string) bool {
	if _, err := loadExchangeCountries(); err != nil {
		return false
	}
	_, ok := exchangeSupportedCurrencies[strings.ToUpper(strings.TrimSpace(code))]
	return ok
}

func matchExchangeCountry(item exchangeCountryItem, keyword string) bool {
	if strings.TrimSpace(keyword) == "" {
		return true
	}
	q := strings.ToLower(strings.TrimSpace(keyword))
	return strings.Contains(strings.ToLower(item.CountryCode), q) ||
		strings.Contains(strings.ToLower(item.CountryName), q) ||
		strings.Contains(strings.ToLower(item.EnglishName), q) ||
		strings.Contains(strings.ToLower(item.CurrencyCode), q) ||
		strings.Contains(strings.ToLower(item.CurrencyName), q)
}

func parseTruthy(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

func fetchExchangeQuote(ctx context.Context, from, to string) (frankfurterQuoteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	url := fmt.Sprintf("https://api.frankfurter.dev/v2/rate/%s/%s", strings.ToUpper(from), strings.ToUpper(to))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return frankfurterQuoteResponse{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return frankfurterQuoteResponse{}, fmt.Errorf("汇率服务请求失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr frankfurterQuoteResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && strings.TrimSpace(apiErr.Message) != "" {
			return frankfurterQuoteResponse{}, errors.New(strings.TrimSpace(apiErr.Message))
		}
		return frankfurterQuoteResponse{}, fmt.Errorf("汇率服务返回异常状态: %d", resp.StatusCode)
	}
	var out frankfurterQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return frankfurterQuoteResponse{}, fmt.Errorf("解析汇率结果失败: %w", err)
	}
	if out.Rate <= 0 {
		return frankfurterQuoteResponse{}, errors.New("暂未获取到有效汇率")
	}
	if strings.TrimSpace(out.Date) == "" {
		out.Date = time.Now().Format("2006-01-02")
	}
	return out, nil
}

func fetchExchangeTrend(ctx context.Context, from, to string, days int) ([]exchangeTrendPoint, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	if days < 2 {
		days = 2
	}
	lookbackDays := days * 2
	if lookbackDays < 14 {
		lookbackDays = 14
	}
	if lookbackDays > 60 {
		lookbackDays = 60
	}
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -lookbackDays)
	url := fmt.Sprintf(
		"https://api.frankfurter.dev/v2/rates?base=%s&quotes=%s&from=%s&to=%s",
		strings.ToUpper(from),
		strings.ToUpper(to),
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("汇率趋势服务请求失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("汇率趋势服务返回异常状态: %d", resp.StatusCode)
	}
	var rows []frankfurterQuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&rows); err != nil {
		return nil, fmt.Errorf("解析汇率趋势结果失败: %w", err)
	}
	points := make([]exchangeTrendPoint, 0, len(rows))
	for _, row := range rows {
		if row.Rate <= 0 || strings.TrimSpace(row.Date) == "" {
			continue
		}
		points = append(points, exchangeTrendPoint{
			Date: row.Date,
			Rate: row.Rate,
		})
	}
	if len(points) == 0 {
		return nil, errors.New("暂未获取到可用的趋势数据")
	}
	if len(points) > days {
		points = points[len(points)-days:]
	}
	return points, nil
}
