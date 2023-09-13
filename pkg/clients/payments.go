package clients

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

var paymentClient *ApiClient

func InitPaymentClient() {
	apiUrl := viper.GetString("SIDOOH_PAYMENTS_API_URL")
	paymentClient = New(apiUrl)
}

func GetPaymentClient() *ApiClient {
	return paymentClient
}

type PaymentApiResponse struct {
	ApiResponse
	Data *Payment `json:"data"`
}

type VoucherTypesApiResponse struct {
	ApiResponse
	Data *[]VoucherType `json:"data"`
}

type VoucherTypeApiResponse struct {
	ApiResponse
	Data *VoucherType `json:"data"`
}

type VouchersApiResponse struct {
	ApiResponse
	Data []*Voucher `json:"data"`
}

type VoucherApiResponse struct {
	ApiResponse
	Data *Voucher `json:"data"`
}

type VoucherTransactionsApiResponse struct {
	ApiResponse
	Data *[]VoucherTransaction `json:"data"`
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

func (api *ApiClient) CreateFloatAccount(enterpriseId, accountId int) (*FloatAccount, error) {
	var apiResponse = new(FloatAccountApiResponse)

	jsonData, err := json.Marshal(map[string]string{
		"initiator":  "ENTERPRISE",
		"reference":  strconv.Itoa(enterpriseId),
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
		"reference":      "ENTERPRISE",
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// VOUCHER TYPES
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (api *ApiClient) FetchVoucherTypes(accountId int) (*[]VoucherType, error) {
	var apiResponse = new(VoucherTypesApiResponse)

	err := api.NewRequest(http.MethodGet, "/voucher-types?account_id="+strconv.Itoa(accountId), nil).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) FetchVoucherType(accountId, voucherTypeId int) (*VoucherType, error) {
	var apiResponse = new(VoucherTypeApiResponse)

	var endpoint = "/voucher-types/" + strconv.Itoa(voucherTypeId) + "?account_id=" + strconv.Itoa(accountId) + "&with=vouchers"
	err := api.NewRequest(http.MethodGet, endpoint, nil).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) CreateVoucherType(accountId int, name string) (*VoucherType, error) {
	var apiResponse = new(VoucherTypeApiResponse)

	jsonData, err := json.Marshal(map[string]string{
		"initiator":  "ENTERPRISE",
		"name":       name,
		"account_id": strconv.Itoa(accountId),
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/voucher-types", dataBytes).Send(apiResponse)

	return apiResponse.Data, err
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// VOUCHERS
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (api *ApiClient) FetchVouchers(accountId int) ([]*Voucher, error) {
	var apiResponse = new(VouchersApiResponse)

	err := api.NewRequest(http.MethodGet, "/vouchers?account_id="+strconv.Itoa(accountId), nil).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) DisburseVoucher(accountId, floatAccountId, voucherId, amount int) (*Payment, error) {
	var apiResponse = new(PaymentApiResponse)

	jsonData, err := json.Marshal(map[string]interface{}{
		"account_id":     accountId,
		"amount":         amount,
		"description":    "Voucher Disbursement",
		"reference":      "ENTERPRISE",
		"source":         "FLOAT",
		"source_account": floatAccountId,
		"voucher":        voucherId,
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/vouchers/credit", dataBytes).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) CreateVoucher(enterpriseAccountId, voucherTypeId int) (*Voucher, error) {
	var apiResponse = new(VoucherApiResponse)

	jsonData, err := json.Marshal(map[string]interface{}{
		"account_id":      enterpriseAccountId,
		"voucher_type_id": voucherTypeId,
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/vouchers", dataBytes).Send(apiResponse)

	return apiResponse.Data, err
}

func (api *ApiClient) FetchVoucherTransactions(accountId int, limit int) (*[]VoucherTransaction, error) {
	var apiResponse = new(VoucherTransactionsApiResponse)

	var endpoint = "/voucher-transactions?with=voucher&account_id=" + strconv.Itoa(accountId)
	if limit > 0 {
		endpoint += "&limit=" + strconv.Itoa(limit)
	}

	err := api.NewRequest(http.MethodGet, endpoint, nil).Send(apiResponse)

	return apiResponse.Data, err
}
