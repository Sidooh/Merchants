package entities

type Earning struct {
	ModelID

	Amount float32 `json:"amount" gorm:"not null;type:decimal(7,2);"`
	Type   string  `json:"type" gorm:"size:16;"` //SELF / INVITE / SYSTEM

	TransactionId uint `json:"transaction_id" gorm:"not null;uniqueIndex"`

	Transaction Transaction `json:"-"`

	MerchantId uint `json:"merchantId" gorm:"uniqueIndex"`

	Merchant Merchant `json:"-"`

	ModelTimeStamps
}
