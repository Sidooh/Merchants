package clients

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"merchants.sidooh/pkg/cache"
	"merchants.sidooh/utils"
	"net/http"
	"strconv"
	"time"
)

var paymentClient *ApiClient

var paymentsCache cache.ICache[string, interface{}]

func InitPaymentClient() {
	apiUrl := viper.GetString("SIDOOH_PAYMENTS_API_URL")
	paymentClient = New(apiUrl)
	paymentClient.client = &http.Client{Timeout: 120 * time.Second}

	paymentsCache = cache.New[string, interface{}]()
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

func (api *ApiClient) CreateFloatAccount(merchantId, accountId, code int) (*FloatAccount, error) {
	var apiResponse = new(FloatAccountApiResponse)

	jsonData, err := json.Marshal(map[string]string{
		"initiator":   "MERCHANT",
		"reference":   strconv.Itoa(merchantId),
		"description": strconv.Itoa(code),
		"account_id":  strconv.Itoa(accountId),
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

func (api *ApiClient) FetchFloatAccount(id string) (*FloatAccount, error) {
	var apiResponse = new(FloatAccountApiResponse)

	var endpoint = "/float-accounts/" + id
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

func (api *ApiClient) BuyMpesaFloat(accountId uint, amount int, agent, store, source, sourceAccount string) (*utils.Payment, error) {
	var apiResponse = new(utils.PaymentApiResponse)

	jsonData, err := json.Marshal(map[string]interface{}{
		"account_id":  accountId,
		"amount":      amount,
		"description": "Mpesa Float Purchase",
		//"reference": "test",
		"source":         source,
		"source_account": sourceAccount,
		"ipn":            viper.GetString("APP_URL") + "/api/v1/payments/ipn",
		"merchant_type":  "MPESA_STORE",
		"agent":          agent,
		"store":          store,
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/payments/mpesa-float", dataBytes).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) MpesaWithdraw(accountId, floatAccountId uint, amount int, phone string) (*utils.Payment, error) {
	var apiResponse = new(utils.PaymentApiResponse)

	jsonData, err := json.Marshal(map[string]interface{}{
		"account_id":          accountId,
		"amount":              amount,
		"description":         "Mpesa Withdrawal",
		"source":              "MPESA",
		"source_account":      phone,
		"ipn":                 viper.GetString("APP_URL") + "/api/v1/payments/ipn",
		"destination":         "FLOAT",
		"destination_account": floatAccountId,
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/payments/mpesa-withdraw", dataBytes).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) FloatPurchase(accountId, floatAccountId uint, amount int, phone string) (*utils.Payment, error) {
	var apiResponse = new(utils.PaymentApiResponse)

	jsonData, err := json.Marshal(map[string]interface{}{
		"account_id":          accountId,
		"amount":              amount,
		"description":         "Float Credit",
		"source":              "MPESA",
		"source_account":      phone,
		"ipn":                 viper.GetString("APP_URL") + "/api/v1/payments/ipn",
		"destination":         "FLOAT",
		"destination_account": floatAccountId,
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/payments/merchant-float", dataBytes).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) GetWithdrawalCharges() ([]utils.AmountCharge, error) {
	endpoint := "/charges/withdrawal"
	apiResponse := new(utils.ChargesApiResponse)

	// TODO: Configure caching
	//charges := cache.Cache.Get[[]utils.AmountCharge](endpoint)
	//if charges != nil && len(charges) > 0 {
	//	return *charges, nil
	//}

	if err := api.NewRequest(http.MethodGet, endpoint, nil).Send(&apiResponse); err != nil {
		return nil, err
	}
	//cache.Cache.Set(endpoint, apiResponse.Data, 28*24*time.Hour)

	return *apiResponse.Data, nil
}

func (api *ApiClient) GetMpesaCollectionCharges() ([]utils.AmountCharge, error) {
	endpoint := "/charges/mpesa-collection"
	apiResponse := new(utils.ChargesApiResponse)

	// TODO: Configure caching
	cacheValue := paymentsCache.Get(endpoint)
	if cacheValue != nil {
		charges := *(*cacheValue).(*[]utils.AmountCharge)
		if charges != nil && len(charges) > 0 {
			return charges, nil
		}
	}

	if err := api.NewRequest(http.MethodGet, endpoint, nil).Send(&apiResponse); err != nil {
		return nil, err
	}
	paymentsCache.Set(endpoint, apiResponse.Data, 28*24*time.Hour)

	return *apiResponse.Data, nil
}

func (api *ApiClient) Withdraw(accountId, floatAccountId uint, amount int, destination, destinationAccount string) (*utils.Payment, error) {
	var apiResponse = new(utils.PaymentApiResponse)

	jsonData, err := json.Marshal(map[string]interface{}{
		"account_id":  accountId,
		"amount":      amount,
		"description": "Merchant Withdrawal",
		//"reference": "test",
		"source":              "FLOAT",
		"source_account":      floatAccountId,
		"ipn":                 viper.GetString("APP_URL") + "/api/v1/payments/ipn",
		"destination":         destination,
		"destination_account": destinationAccount,
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/payments/withdraw", dataBytes).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) GetMpesaWithdrawalCommissions() ([]utils.AmountCharge, error) {
	return []utils.AmountCharge{
		{Min: 50, Max: 100, Charge: 5},
		{Min: 101, Max: 1500, Charge: 13},
		{Min: 1501, Max: 2500, Charge: 10},
		{Min: 2501, Max: 3500, Charge: 23},
		{Min: 3501, Max: 5000, Charge: 30},
		{Min: 5001, Max: 7500, Charge: 35},
		{Min: 7501, Max: 10000, Charge: 40},
		{Min: 10001, Max: 15000, Charge: 73},
		{Min: 15001, Max: 20000, Charge: 80},
		{Min: 20001, Max: 35000, Charge: 86},
		{Min: 35001, Max: 50000, Charge: 121},
		{Min: 50001, Max: 150000, Charge: 133},
	}, nil
}

func (api *ApiClient) GetMpesaWithdrawalInviterCommissions() ([]utils.AmountCharge, error) {
	return []utils.AmountCharge{
		{Min: 50, Max: 100, Charge: 1},
		{Min: 101, Max: 1500, Charge: 3},
		{Min: 1501, Max: 2500, Charge: 2},
		{Min: 2501, Max: 3500, Charge: 5},
		{Min: 3501, Max: 5000, Charge: 6},
		{Min: 5001, Max: 7500, Charge: 7},
		{Min: 7501, Max: 10000, Charge: 8},
		{Min: 10001, Max: 15000, Charge: 15},
		{Min: 15001, Max: 20000, Charge: 16},
		{Min: 20001, Max: 35000, Charge: 18},
		{Min: 35001, Max: 50000, Charge: 25},
		{Min: 50001, Max: 150000, Charge: 27},
	}, nil
}
