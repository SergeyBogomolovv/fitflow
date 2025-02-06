package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/SergeyBogomolovv/fitflow/internal/config"
	"github.com/SergeyBogomolovv/fitflow/internal/db"
	"github.com/SergeyBogomolovv/fitflow/internal/delivery/cli"
	repo "github.com/SergeyBogomolovv/fitflow/internal/repo/admin"
	svc "github.com/SergeyBogomolovv/fitflow/internal/service/admin"
)

func main() {
	conf := config.NewCLIConfig()
	db := db.MustNew(conf.PostgresURL)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	repo := repo.NewAdminRepo(db)
	svc := svc.NewAdminService(logger, repo)
	app := cli.NewAdminCLI(svc)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Run(ctx)
}
