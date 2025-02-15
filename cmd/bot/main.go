package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/SergeyBogomolovv/fitflow/config"
	"github.com/SergeyBogomolovv/fitflow/internal/delivery/telegram"
	postRepo "github.com/SergeyBogomolovv/fitflow/internal/repo/post"
	userRepo "github.com/SergeyBogomolovv/fitflow/internal/repo/user"
	postSvc "github.com/SergeyBogomolovv/fitflow/internal/service/post"
	userSvc "github.com/SergeyBogomolovv/fitflow/internal/service/user"
	"github.com/SergeyBogomolovv/fitflow/pkg/bot"
	"github.com/SergeyBogomolovv/fitflow/pkg/db"
	"github.com/SergeyBogomolovv/fitflow/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	confPath := flag.String("config", "./config/config.yml", "path to config file")
	flag.Parse()
	conf := config.MustNewConfig(*confPath)

	logger := logger.MustNew(conf.Log.Level, os.Stdout)

	db := db.MustNew(conf.PG.URL)
	logger.Info("database connected")

	bot := bot.MustNew(conf.TG.Token)
	logger.Info("telegram connected")

	userRepo := userRepo.New(db)
	postsRepo := postRepo.New(db)
	logger.Info("init repositories")

	userSvc := userSvc.New(logger, userRepo)
	postSvc := postSvc.New(logger, postsRepo)
	logger.Info("init services")

	telegram := telegram.New(logger, bot, postSvc, userSvc)
	telegram.Init()
	logger.Info("init handlers")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		bot.Stop()
		db.Close()
		logger.Info("bot stopped")
	}()

	logger.Info("starting bot", slog.String("name", bot.Me.FirstName))
	telegram.RunScheduler(ctx, conf.TG.BroadcastSpec, conf.TG.LevelSpec)
	bot.Start()
	wg.Wait()
}

func init() {
	godotenv.Load()
}
