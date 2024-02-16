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

func TestInitNotifyClient(t *testing.T) {
	viper.Set("SIDOOH_NOTIFY_API_URL", "test.test")
	InitNotifyClient()

	assert.NotNil(t, notifyClient, "account client is nil")
	assert.NotNil(t, notifyClient.client, "http client is nil")
	assert.Nil(t, notifyClient.request, "request is not nil")
	assert.NotNil(t, notifyClient.cache, "cache is nil")

	assert.Equal(t, "test.test", notifyClient.baseUrl)

	notifyClient = nil
}

func TestGetNotifyClient(t *testing.T) {
	api := GetNotifyClient()
	assert.Nil(t, api)

	InitNotifyClient()
	api = GetNotifyClient()
	assert.NotNil(t, api)
}

func notificationSentRequest() RoundTripFunc {
	return func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader(`{"result":1,"data":{"id":1}}`)),
			// Must be set to non-nil value, or it panics
			//Header: make(http.Header),
		}
	}
}

func notificationFailedRequest() RoundTripFunc {
	return func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: 400,
			// Send response to be tested
			Body: ioutil.NopCloser(strings.NewReader(`{"result":0,"message":"Something went wrong, please try again."}`)),
			// Must be set to non-nil value, or it panics
			//Header: make(http.Header),
		}
	}
}

func TestApiClient_SendNotification(t *testing.T) {
	InitNotifyClient()
	api := GetNotifyClient()

	type args struct {
		channel     string
		event       string
		destination string
		message     string
	}
	tests := []struct {
		name    string
		apiMock RoundTripFunc
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{"notification sent", notificationSentRequest(), args{
			channel:     "SMS",
			event:       "TEST",
			destination: "25412345678",
			message:     "Test",
		}, assert.NoError},
		{"notification not sent", notificationFailedRequest(), args{
			channel:     "SMS",
			event:       "TEST",
			destination: "25412345678",
			message:     "Test",
		}, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api.client = &http.Client{Transport: tt.apiMock}
			tt.wantErr(t, api.SendNotification(tt.args.channel, tt.args.event, tt.args.destination, tt.args.message), fmt.Sprintf("SendNotification(%v, %v, %v, %v)", tt.args.channel, tt.args.event, tt.args.destination, tt.args.message))
		})
	}
}
