package model

// Response represends service response.
type Response struct {
	ID string `json:"id"`
	CustomerID string `json:"customer_id"`
	Accepted bool `json:"accepted"`
}
