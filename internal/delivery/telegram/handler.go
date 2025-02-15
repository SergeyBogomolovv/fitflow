package telegram

import (
	"context"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/SergeyBogomolovv/fitflow/pkg/state"
	tele "gopkg.in/telebot.v4"
)

type UserService interface {
	EnsureUserExists(ctx context.Context, id int64) error
	UpdateSubscribed(ctx context.Context, id int64, subscribed bool) error
	UpdateUserLvl(ctx context.Context, id int64, lvl domain.UserLvl) error
	SubscribersIds(ctx context.Context, lvl domain.UserLvl) ([]int64, error)
}

type PostService interface {
	PickLatest(ctx context.Context, audience domain.UserLvl) (domain.Post, error)
	MarkAsPosted(ctx context.Context, id int64) error
}

type handler struct {
	logger *slog.Logger
	bot    *tele.Bot
	users  UserService
	posts  PostService
	state  state.State
}

func New(logger *slog.Logger, bot *tele.Bot, posts PostService, users UserService) *handler {
	state := state.NewState()
	return &handler{logger, bot, users, posts, state}
}

func (h *handler) Init() {
	h.bot.Handle(cmdStart, h.handleStart)
	h.bot.Handle(cmdAbout, h.handleAbout)
	h.bot.Handle(cmdSubscribe, h.handleSubscribe)
	h.bot.Handle(cmdUnsubscribe, h.handleUnsubscribe)
	h.bot.Handle(cmdTest, h.handleStartTest)
	h.bot.Handle(tele.OnText, h.handleText)
	h.bot.Handle(cmdCancel, h.handleCancel)
}

func (h *handler) handleStart(c tele.Context) error {
	userID := c.Sender().ID

	if err := h.users.EnsureUserExists(context.TODO(), userID); err != nil {
		return c.Send("Произошла непредвиденная ошибка.")
	}
	return c.Send(startMessage, defaultKeyboard, tele.ModeMarkdown)
}

func (h *handler) handleAbout(c tele.Context) error {
	return c.Send(aboutMessage, tele.ModeMarkdown)
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
