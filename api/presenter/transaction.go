package presenter

type Transaction struct {
	Id          uint    `json:"id"`
	Description string  `json:"description"`
	Destination string  `json:"destination"`
	Amount      float32 `json:"amount"`
	MerchantId  uint    `json:"merchant"`
	Product     string  `json:"product"`
	CreatedAt   string  `json:"created_at"`
}
