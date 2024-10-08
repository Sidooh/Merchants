package presenter

type Merchant struct {
	Id             uint   `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	IdNumber       string `json:"id_number"`
	Phone          string `json:"phone"`
	BusinessName   string `json:"business_name"`
	Code           string `json:"code"`
	AccountId      uint   `json:"account_id"`
	FloatAccountId uint   `json:"float_account_id"`
}
