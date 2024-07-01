package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/datatypes"
	"net/http"
	"time"
)

var savingsClient *ApiClient

func InitSavingsClient() {
	apiUrl := viper.GetString("SIDOOH_SAVINGS_API_URL")
	savingsClient = New(apiUrl)
	savingsClient.client = &http.Client{Timeout: 60 * time.Second}
}

func GetSavingsClient() *ApiClient {
	return savingsClient
}

type InvestmentTransaction struct {
}

type InvestmentsApiResponse struct {
	ApiResponse
	Data map[string]map[string][]InvestmentTransaction `json:"data"`
}

type Investment struct {
	AccountId        uint    `json:"account_id"`
	CashbackAmount   float32 `json:"cashback_amount"`
	CommissionAmount float32 `json:"commission_amount"`
}

type WithdrawalApiResponse struct {
	ApiResponse
	Data *Withdrawal `json:"data"`
}

type Withdrawal struct {
	Type              string         `json:"type"`
	Description       string         `json:"description"`
	Amount            int            `json:"amount"`
	PersonalAccountId uint           `json:"personal_account_id,string"`
	Extra             datatypes.JSON `json:"extra"`
	Id                uint           `json:"id,string"`
	Status            string         `json:"status"`
}

type PersonalAccountApiResponse struct {
	ApiResponse
	Data []PersonalAccount `json:"data"`
}

type PersonalAccount struct {
	Id          string      `json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	Type        string      `json:"type"`
	Description interface{} `json:"description"`
	Balance     float64     `json:"balance"`
	Status      string      `json:"status"`
	AccountId   string      `json:"account_id"`
}

func (api *ApiClient) SaveEarnings(investments []Investment) (map[string]map[string][]InvestmentTransaction, error) {
	res := new(InvestmentsApiResponse)

	jsonData, err := json.Marshal(investments)
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/accounts/merchant-earnings", dataBytes).Send(&res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (api *ApiClient) WithdrawSavings(personalAccId, destination, account, reference string, amount int) (*Withdrawal, error) {
	res := new(WithdrawalApiResponse)

	jsonData, err := json.Marshal(map[string]interface{}{
		"amount":              amount,
		"destination":         destination,
		"destination_account": account,
		"ipn":                 viper.GetString("APP_URL") + "/api/v1/savings/ipn",
		"reference":           reference,
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/personal-accounts/"+personalAccId+"/withdraw", dataBytes).Send(&res)
	fmt.Println(res, err)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (api *ApiClient) GetPersonalAccounts(accountId string) ([]PersonalAccount, error) {
	res := new(PersonalAccountApiResponse)

	err := api.NewRequest(http.MethodGet, "/accounts/"+accountId+"/personal-accounts", nil).Send(&res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}
