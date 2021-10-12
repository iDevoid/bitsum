package model

import "time"

// Coin is the data structure that contains the datetime and amount for historical record
type Coin struct {
	DateTime time.Time `json:"datetime"`
	Amount   float64   `json:"amount"`
}

// FilterDate contains start and end date time for getting the history
type FilterDate struct {
	StartDateTime time.Time `json:"start_date_time"`
	EndDateTime   time.Time `json:"end_date_time"`
}

// Response is for API response
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
