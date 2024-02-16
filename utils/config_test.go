package utils

import (
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestSetupConfig(t *testing.T) {
	err := os.Setenv("TEST", "test")
	if err != nil {
		t.Errorf("Setting env failed, got: %s, want: %v.", err, nil)
	}

	SetupConfig(".")

	test := viper.GetString("TEST")
	if test != "test" {
		t.Errorf("Incorrect env, got: %s, want: %v.", test, "test")
	}

	test = viper.Get("TEST").(string)
	if test != "test" {
		t.Errorf("Incorrect env, got: %s, want: %v.", test, "test")
	}
}
