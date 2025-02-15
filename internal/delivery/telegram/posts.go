package telegram

import (
	"context"
	"log/slog"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/robfig/cron/v3"
	tele "gopkg.in/telebot.v4"
)

func (h *handler) RunScheduler(ctx context.Context, broadcastSpec, levelSpec string) {
	const op = "telegram.RunScheduler"
	logger := h.logger.With(slog.String("op", op))
	logger.Info("starting posts scheduler")

	cron := cron.New(cron.WithSeconds())

	broadcastID, err := cron.AddFunc(broadcastSpec, func() {
		h.notifySubscribers(ctx, domain.UserLvlDefault)
	})
	if err != nil {
		logger.Error("failed to add cron job", "error", err, "id", broadcastID)
		return
	}

	lvlID, err := cron.AddFunc(levelSpec, func() {
		h.notifySubscribers(ctx, domain.UserLvlBeginner)
		h.notifySubscribers(ctx, domain.UserLvlIntermediate)
		h.notifySubscribers(ctx, domain.UserLvlAdvanced)
	})
	if err != nil {
		logger.Error("failed to add cron job", "error", err, "id", lvlID)
		return
	}

	cron.Start()
}

func (h *handler) notifySubscribers(ctx context.Context, lvl domain.UserLvl) {
	const op = "telegram.notifySubscribers"
	logger := h.logger.With(slog.String("op", op))

	post, err := h.posts.PickLatest(ctx, lvl)
	if err != nil {
		return
	}
	subscribers, err := h.users.SubscribersIds(ctx, lvl)
	if err != nil {
		return
	}
	if count := h.sendPost(subscribers, post); count > 0 {
		h.posts.MarkAsPosted(ctx, post.ID)
		logger.Info("notified subscribers", "count", count)
	}
}

func (h *handler) sendPost(subscribers []int64, post domain.Post) int {
	const op = "telegram.sendPost"
	logger := h.logger.With(slog.String("op", op))

	count := 0
	for _, id := range subscribers {
		if err := h.sendMessage(tele.ChatID(id), post); err != nil {
			logger.Error("failed to send post", "subscriber_id", id, "error", err)
		} else {
			count++
		}
	}
	return count
}

func (h *handler) sendMessage(chatID tele.ChatID, post domain.Post) error {
	if len(post.Images) > 0 {
		var album tele.Album
		for _, url := range post.Images {
			album = append(album, &tele.Photo{File: tele.FromURL(url)})
		}
		album.SetCaption(post.Content)
		_, err := h.bot.SendAlbum(chatID, album, tele.ModeMarkdown)
		return err
	}
	_, err := h.bot.Send(chatID, post.Content, tele.ModeMarkdown)
	return err
}
