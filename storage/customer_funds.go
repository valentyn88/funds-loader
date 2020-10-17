package storage

import (
	"sync"
	"time"
)

// CustomerFunds represents customer funds.
type CustomerFunds struct {
	WeeklyAmount float64
	DayAmount    float64
	DayAttempts  int
	LastLoadDate time.Time
}

// Storage represents storage of user funds.
type Storage struct {
	mtx  sync.RWMutex
	data map[string]CustomerFunds
}

// NewStorage init new Storage.
func NewStorage(mtx sync.RWMutex, data map[string]CustomerFunds) *Storage {
	return &Storage{mtx: mtx, data: data}
}

// GetFundsByCustomerID gets customer funds by customerID.
func (s *Storage) GetFundsByCustomerID(cutomerID string) (CustomerFunds, bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	custFunds, ok := s.data[cutomerID]
	return custFunds, ok
}

// Save saves customer funds.
func (s *Storage) Save(customerID string, amount float64, lastLoadDate time.Time) CustomerFunds {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	customerFunds := s.data[customerID]
	customerFunds.WeeklyAmount += amount
	customerFunds.DayAmount += amount
	customerFunds.LastLoadDate = lastLoadDate
	customerFunds.DayAttempts += 1
	s.data[customerID] = customerFunds

	return customerFunds
}

// UpdToDefVals updates customer account data.
func (s *Storage) UpdToDefVals(customerID string, wAmount *float64, dAmount float64, lastLoadDate time.Time) CustomerFunds {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	customerFunds := s.data[customerID]
	if wAmount != nil {
		customerFunds.WeeklyAmount = *wAmount
	}
	customerFunds.DayAmount = dAmount
	customerFunds.DayAttempts = 0
	customerFunds.LastLoadDate = lastLoadDate
	s.data[customerID] = customerFunds

	return customerFunds
}

// UpdAttempts updates customer attempts.
func (s *Storage) UpdAttempts(customerID string) CustomerFunds {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	customerFunds := s.data[customerID]
	customerFunds.DayAttempts += 1
	s.data[customerID] = customerFunds

	return customerFunds
}
