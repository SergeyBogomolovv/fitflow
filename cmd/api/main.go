package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/SergeyBogomolovv/fitflow/config"
	_ "github.com/SergeyBogomolovv/fitflow/docs"
	authHandler "github.com/SergeyBogomolovv/fitflow/internal/delivery/http/auth"
	adminRepo "github.com/SergeyBogomolovv/fitflow/internal/repo/admin"
	authSvc "github.com/SergeyBogomolovv/fitflow/internal/service/auth"
	"github.com/SergeyBogomolovv/fitflow/pkg/db"
	"github.com/SergeyBogomolovv/fitflow/pkg/logger"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title FitFlow API
// @version 0.0.1
// @description Описание API для сервиса FitFlow
func main() {
	confPath := flag.String("config", "./config/config.yml", "path to config file")
	flag.Parse()
	conf := config.MustNewConfig(*confPath)

	logger := logger.MustNew(conf.Log.Level, os.Stdout)

	db := db.MustNew(conf.PG.URL)
	logger.Info("database connected")

	router := http.NewServeMux()
	router.Handle("/api/docs/", httpSwagger.WrapHandler)

	adminRepo := adminRepo.New(db)
	logger.Info("init repositories")

	authSvc := authSvc.New(logger, adminRepo, conf.JWT.Secret, conf.JWT.TTL)
	logger.Info("init services")

	authHandler := authHandler.New(logger, authSvc)
	authHandler.Init(router)
	logger.Info("init handlers")

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.HTTP.Host, conf.HTTP.Port),
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
		db.Close()
	}()

	logger.Info("starting server")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to start server: %s", err)
	}
	wg.Wait()
}

func init() {
	godotenv.Load()
}
