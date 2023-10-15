package entities

type EarningAccountTransaction struct {
	ModelID

	Type   string  `json:"type" gorm:"not null;size:16;"` // CREDIT / DEBIT
	Amount float32 `json:"amount" gorm:"not null;type:decimal(12,2);"`

	EarningAccountId uint `json:"earning_account_id" gorm:"not null;"`

	EarningAccount EarningAccount

	ModelTimeStamps
}
