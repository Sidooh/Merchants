package entities

type MpesaAgentStoreAccount struct {
	ModelID

	Agent string `json:"agent" gorm:"not null;size:16;uniqueIndex:idx_account;index:idx_store"`
	Store string `json:"status" gorm:"not null;size:16;uniqueIndex:idx_account;index:idx_store"`
	Name  string `json:"name" gorm:"size:64;index:idx_store"`

	MerchantId uint `json:"merchant_id" gorm:"not null;uniqueIndex:idx_account"`

	Merchant Merchant `json:"-"`

	ModelTimeStamps
}
