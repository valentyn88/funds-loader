package service

import (
	"reflect"
	"sync"
	"testing"

	"github.com/valentyn88/funds-loader/model"
	"github.com/valentyn88/funds-loader/storage"
)

func TestFundsLoader_Load(t *testing.T) {
	storage := storage.NewStorage(sync.RWMutex{}, map[string]storage.CustomerFunds{})

	t.Run("exceeded weekly limit new record", func(t *testing.T) {
		fLoader := NewFundsLoader(100.00, 20.00, 3,
			storage)

		got := fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "234",
			LoadAmount: "$105.57",
			Time:       "2000-01-01T00:00:00Z",
		})

		expected := model.Response{ID: "123", CustomerID: "234", Accepted: false}

		assertResponse(t, got, expected)
	})

	t.Run("exceeded daily limit new record", func(t *testing.T) {
		fLoader := NewFundsLoader(100.00, 20.00, 3,
			storage)

		got := fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "2345",
			LoadAmount: "$25.57",
			Time:       "2000-01-01T00:00:00Z",
		})

		expected := model.Response{ID: "123", CustomerID: "2345", Accepted: false}

		assertResponse(t, got, expected)
	})

	t.Run("save new record successfully", func(t *testing.T) {
		fLoader := NewFundsLoader(100.00, 20.00, 3,
			storage)

		got := fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "23456",
			LoadAmount: "$15.57",
			Time:       "2000-01-01T00:00:00Z",
		})
		expected := model.Response{ID: "123", CustomerID: "23456", Accepted: true}

		assertResponse(t, got, expected)
	})

	t.Run("exceeded weekly limit record exists", func(t *testing.T) {
		fLoader := NewFundsLoader(40.00, 20.00, 3,
			storage)

		fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "234567",
			LoadAmount: "$20.00",
			Time:       "2000-01-01T00:00:00Z",
		})

		got := fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "234567",
			LoadAmount: "$20.01",
			Time:       "2000-01-02T00:00:00Z",
		})

		expected := model.Response{ID: "123", CustomerID: "234567", Accepted: false}

		assertResponse(t, got, expected)
	})

	t.Run("exceeded weekly limit record exists", func(t *testing.T) {
		fLoader := NewFundsLoader(40.00, 20.00, 3,
			storage)

		fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "234567",
			LoadAmount: "$20.00",
			Time:       "2000-01-01T00:00:00Z",
		})

		got := fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "234567",
			LoadAmount: "$20.01",
			Time:       "2000-01-02T00:00:00Z",
		})

		expected := model.Response{ID: "123", CustomerID: "234567", Accepted: false}

		assertResponse(t, got, expected)
	})

	t.Run("exceeded daily limit record exists", func(t *testing.T) {
		fLoader := NewFundsLoader(100.00, 20.00, 3,
			storage)

		fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "2345678",
			LoadAmount: "$10.00",
			Time:       "2000-01-01T00:00:00Z",
		})

		got := fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "2345678",
			LoadAmount: "$10.01",
			Time:       "2000-01-01T05:11:00Z",
		})

		expected := model.Response{ID: "123", CustomerID: "2345678", Accepted: false}

		assertResponse(t, got, expected)
	})

	t.Run("exceeded daily attempts limit record exists", func(t *testing.T) {
		fLoader := NewFundsLoader(100.00, 50.00, 3,
			storage)

		fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "23456789",
			LoadAmount: "$10.00",
			Time:       "2000-01-01T00:00:00Z",
		})

		fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "23456789",
			LoadAmount: "$10.00",
			Time:       "2000-01-01T05:11:00Z",
		})

		fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "23456789",
			LoadAmount: "$10.00",
			Time:       "2000-01-01T06:11:00Z",
		})

		got := fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "23456789",
			LoadAmount: "$10.00",
			Time:       "2000-01-01T23:11:00Z",
		})

		expected := model.Response{ID: "123", CustomerID: "23456789", Accepted: false}

		assertResponse(t, got, expected)
	})

	t.Run("exceeded daily attempts limit record exists", func(t *testing.T) {
		fLoader := NewFundsLoader(100.00, 20.00, 3,
			storage)

		fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "2345678910",
			LoadAmount: "$10.00",
			Time:       "2000-01-01T00:00:00Z",
		})

		fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "2345678910",
			LoadAmount: "$10.00",
			Time:       "2000-01-01T05:11:00Z",
		})

		got := fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "2345678910",
			LoadAmount: "$20.01",
			Time:       "2000-01-02T23:11:00Z",
		})

		expected := model.Response{ID: "123", CustomerID: "2345678910", Accepted: false}

		assertResponse(t, got, expected)
	})

	t.Run("exceeded weekly limit record exists", func(t *testing.T) {
		fLoader := NewFundsLoader(40.00, 20.00, 3,
			storage)

		fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "1023456789",
			LoadAmount: "$10.00",
			Time:       "2000-01-01T00:00:00Z",
		})

		fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "1023456789",
			LoadAmount: "$20.15",
			Time:       "2000-01-10T23:11:00Z",
		})

		got := fLoader.Load(model.LoadedFunds{
			ID:         "123",
			CustomerID: "1023456789",
			LoadAmount: "$20.01",
			Time:       "2000-01-10T23:16:00Z",
		})

		expected := model.Response{ID: "123", CustomerID: "1023456789", Accepted: false}

		assertResponse(t, got, expected)
	})

	t.Run("exceeded daily limit record exists", func(t *testing.T) {
		testCases := []struct {
			loadedFunds model.LoadedFunds
			expected    model.Response
		}{
			{
				loadedFunds: model.LoadedFunds{
					ID:         "21336",
					CustomerID: "477",
					LoadAmount: "$3725.46",
					Time:       "2000-01-02T06:41:00Z",
				},
				expected: model.Response{
					ID:         "21336",
					CustomerID: "477",
					Accepted:   true,
				},
			},
			{
				loadedFunds: model.LoadedFunds{
					ID:         "21336",
					CustomerID: "477",
					LoadAmount: "$5678.22",
					Time:       "2000-01-03T07:13:48Z",
				},
				expected: model.Response{
					ID:         "21336",
					CustomerID: "477",
					Accepted:   false,
				},
			},
		}

		fLoader := NewFundsLoader(WeeklyAmountLimit, DailyAmountLimit, AttemptsLimit,
			storage)

		for _, tc := range testCases {
			got := fLoader.Load(tc.loadedFunds)
			assertResponse(t, got, tc.expected)
		}

	})
}

func assertResponse(t *testing.T, got, expected model.Response) {
	t.Helper()

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected %v got %v", expected, got)
	}
}
