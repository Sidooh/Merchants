package entities

type MpesaAgentStoreAccount struct {
	ModelID

	Agent string `json:"agent" gorm:"not null;size:16;uniqueIndex:idx_account"`
	Store string `json:"status" gorm:"not null;size:16;uniqueIndex:idx_account"`
	Name  string `json:"name" gorm:"size:64"`

	MerchantId uint `json:"merchant_id" gorm:"not null;uniqueIndex:idx_account"`

	Merchant Merchant `json:"-"`

	ModelTimeStamps
}
