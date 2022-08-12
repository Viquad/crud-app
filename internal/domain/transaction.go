package domain

import "time"

type Transaction struct {
	Id          int64
	AccountId   int64
	Amount      float64
	Description string
	Date        time.Time
}
