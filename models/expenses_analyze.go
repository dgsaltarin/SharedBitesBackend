package models

type ExpensesAnalyze struct {
	Items []Expense `json:"items"`
	Date  string    `json:"date"`
	Total float64   `json:"total"`
}

type Expense struct {
	Item  string  `json:"item"`
	Price float64 `json:"price"`
}
