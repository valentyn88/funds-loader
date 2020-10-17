package model

import (
	"strconv"
	"strings"
)

// LoadedFunds represent client loaded funds.
type LoadedFunds struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	LoadAmount string `json:"load_amount"`
	Time       string `json:"time"`
}

// LoadAmount2Float64 converts string to float64.
func (l *LoadedFunds) LoadAmount2Float64() (float64, error) {
	lAmount := strings.TrimPrefix(l.LoadAmount, "$")
	return strconv.ParseFloat(lAmount, 64)
}
