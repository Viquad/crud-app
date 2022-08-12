package domain

import "time"

type Account struct {
	Id         int64     `json:"id"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Balance    int64     `json:"balance"`
	Currency   string    `json:"currency"`
	LastUpdate time.Time `json:"lastUpdate"`
}

type AccountUpdateInput struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Balance   *int64  `json:"balance"`
	Currency  *string `json:"currency"`
}
