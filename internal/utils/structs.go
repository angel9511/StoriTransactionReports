package utils

type Transaction struct {
	ID     int     `json:"id"`
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}
