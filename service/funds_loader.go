package service

import (
	"github.com/valentyn88/funds-loader/model"
	"github.com/valentyn88/funds-loader/storage"
	"github.com/valentyn88/funds-loader/tools"
)

const (
	// WeeklyAmountLimit allowed amount per week.
	WeeklyAmountLimit = 20000.00
	// DailyAmountLimit allowed amount per day.
	DailyAmountLimit = 5000.00
	// AttemptsLimit allowed attempts per day.
	AttemptsLimit = 3
)

// FundsLoader loads funds a customer account.
type FundsLoader struct {
	weeklyAmountLimit  float64
	dailyAmountLimit   float64
	dailyAttemptsLimit int
	storage            *storage.Storage
}

// NewFundsLoader creates new FundLoader.
func NewFundsLoader(wAmntLim, dAmntLim float64, dAttemptsLim int, strg *storage.Storage) *FundsLoader {
	return &FundsLoader{
		weeklyAmountLimit:  wAmntLim,
		dailyAmountLimit:   dAmntLim,
		dailyAttemptsLimit: dAttemptsLim,
		storage:            strg,
	}
}

// Load loads money on customer account.
func (f *FundsLoader) Load(funds model.LoadedFunds) model.Response {
	customerFunds, ok := f.storage.GetFundsByCustomerID(funds.CustomerID)
	if !ok {
		return f.loadNewRecord(funds)
	}

	return f.updateExistsRecord(customerFunds, funds)
}

func (f FundsLoader) updateExistsRecord(customerFunds storage.CustomerFunds, funds model.LoadedFunds) model.Response {
	resp := model.Response{ID: funds.ID, CustomerID: funds.CustomerID, Accepted: false}
	t, err := tools.ParseDate(funds.Time)
	if err != nil {
		return resp
	}

	loadedAmount, err := funds.LoadAmount2Float64()
	if err != nil {
		return resp
	}

	if tools.IsSameWeek(t, customerFunds.LastLoadDate) {
		// Check weekly limit
		if tools.IsAmountExceeded(f.weeklyAmountLimit, customerFunds.WeeklyAmount+loadedAmount) {
			f.storage.UpdAttempts(funds.ID)
			return resp
		}
		// Check daily limits
		if tools.IsSameDay(t, customerFunds.LastLoadDate) {
			// Check daily amount
			if tools.IsAmountExceeded(f.dailyAmountLimit, customerFunds.DayAmount+loadedAmount) {
				f.storage.UpdAttempts(funds.ID)
				return resp
			}

			// Check daily attempts
			if customerFunds.DayAttempts+1 > f.dailyAttemptsLimit {
				f.storage.UpdAttempts(funds.ID)
				return resp
			}
		} else {
			customerFunds = f.storage.UpdToDefVals(funds.CustomerID, nil, 0.00, t)

			// Check daily amount
			if tools.IsAmountExceeded(f.dailyAmountLimit, customerFunds.DayAmount+loadedAmount) {
				f.storage.UpdAttempts(funds.ID)
				return resp
			}
		}

		f.storage.Save(funds.CustomerID, loadedAmount, t)
	} else {
		wAmount := 0.00
		customerFunds = f.storage.UpdToDefVals(funds.CustomerID, &wAmount, 0.00, t)
		// Check weekly limit
		if tools.IsAmountExceeded(f.weeklyAmountLimit, customerFunds.WeeklyAmount+loadedAmount) {
			f.storage.UpdAttempts(funds.ID)
			return resp
		}

		// Check daily amount
		if tools.IsAmountExceeded(f.dailyAmountLimit, customerFunds.DayAmount+loadedAmount) {
			f.storage.UpdAttempts(funds.ID)
			return resp
		}

		f.storage.Save(funds.CustomerID, loadedAmount, t)
	}

	resp.Accepted = true
	return resp
}

func (f FundsLoader) loadNewRecord(funds model.LoadedFunds) model.Response {
	resp := model.Response{ID: funds.ID, CustomerID: funds.CustomerID, Accepted: false}
	t, err := tools.ParseDate(funds.Time)
	if err != nil {
		return resp
	}

	loadedAmount, err := funds.LoadAmount2Float64()
	if err != nil {
		return resp
	}

	if tools.IsAmountExceeded(f.weeklyAmountLimit, loadedAmount) {
		return resp
	}

	if tools.IsAmountExceeded(f.dailyAmountLimit, loadedAmount) {
		return resp
	}

	f.storage.Save(funds.CustomerID, loadedAmount, t)

	resp.Accepted = true
	return resp
}
