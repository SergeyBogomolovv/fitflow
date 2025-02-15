package telegram

import (
	"context"

	tele "gopkg.in/telebot.v4"
)

func (h *handler) handleSubscribe(c tele.Context) error {
	const errText = "Произошла ошибка при подписке на рассылку."
	userID := c.Sender().ID

	ctx := context.TODO()
	if err := h.users.EnsureUserExists(ctx, userID); err != nil {
		return c.Send(errText)
	}

	if err := h.users.UpdateSubscribed(ctx, userID, true); err != nil {
		return c.Send(errText)
	}
	return c.Send("Вы подписались на рассылку.")
}

func (h *handler) handleUnsubscribe(c tele.Context) error {
	const errText = "Произошла ошибка при отписке от рассылки."
	userID := c.Sender().ID

	ctx := context.TODO()
	if err := h.users.EnsureUserExists(ctx, userID); err != nil {
		return c.Send(errText)
	}

	if err := h.users.UpdateSubscribed(ctx, userID, false); err != nil {
		return c.Send(errText)
	}
	return c.Send("Вы отписались от рассылки.")
}
