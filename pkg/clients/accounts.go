package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"net/http"
)

var accountClient *ApiClient

func InitAccountClient() {
	accountsApiUrl := viper.GetString("SIDOOH_ACCOUNTS_API_URL")
	accountClient = New(accountsApiUrl)
}

func GetAccountClient() *ApiClient {
	return accountClient
}

type AccountApiResponse struct {
	ApiResponse
	Data *Account `json:"data"`
}

func (api *ApiClient) CreateAccount(phone string) (*Account, error) {
	var apiResponse = new(AccountApiResponse)

	jsonData, err := json.Marshal(map[string]string{"phone": phone})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/accounts", dataBytes).Send(apiResponse)
	if err != nil {
		return nil, err
	}
	if apiResponse.Result == 0 {
		return nil, errors.New(apiResponse.Message)
	}

	return apiResponse.Data, nil
}

func (api *ApiClient) GetAccount(phone string) (*Account, error) {
	var apiResponse = new(AccountApiResponse)

	err := api.NewRequest(http.MethodGet, "/accounts/phone/"+phone, nil).Send(apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

func (api *ApiClient) GetOrCreateAccount(phone string) (*Account, error) {
	account, err := api.CreateAccount(phone)
	if err != nil {
		account, err = api.GetAccount(phone)
		if err != nil {
			return nil, err
		}
	}

	return account, nil
}

func (api *ApiClient) GetAccountById(id string) (*Account, error) {
	var apiResponse = new(AccountApiResponse)

	err := api.NewRequest(http.MethodGet, "/accounts/"+id, nil).Send(apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}
