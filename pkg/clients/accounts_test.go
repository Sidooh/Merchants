package clients

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestInitAccountClient(t *testing.T) {
	viper.Set("SIDOOH_ACCOUNTS_API_URL", "test.test")
	InitAccountClient()

	assert.NotNil(t, accountClient, "account client is nil")
	assert.NotNil(t, accountClient.client, "http client is nil")
	assert.Nil(t, accountClient.request, "request is not nil")
	assert.NotNil(t, accountClient.cache, "cache is nil")

	assert.Equal(t, "test.test", accountClient.baseUrl)

	accountClient = nil
}

func TestGetAccountClient(t *testing.T) {
	api := GetAccountClient()
	assert.Nil(t, api)

	InitAccountClient()
	api = GetAccountClient()
	assert.NotNil(t, api)
}

func accountFoundRequest() RoundTripFunc {
	return func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader(`{"result":1,"data":{"id":1,"phone":"25412345678","active":true,"inviter_id":1,"user_id":50}}`)),
			// Must be set to non-nil value, or it panics
			//Header: make(http.Header),
		}
	}
}

func accountNotFoundRequest() RoundTripFunc {
	return func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: 500,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader(`{"result":0,"message":"Something went wrong, please try again."}`)),
			// Must be set to non-nil value, or it panics
			//Header: make(http.Header),
		}
	}
}

func TestApiClient_CreateAccount(t *testing.T) {
	InitAccountClient()
	api := GetAccountClient()

	type args struct {
		phone string
	}
	tests := []struct {
		name    string
		apiMock RoundTripFunc
		args    args
		want    *Account
		wantErr assert.ErrorAssertionFunc
	}{
		{"account is not found", accountNotFoundRequest(), args{phone: "25412345678"}, nil, assert.Error},
		{"account is found", accountFoundRequest(), args{phone: "25412345678"}, &Account{
			Id:     1,
			Phone:  "25412345678",
			Active: true,
		}, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api.client = &http.Client{Transport: tt.apiMock}
			got, err := api.CreateAccount(tt.args.phone)
			if !tt.wantErr(t, err, fmt.Sprintf("CreateAccount(%v)", tt.args.phone)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CreateAccount(%v)", tt.args.phone)
		})
	}
}

func TestApiClient_GetAccount(t *testing.T) {
	InitAccountClient()
	api := GetAccountClient()

	type args struct {
		phone string
	}
	tests := []struct {
		name    string
		apiMock RoundTripFunc
		args    args
		want    *Account
		wantErr assert.ErrorAssertionFunc
	}{
		{"account is not found", accountNotFoundRequest(), args{phone: "25412345678"}, nil, assert.Error},
		{"account is found", accountFoundRequest(), args{phone: "25412345678"}, &Account{
			Id:     1,
			Phone:  "25412345678",
			Active: true,
		}, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Removed t.Parallel cause of race condition, confirm where race cond is found
			api.client = &http.Client{Transport: tt.apiMock}
			got, err := api.GetAccount(tt.args.phone)
			if !tt.wantErr(t, err, fmt.Sprintf("GetAccount(%v)", tt.args.phone)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAccount(%v)", tt.args.phone)
		})
	}
}
