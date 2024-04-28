package entity

import "time"

type Bill struct {
	ID          string
	Items       []Item
	Total       float64
	People      int
	SplitEqualy bool
	UserID      string
	Date        time.Time
}

type Item struct {
	Name  string
	Price float64
}
