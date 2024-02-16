package presenter

import "time"

type Transaction struct {
	Id          uint      `json:"id"`
	Description string    `json:"description"`
	Destination string    `json:"destination"`
	Status      string    `json:"status"`
	Amount      float32   `json:"amount"`
	MerchantId  uint      `json:"merchant"`
	Product     string    `json:"product"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Payment     *Payment  `json:"payment,omitempty"`
}
