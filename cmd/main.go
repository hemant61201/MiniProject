package main

import (
	"MiniProject/internal/config"
	"MiniProject/internal/status"
	"MiniProject/internal/storage/sqlite"
	"MiniProject/internal/types"
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

	router.POST("/devices", func(context *gin.Context) {

		slog.Info("Registering new device...")

		var device types.Device

		if err := context.ShouldBindJSON(&device); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		status.CheckStatus(&device)

		result, err := storage.RegisterDevice(&device)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		context.JSON(http.StatusOK, gin.H{
			"registration id": result,
		})

		slog.Info("Register new device successfully")
	})

	// get devices list

	router.GET("/devices", func(context *gin.Context) {

		slog.Info("Getting all devices...")

		result, err := storage.GetDeviceList()

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		context.JSON(http.StatusOK, gin.H{
			"devices": result,
		})

		slog.Info("Getting all devices successfully")
	})

	// get device by id

	router.GET("/devices/:id", func(context *gin.Context) {

		slog.Info("Getting device by id...")

		id, err := strconv.ParseInt(context.Param("id"), 10, 64)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		result, err := storage.GetDevice(id)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		context.JSON(http.StatusOK, gin.H{
			"device": result,
		})

		slog.Info("Getting device by id successfully")
	})

	// update device

	router.PUT("/devices/:id", func(context *gin.Context) {

		slog.Info("Updating device with id...")

		id, err := strconv.ParseInt(context.Param("id"), 10, 64)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		idPresent, err := storage.CheckDevice(id)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		if idPresent {

			var deviceInput types.UpdateDeviceInput

			if err := context.ShouldBindJSON(&deviceInput); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			}

			result, err := storage.UpdateDevice(id, &deviceInput)

			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}

			if result == 0 {
				context.JSON(http.StatusOK, gin.H{
					"error": "Device is not Updated",
				})
			} else {
				context.JSON(http.StatusOK, gin.H{
					"result": "Device is Updated successfully...",
				})
			}

		} else {
			context.JSON(http.StatusOK, gin.H{
				"error": "Device not found",
			})
		}

		slog.Info("Device updated successfully...")
	})

	// delete device

	router.DELETE("/devices/:id", func(context *gin.Context) {

		slog.Info("Deleting device with id...")

		id, err := strconv.ParseInt(context.Param("id"), 10, 64)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		result, err := storage.DeleteDevice(id)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		if result == 0 {
			context.JSON(http.StatusOK, gin.H{
				"error": "Device not found",
			})
		} else {
			context.JSON(http.StatusOK, gin.H{
				"result": "Device deleted successfully",
			})
		}

		slog.Info("Device deleted successfully")
	})

	// device monitoring data

	router.GET("/devices/:id/monitoring", func(context *gin.Context) {

		slog.Info("Getting device monitoring info by id...")

		id, err := strconv.ParseInt(context.Param("id"), 10, 64)

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		result, err := storage.GetDevice(id)

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		for i := range result {
			result[i].Metadata.CPU = fmt.Sprintf("%d%%", rand.Intn(100))
			result[i].Metadata.Memory = fmt.Sprintf("%dMB", rand.Intn(16000))
			result[i].Metadata.Disk = fmt.Sprintf("%dGB", rand.Intn(512))
		}

		context.JSON(http.StatusOK, gin.H{
			"device": result,
		})

		slog.Info("Getting device monitoring info by id successfully")

	})

	server := http.Server{
		Addr:    config.Addr,
		Handler: router.Handler(),
	}

	slog.Info("server started", slog.String("address", config.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
