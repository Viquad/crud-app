package domain

import "time"

// TODO: still not implemented
type Transaction struct {
	Id          int64
	AccountId   int64
	Amount      float64
	Description string
	Date        time.Time
}
