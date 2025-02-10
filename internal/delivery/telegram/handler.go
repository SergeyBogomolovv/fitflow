package telegram

import (
	"context"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/SergeyBogomolovv/fitflow/pkg/state"
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
	state  state.State
}

func New(logger *slog.Logger, bot *tele.Bot, users UserService) *handler {
	state := state.NewState()
	return &handler{logger, bot, users, state}
}

func (h *handler) Handle() {
	h.bot.Handle(cmdStart, h.handleStart)
	h.bot.Handle(cmdAbout, h.handleAbout)
	h.bot.Handle(cmdSubscribe, h.handleSubscribe)
	h.bot.Handle(cmdUnsubscribe, h.handleUnsubscribe)
	h.bot.Handle(cmdTest, h.handleStartTest)
	h.bot.Handle(tele.OnText, h.handleText)
	h.bot.Handle(cmdCancel, h.handleCancel)
}

func (h *handler) handleStart(c tele.Context) error {
	const op = "telegram.handleStart"
	userID := c.Sender().ID

	logger := h.logger.With(slog.String("op", op), slog.Int64("id", userID))

	if err := h.users.SaveUser(context.TODO(), userID); err != nil {
		logger.Error("failed to save user")
	}
	return c.Send(startMessage, defaultKeyboard)
}

func (h *handler) handleAbout(c tele.Context) error {
	return c.Send(aboutMessage)
}

func (h *handler) handleText(c tele.Context) error {
	userID := c.Sender().ID
	state := h.state.Get(userID)
	switch state := state.(type) {
	case *UserTestState:
		return h.handleTestAnswer(c, state)
	default:
		return c.Send(unknownMessage)
	}
}

func (h *handler) handleCancel(c tele.Context) error {
	h.state.Delete(c.Sender().ID)
	return c.Send("Действие отменено.", defaultKeyboard)
}

func (h *handler) handleSubscribe(c tele.Context) error {
	const op = "telegram.handleSubscribe"
	userID := c.Sender().ID

	logger := h.logger.With(slog.String("op", op), slog.Int64("id", userID))

	ctx := context.TODO()
	if err := h.users.SaveUser(ctx, userID); err != nil {
		logger.Error("failed to save user")
	}

	if err := h.users.UpdateSubscribed(ctx, userID, true); err != nil {
		logger.Error("failed to subscribe user")
		return c.Send("Произошла ошибка при подписке на рассылку.")
	}
	return c.Send("Вы подписались на рассылку.")
}

func (h *handler) handleUnsubscribe(c tele.Context) error {
	const op = "telegram.handleUnsubscribe"
	userID := c.Sender().ID

	logger := h.logger.With(slog.String("op", op), slog.Int64("id", userID))

	ctx := context.TODO()
	if err := h.users.SaveUser(ctx, userID); err != nil {
		logger.Error("failed to save user")
	}

	if err := h.users.UpdateSubscribed(ctx, userID, false); err != nil {
		logger.Error("failed to unsubscribe user")
		return c.Send("Произошла ошибка при отписки от рассылки.")
	}
	return c.Send("Вы отписались от рассылки.")
}
