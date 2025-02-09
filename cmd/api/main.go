package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SergeyBogomolovv/fitflow/config"
	authHandler "github.com/SergeyBogomolovv/fitflow/internal/delivery/http/auth"
	adminRepo "github.com/SergeyBogomolovv/fitflow/internal/repo/admin"
	authSvc "github.com/SergeyBogomolovv/fitflow/internal/service/auth"
	"github.com/SergeyBogomolovv/fitflow/pkg/db"
	"github.com/SergeyBogomolovv/fitflow/pkg/logger"
)

func main() {
	conf := config.MustNewConfig("./config/config.yml")
	db := db.MustNew(conf.PG.URL)
	defer db.Close()

	logger := logger.MustNew(conf.Log.Level, os.Stdout)

	adminRepo := adminRepo.New(db)
	authSvc := authSvc.New(logger, adminRepo, conf.JWT.Secret, conf.JWT.TTL)
	authHandler := authHandler.New(logger, authSvc)

	router := http.NewServeMux()
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
