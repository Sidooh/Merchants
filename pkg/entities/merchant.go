package entities

type Merchant struct {
	ModelID

	FirstName string `json:"firstname" gorm:"size:32"`
	LastName  string `json:"lastname" gorm:"size:32"`
	IdNumber  string `json:"id_number" gorm:"unique; size 16"`

	BusinessName *string `json:"business_name" gorm:"size:128"`
	Code         *string `json:"code" gorm:"unique; size:8"`

	AccountId      uint    `json:"account_id" gorm:"uniqueIndex"`
	FloatAccountId *uint   `json:"-" gorm:"uniqueIndex"`
	LocationId     *uint   `json:"-"`
	Landmark       *string `json:"-" gorm:"size:128"`

	ModelTimeStamps
}
