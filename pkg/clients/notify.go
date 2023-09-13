package clients

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
)

var notifyClient *ApiClient

func InitNotifyClient() {
	apiUrl := viper.GetString("SIDOOH_NOTIFY_API_URL")
	notifyClient = New(apiUrl)
}

func GetNotifyClient() *ApiClient {
	return notifyClient
}

func (api *ApiClient) SendSMS(event, phone, message string) error {
	return api.SendNotification("sms", event, phone, message)
}

func (api *ApiClient) SendMail(event, email, message string) error {
	return api.SendNotification("mail", event, email, message)
}

func (api *ApiClient) SendNotification(channel, event, destination, message string) error {
	var apiResponse = new(ApiResponse)

	jsonData, err := json.Marshal(map[string]interface{}{
		"channel":     channel,
		"destination": []string{destination},
		"event_type":  event,
		"content":     message,
	})
	dataBytes := bytes.NewBuffer(jsonData)

	err = api.NewRequest(http.MethodPost, "/notifications", dataBytes).Send(apiResponse)

	return err
}
