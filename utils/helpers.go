package utils

import (
	"encoding/json"
)

func ConvertStruct(from interface{}, to interface{}) {
	record, _ := json.Marshal(from)
	_ = json.Unmarshal(record, to)
}
