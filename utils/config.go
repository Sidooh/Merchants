package utils

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Environment string
	JwtKey      string
	Port        int
}

func SetupConfig(path string) {
	// Set the path to look for the configurations file
	viper.AddConfigPath(path)

	// Set the file name of the configurations file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			log.Fatal("Fatal error: ", err)
		}
	}
}
