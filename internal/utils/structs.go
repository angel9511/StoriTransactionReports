package utils

type Transaction struct {
	ID     int     `json:"id"`
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

type SummaryRequestPayload struct {
	Recipient    string `json:"recipient"`
	Transactions string `json:"transactions"`
}
