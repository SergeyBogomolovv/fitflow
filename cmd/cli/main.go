package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/SergeyBogomolovv/fitflow/config"
	"github.com/SergeyBogomolovv/fitflow/internal/delivery/cli"
	adminRepo "github.com/SergeyBogomolovv/fitflow/internal/repo/admin"
	adminSvc "github.com/SergeyBogomolovv/fitflow/internal/service/admin"
	"github.com/SergeyBogomolovv/fitflow/pkg/db"
	"github.com/SergeyBogomolovv/fitflow/pkg/logger"
)

func main() {
	conf := config.MustNewConfig("./config/config.yml")
	db := db.MustNew(conf.PG.URL)
	defer db.Close()

	logger := logger.MustNew(conf.Log.Level, io.Discard)
	repo := adminRepo.New(db)
	svc := adminSvc.New(logger, repo)
	app := cli.NewAdminCLI(svc)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Run(ctx)
}
