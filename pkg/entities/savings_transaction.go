package entities

import "gorm.io/datatypes"

type SavingsTransaction struct {
	ModelID

	Type              string  `json:"type" gorm:"size:32;"`
	Amount            float32 `json:"amount" gorm:"not null;type:decimal(10,2);"`
	Status            string  `json:"status" gorm:"size:16; default:PENDING"`
	Description       string  `json:"description" gorm:"size:128"`
	PersonalAccountId uint    `json:"personal_account_id" gorm:"not_null"`

	Extra         datatypes.JSON `json:"extra"`
	TransactionId uint           `json:"transaction_id" gorm:"not null"`

	Transaction Transaction `json:"-"`

	SavingsId uint `json:"savings_id" gorm:"unique,not null"`

	ModelTimeStamps
}
