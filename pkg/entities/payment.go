package entities

import "gorm.io/datatypes"

type Payment struct {
	ModelID

	Amount      float32 `json:"amount" gorm:"not null;type:decimal(10,2);"`
	Status      string  `json:"status" gorm:"size:16; default:PENDING"`
	Description string  `json:"description" gorm:"size:64"`

	Destination   datatypes.JSON `json:"destination"`
	TransactionId uint           `json:"transaction_id" gorm:"not null"`

	Transaction Transaction `json:"-"`

	PaymentId uint `json:"payment_id" gorm:"not null"`

	ModelTimeStamps
}
