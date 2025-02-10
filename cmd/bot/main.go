package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/SergeyBogomolovv/fitflow/config"
	"github.com/SergeyBogomolovv/fitflow/internal/delivery/telegram"
	userRepo "github.com/SergeyBogomolovv/fitflow/internal/repo/user"
	userSvc "github.com/SergeyBogomolovv/fitflow/internal/service/user"
	"github.com/SergeyBogomolovv/fitflow/pkg/bot"
	"github.com/SergeyBogomolovv/fitflow/pkg/db"
	"github.com/SergeyBogomolovv/fitflow/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	conf := config.MustNewConfig("./config/config.yml")

	logger := logger.MustNew(conf.Log.Level, os.Stdout)
	db := db.MustNew(conf.PG.URL)
	b := bot.MustNew(conf.TG.Token)

	userRepo := userRepo.New(db)
	userSvc := userSvc.New(logger, userRepo)
	bot := telegram.New(logger, b, userSvc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		b.Stop()
		db.Close()
		logger.Info("bot stopped")
	}()

	logger.Info("starting bot", slog.String("name", b.Me.FirstName))
	bot.Handle()
	b.Start()
	wg.Wait()
}

func init() {
	godotenv.Load()
}
