package clients

import (
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var savingsClient *ApiClient

func InitSavingsClient() {
	apiUrl := viper.GetString("SIDOOH_SAVINGS_API_URL")
	notifyClient = New(apiUrl)
	notifyClient.client = &http.Client{Timeout: 60 * time.Second}
}

func GetSavingsClient() *ApiClient {
	return savingsClient
}

type SavingsAccountApiResponse struct {
	ApiResponse
	Data []string `json:"data"`
}

func (api *ApiClient) SaveEarnings() ([]string, error) {
	res := new(SavingsAccountApiResponse)

	err := api.NewRequest(http.MethodGet, "/accounts/earnings", nil).Send(&res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}
