package entities

type EarningAccount struct {
	ModelID

	Type   string  `json:"type" gorm:"size:16;"` // CASHBACK / COMMISSION
	Amount float32 `json:"amount" gorm:"not null;type:decimal(12,2);"`

	MerchantId uint `json:"merchantId" gorm:"not null;"`

	Merchant Merchant `json:"-"`

	ModelTimeStamps
}
