package entities

type Transaction struct {
	ModelID

	Amount      float32 `json:"amount" gorm:"not null;type:decimal(10,2);"`
	Status      string  `json:"status" gorm:"size:16; default:PENDING"`
	Description string  `json:"description" gorm:"size:64"`

	Destination *string `json:"destination" gorm:"size:64"`
	MerchantId  uint    `json:"merchant_id" gorm:"not null"`
	Product     string  `json:"product" gorm:"not null;size:16"`

	Merchant Merchant `json:"-"`

	ModelTimeStamps
}
