package presenter

type MpesaAgentStoreAccount struct {
	Id    uint   `json:"id"`
	Agent string `json:"agent"`
	Store string `json:"store"`
	Name  string `json:"name"`
}
