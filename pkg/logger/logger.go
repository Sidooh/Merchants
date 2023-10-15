package logger

import (
	"github.com/spf13/viper"
	"log/slog"
	"merchants.sidooh/utils"
	"os"
)

var ClientLog = slog.Default()

func Init() {
	ClientLog = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	env := viper.GetString("APP_ENV")
	logger := viper.GetString("LOGGER")

	if env != "TEST" {
		if logger == "GCP" {
			//	ClientLog.SetFormatter(NewGCEFormatter(false))
		} else {
			ClientLog = slog.New(slog.NewJSONHandler(utils.GetLogFile("client.log"), nil))
		}
	}
}
