package datastore

import (
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestDBInit(t *testing.T) {
	viper.Set("APP_ENV", "TEST")

	Init()

	if DB.Config.Name() != "sqlite" {
		t.Errorf("DB Init() misconfiguration: %s", DB.Config.Name())
	}

	if DB.Error != nil {
		t.Errorf("DB Init() failed: %s", DB.Error)
	}

	os.Remove("test.db")
}
