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

	slog.Info("Registering new device...")

	go controller.RegisterDevice(router, &storage)

	slog.Info("Register new device successfully")

	// get devices list

	slog.Info("Getting all devices...")

	go controller.GetDeviceListResult(router, &storage)

	slog.Info("Getting all devices successfully")

	// get device by id

	slog.Info("Getting device by id...")

	go controller.GetDevice(router, &storage)

	slog.Info("Getting device by id successfully")

	// update device

	slog.Info("Updating device with id...")

	go controller.UpdateDevice(router, &storage)

	slog.Info("Device updated successfully...")

	// delete device

	slog.Info("Deleting device with id...")

	go controller.DeleteDevice(router, &storage)

	slog.Info("Device deleted successfully")

	// device monitoring data

	slog.Info("Getting device monitoring info by id...")

	go controller.GetMonitoringResult(router, &storage)

	slog.Info("Getting device monitoring info by id successfully")

	// Get Server

	server := server_utils.GetServer(router, config)

	// start server

	server_utils.StartServer(&server, config)

	// shut down server

	server_utils.StopServer(&server)
}
