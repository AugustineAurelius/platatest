package repository

import "time"

type Price struct {
	Value float64
	Date  time.Time
}

type Prices []Price

type Currency struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Currencies []Currency
