package entities

type EarningAccount struct {
	ModelID

	Type   string  `json:"type" gorm:"not null;size:16;uniqueIndex:idx_earning_accounts"` // CASHBACK / COMMISSION
	Amount float32 `json:"amount" gorm:"not null;type:decimal(12,2);"`

	AccountId uint `json:"accountId" gorm:"not null;uniqueIndex:idx_earning_accounts"`

	ModelTimeStamps
}
