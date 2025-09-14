package main

import (
	"MiniProject/internal/config"
	"MiniProject/internal/controller"
	"MiniProject/internal/server_utils"
	"MiniProject/internal/storage/sqlite"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {

	// load config

	slog.Info("Loading config...")

	config := config.MustLoad()

	slog.Info("Config loaded successfully.")

	// setup database

	slog.Info("Connecting to database...")

	storage, err := sqlite.NewSqlite(config)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Database initialized", slog.String("env", config.Env), slog.String("version", "1.0.0"))

	// setup router

	router := gin.New()

	// register device

	go controller.RegisterDevice(router, &storage)

	// get devices list

	go controller.GetDeviceListResult(router, &storage)

	// get device by id

	go controller.GetDevice(router, &storage)

	// update device

	go controller.UpdateDevice(router, &storage)

	// delete device

	go controller.DeleteDevice(router, &storage)

	// device monitoring data

	go controller.GetMonitoringResult(router, &storage)

	// Get Server

	server := server_utils.GetServer(router, config)

	// start server

	server_utils.StartServer(&server, config)

	// shut down server

	server_utils.StopServer(&server)
}
