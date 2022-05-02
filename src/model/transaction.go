package model

import "time"

type Transaction struct {
	Amount   int
	Merchant string
	Time     time.Time
}
