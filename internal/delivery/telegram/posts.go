package telegram

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"gopkg.in/telebot.v4"
	tele "gopkg.in/telebot.v4"
)

func (h *handler) RunScheduler(ctx context.Context, delay time.Duration) {
	const op = "telegram.RunScheduler"
	logger := h.logger.With(slog.String("op", op))
	logger.Info("starting notification service", "delay", delay)

	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("shutting down notification service")
			return
		case <-ticker.C:
			h.notifySubscribers(ctx, domain.UserLvlDefault)
			h.notifySubscribers(ctx, domain.UserLvlBeginner)
			h.notifySubscribers(ctx, domain.UserLvlIntermediate)
			h.notifySubscribers(ctx, domain.UserLvlAdvanced)
		}
	}
}

func (h *handler) notifySubscribers(ctx context.Context, lvl domain.UserLvl) {
	const op = "telegram.NotifySubscribers"
	logger := h.logger.With(slog.String("op", op))

	post, err := h.posts.PickLatest(ctx, lvl)
	if err != nil {
		if errors.Is(err, domain.ErrNoPosts) {
			logger.Debug("no posts")
			return
		}
		logger.Error("failed to pick latest post", "error", err)
		return
	}

	var subscribers []int64
	if lvl == domain.UserLvlDefault {
		subscribers, err = h.users.SubscribersIds(ctx)
		if err != nil {
			return
		}
	} else {
		subscribers, err = h.users.SubscribersIdsByLvl(ctx, lvl)
		if err != nil {
			return
		}
	}

	count := h.sendPost(subscribers, post)
	if count > 0 {
		h.posts.MarkAsPosted(ctx, post.ID)
		logger.Info("notified subscribers", "count", count)
	}
}

func (h *handler) sendPost(subscribers []int64, post *domain.Post) int {
	const op = "telegram.sendPost"
	logger := h.logger.With(slog.String("op", op))

	count := 0
	if len(post.Images) > 0 {
		var album tele.Album
		for _, url := range post.Images {
			album = append(album, &tele.Photo{File: tele.FromURL(url)})
		}
		album.SetCaption(post.Content)
		for _, id := range subscribers {
			if _, err := h.bot.SendAlbum(telebot.ChatID(id), album, tele.ModeMarkdown); err != nil {
				logger.Error("failed to send post to subscriber", "id", id, "err", err)
			} else {
				count++
			}
		}
	} else {
		for _, id := range subscribers {
			if _, err := h.bot.Send(telebot.ChatID(id), post.Content, tele.ModeMarkdown); err != nil {
				logger.Error("failed to send post to subscriber", "id", id, "err", err)
			} else {
				count++
			}
		}
	}
	return count
}
