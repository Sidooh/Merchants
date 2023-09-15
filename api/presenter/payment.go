package presenter

type Payment struct {
	Id            uint    `json:"id"`
	Description   string  `json:"description"`
	Destination   string  `json:"destination"`
	Amount        float32 `json:"amount"`
	TransactionId uint    `json:"transaction_id"`
	Status        string  `json:"status"`
}
