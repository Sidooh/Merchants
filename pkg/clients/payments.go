package clients

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"merchants.sidooh/utils"
	"net/http"
	"strconv"
	"time"
)

var paymentClient *ApiClient

func InitPaymentClient() {
	apiUrl := viper.GetString("SIDOOH_PAYMENTS_API_URL")
	paymentClient = New(apiUrl)
	paymentClient.client = &http.Client{Timeout: 60 * time.Second}
}

func GetPaymentClient() *ApiClient {
	return paymentClient
}

type FloatAccountApiResponse struct {
	ApiResponse
	Data *FloatAccount `json:"data"`
}

type FloatAccountTransactionsApiResponse struct {
	ApiResponse
	Data *[]FloatAccountTransaction `json:"data"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// FLOAT ACCOUNTS
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (api *ApiClient) CreateFloatAccount(merchantId, accountId int) (*FloatAccount, error) {
	var apiResponse = new(FloatAccountApiResponse)

	jsonData, err := json.Marshal(map[string]string{
		"initiator":  "MERCHANT",
		"reference":  strconv.Itoa(merchantId),
		"account_id": strconv.Itoa(accountId),
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/float-accounts", dataBytes).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) CreditFloatAccount(accountId, floatAccountId, amount, phone int) (*interface{}, error) {
	var apiResponse = new(ApiResponse)

	jsonData, err := json.Marshal(map[string]interface{}{
		"account_id":     accountId,
		"amount":         amount,
		"description":    "Float Credit",
		"reference":      "MERCHANT",
		"source":         "MPESA",
		"source_account": phone,
		"float_account":  floatAccountId,
	})
	dataBytes := bytes.NewBuffer(jsonData)

	var endpoint = "/float-accounts/credit"
	err = api.NewRequest(http.MethodPost, endpoint, dataBytes).Send(apiResponse)

	return &apiResponse.Data, err
}

func (api *ApiClient) FetchFloatAccount(accountId int) (*FloatAccount, error) {
	var apiResponse = new(FloatAccountApiResponse)

	var endpoint = "/float-accounts/" + strconv.Itoa(accountId)
	err := api.NewRequest(http.MethodGet, endpoint, nil).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) FetchFloatAccountTransactions(accountId int, limit int) (*[]FloatAccountTransaction, error) {
	var apiResponse = new(FloatAccountTransactionsApiResponse)

	var endpoint = "/float-account-transactions?float_account_id=" + strconv.Itoa(accountId)
	if limit > 0 {
		endpoint += "&limit=" + strconv.Itoa(limit)
	}

	err := api.NewRequest(http.MethodGet, endpoint, nil).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) BuyMpesaFloat(accountId, floatAccountId uint, amount int, agent, store string) (*utils.Payment, error) {
	var apiResponse = new(utils.PaymentApiResponse)

	jsonData, err := json.Marshal(map[string]interface{}{
		"account_id":  accountId,
		"amount":      amount,
		"description": "Mpesa Float Purchase",
		//"reference": "test",
		"source":         "FLOAT",
		"source_account": floatAccountId,
		"ipn":            viper.GetString("APP_URL") + "/api/v1/payments/ipn",
		"merchant_type":  "MPESA_STORE",
		"agent":          agent,
		"store":          store,
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/payments/mpesa-float", dataBytes).Send(apiResponse)

	return apiResponse.Data, err
}
