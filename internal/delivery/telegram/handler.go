package telegram

import (
	"log/slog"

	tele "gopkg.in/telebot.v4"
)

type handler struct {
	logger *slog.Logger
	bot    *tele.Bot
}

func New(logger *slog.Logger, bot *tele.Bot) *handler {
	return &handler{logger, bot}
}

func (h *handler) Handle() {}
