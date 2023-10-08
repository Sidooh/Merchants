package clients

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
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

type SavingsAccountApiResponse struct {
	ApiResponse
	Data []string `json:"data"`
}

type Investment struct {
	AccountId        uint    `json:"account_id"`
	CashbackAmount   float32 `json:"cashback_amount"`
	CommissionAmount float32 `json:"commission_amount"`
}

func (api *ApiClient) SaveEarnings(investments []Investment) ([]string, error) {
	res := new(SavingsAccountApiResponse)

	jsonData, err := json.Marshal(investments)
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/accounts/merchant-earnings", dataBytes).Send(&res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}
