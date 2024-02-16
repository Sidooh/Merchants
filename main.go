package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"merchants.sidooh/api"
	"merchants.sidooh/pkg/cache"
	"merchants.sidooh/pkg/clients"
	"merchants.sidooh/pkg/datastore"
	"merchants.sidooh/pkg/logger"
	"merchants.sidooh/utils"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	utils.SetupConfig(".")

	jwtKey := viper.GetString("JWT_KEY")
	if len(jwtKey) == 0 {
		panic("JWT_KEY is not set")
	}

	logger.Init()
	datastore.Init()
	cache.Init()
	clients.Init()

	app := api.Server()

	port := viper.GetString("PORT")
	if port == "" {
		port = "8000"
	}
	go func() {
		log.Fatal(app.Listen(":" + port))
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	_ = <-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()
}
