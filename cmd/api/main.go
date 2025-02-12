package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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

	db := db.MustNew(conf.PG.URL)
	defer db.Close()

	logger := logger.MustNew(conf.Log.Level, os.Stdout)

	adminRepo := adminRepo.New(db)
	authSvc := authSvc.New(logger, adminRepo, conf.JWT.Secret, conf.JWT.TTL)
	authHandler := authHandler.New(logger, authSvc)

	router := http.NewServeMux()
	router.Handle("/api/docs/", httpSwagger.WrapHandler)
	authHandler.Handle(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.HTTP.Host, conf.HTTP.Port),
		Handler: router,
	}

	logger.Info("starting server")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to start server: %s", err)
	}
}

func init() {
	godotenv.Load()
}
