package presenter

type Earning struct {
	Id     uint    `json:"id"`
	Type   string  `json:"type"`
	Amount float32 `json:"amount"`

	TransactionId uint `json:"transaction"`
	MerchantId    uint `json:"merchant"`
}
