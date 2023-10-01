package presenter

type EarningAccount struct {
	Id     uint    `json:"id"`
	Type   string  `json:"type"`
	Amount float32 `json:"amount"`
}
