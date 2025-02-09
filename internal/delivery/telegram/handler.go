package telegram

import (
	"context"
	"errors"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	tele "gopkg.in/telebot.v4"
)

type UserService interface {
	SaveUser(ctx context.Context, id int64) error
	UpdateSubscribed(ctx context.Context, id int64, subscribed bool) error
	UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error
	SubscribersIdsByLvl(ctx context.Context, lvl domain.UserLvl) ([]int64, error)
	SubscribersIds(ctx context.Context) ([]int64, error)
}
type handler struct {
	logger *slog.Logger
	bot    *tele.Bot
	users  UserService
}

func New(logger *slog.Logger, bot *tele.Bot, users UserService) *handler {
	return &handler{logger, bot, users}
}

func (h *handler) Handle() {
	h.bot.Handle(cmdStart, h.handleStart)
	h.bot.Handle(cmdAbout, h.handleAbout)
	h.bot.Handle(cmdSubscribe, h.handleSubscribe)
	h.bot.Handle(cmdUnsubscribe, h.handleUnsubscribe)
}

func (h *handler) handleStart(ctx tele.Context) error {
	const op = "telegram.handleStart"
	logger := h.logger.With(slog.String("op", op))

	userId := ctx.Sender().ID

	if err := h.users.SaveUser(context.TODO(), userId); err != nil {
		logger.Error("failed to save user", "error", err)
		return errors.New("failed to save user")
	}

	return ctx.Send(startMessage)
}

func (h *handler) handleAbout(ctx tele.Context) error {
	return ctx.Send(aboutMessage)
}

func (h *handler) handleSubscribe(ctx tele.Context) error {
	if err := h.users.UpdateSubscribed(context.TODO(), ctx.Sender().ID, true); err != nil {
		return errors.New("failed to subscribe")
	}
	return ctx.Send("Вы подписались на рассылку.")
}

func (h *handler) handleUnsubscribe(ctx tele.Context) error {
	if err := h.users.UpdateSubscribed(context.TODO(), ctx.Sender().ID, false); err != nil {
		return errors.New("failed to unsubscribe")
	}
	return ctx.Send("Вы отписались от рассылки.")
}
