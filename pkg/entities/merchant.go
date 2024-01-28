package entities

type Merchant struct {
	ModelID

	FirstName string `json:"firstname" gorm:"size:32"`
	LastName  string `json:"lastname" gorm:"size:32"`
	IdNumber  string `json:"id_number" gorm:"unique; size 16"`
	Phone     string `json:"phone" gorm:"not null; uniqueIndex; size:16"`

	BusinessName *string `json:"business_name" gorm:"size:128"`
	Code         *uint   `json:"code" gorm:"unique; size:24"`

	AccountId      uint    `json:"account_id" gorm:"not null; uniqueIndex"`
	FloatAccountId *uint   `json:"-" gorm:"uniqueIndex"`
	LocationId     *uint   `json:"-"`
	Landmark       *string `json:"-" gorm:"size:128"`

	ModelTimeStamps
}
