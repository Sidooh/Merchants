package utils

import "gorm.io/datatypes"

type ApiResponse struct {
	Result  int           `json:"result"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
	Errors  []interface{} `json:"errors"`
}

type Payment struct {
	Id          uint           `json:"id"`
	Amount      float32        `json:"amount,string"`
	Charge      int            `json:"charge"`
	Status      string         `json:"status"`
	Destination datatypes.JSON `json:"destination"`
}

type PaymentApiResponse struct {
	ApiResponse
	Data *Payment `json:"data"`
}
