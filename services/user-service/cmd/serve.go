package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"user-service/internal/bootstrap"
	"user-service/internal/common/response"
	"user-service/internal/common/validator"
	"user-service/internal/constants"
	"user-service/internal/database"
	"user-service/internal/handler"
	"user-service/internal/logger"
	"user-service/internal/middleware"
	"user-service/internal/repository"
	"user-service/internal/routes"
	"user-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start user service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer()
	},
}

func runServer() error {
	postgresDB, err := database.NewPostgres(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	defer func() {
		if err := postgresDB.Close(); err != nil {
			logger.Log.WithError(err).Warn("failed to close postgres connection")
		}
	}()

	logger.Log.WithFields(logrus.Fields{
		"host": cfg.Database.Host,
		"port": cfg.Database.Port,
		"name": cfg.Database.Name,
	}).Info("postgres connection established")

	v, err := validator.New()
	if err != nil {
		return fmt.Errorf("failed to initialize validator: %w", err)
	}

	repositoryRegistry := repository.NewRegistry(postgresDB.DB)
	serviceRegistry := services.NewRegistry(repositoryRegistry, cfg.JWT)
	handlerRegistry := handler.NewRegistry(serviceRegistry, v)

	bootstrap.SetupAppMode(cfg.App)

	router := gin.New()
	router.Use(middleware.RequestID(), middleware.RequestLogger(), middleware.HandlePanic())
	router.HandleMethodNotAllowed = true
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, response.Response{
			Status:  constants.Error,
			Message: "指定されたパスが見つかりません",
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, response.Response{
			Status:  constants.Error,
			Message: "指定されたHTTPメソッドは許可されていません",
		})
	})

	group := router.Group("/api/v1/users")
	routerRegistry := routes.NewRegistry(group, handlerRegistry)
	routerRegistry.RegisterRoutes()

	srv := &http.Server{
		Addr:    cfg.ServerAddress(),
		Handler: router,
	}

	serverErr := make(chan error, 1)

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Log.WithFields(logrus.Fields{
			"environment": cfg.App.Env,
			"address":     cfg.ServerAddress(),
		}).Info("starting user service")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case <-signalCtx.Done():
		logger.Log.Info("shutdown signal received")
	case err := <-serverErr:
		return fmt.Errorf("user service stopped unexpectedly: %w", err)
	}

	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.WithError(err).Warn("failed to gracefully shutdown user service")
		_ = srv.Close()
		return fmt.Errorf("forced to shutdown: %w", err)
	}
	logger.Log.Info("service exited gracefully")

	return nil
}
