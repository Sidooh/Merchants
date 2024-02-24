package presenter

import "gorm.io/datatypes"

type Payment struct {
	Id            uint           `json:"id"`
	Description   string         `json:"description"`
	Destination   datatypes.JSON `json:"destination"`
	Amount        float32        `json:"amount"`
	Charge        float32        `json:"charge"`
	TransactionId uint           `json:"transaction_id"`
	Status        string         `json:"status"`
}
