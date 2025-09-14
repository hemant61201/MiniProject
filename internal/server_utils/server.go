package server_utils

import (
	"MiniProject/internal/config"
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func StartServer(server *http.Server, config *config.Config) {

	slog.Info("server_utils started", slog.String("address", config.Addr))

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("failed to start server_utils")
		}
	}()
}

func StopServer(server *http.Server) {

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	slog.Info("shutting down the server_utils")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server_utils", slog.String("error", err.Error()))
	}

	slog.Info("server_utils shutdown successfully")
}

func GetServer(router *gin.Engine, config *config.Config) http.Server {

	server := http.Server{
		Addr:    config.Addr,
		Handler: router.Handler(),
	}

	return server
}
