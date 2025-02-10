package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	tele "gopkg.in/telebot.v4"
)

func (h *handler) handleStartTest(c tele.Context) error {
	const op = "telegram.handleStartTest"
	userID := c.Sender().ID

	logger := h.logger.With(slog.String("op", op), slog.Int64("id", userID))

	ctx := context.TODO()

	if err := h.users.SaveUser(ctx, userID); err != nil {
		logger.Error("failed to save user")
	}

	state := &UserTestState{CurrentQuestion: 0, Score: 0}
	h.state.Set(userID, state)
	return h.askQuestion(c, state)
}

func (h *handler) askQuestion(c tele.Context, state *UserTestState) error {
	if state.CurrentQuestion >= len(questions) {
		return h.finishTest(c, state)
	}

	q := questions[state.CurrentQuestion]
	markup := &tele.ReplyMarkup{}
	rows := make([]tele.Row, 0, len(q.Answers))

	for question := range q.Answers {
		rows = append(rows, tele.Row{markup.Text(question)})
	}

	markup.Reply(rows...)
	return c.Send(q.Question, markup)
}

func (h *handler) handleTestAnswer(c tele.Context, state *UserTestState) error {
	answer := c.Text()
	q := questions[state.CurrentQuestion]
	score, ok := q.Answers[answer]
	if !ok {
		return c.Send("Выберите один из предложенных вариантов.")
	}

	state.Score += score
	state.CurrentQuestion++
	return h.askQuestion(c, state)
}

func (h *handler) finishTest(c tele.Context, state *UserTestState) error {
	const op = "telegram.handleStartTest"
	userID := c.Sender().ID
	logger := h.logger.With(slog.String("op", op), slog.Int64("id", userID))

	var levelRu string
	var lvl domain.UserLvl
	avgScore := math.Round(float64(state.Score) / float64(len(questions)))
	switch avgScore {
	case 1:
		lvl = domain.UserLvlBeginner
		levelRu = "Новичок"
	case 2:
		lvl = domain.UserLvlIntermediate
		levelRu = "Средний"
	case 3:
		lvl = domain.UserLvlAdvanced
		levelRu = "Продвинутый"
	default:
		lvl = domain.UserLvlDefault
		levelRu = "Не определен"
	}

	if err := h.users.UpdateUserLvl(context.TODO(), userID, lvl); err != nil {
		logger.Error("failed to update user level")
	}

	h.state.Delete(userID)
	return c.Send(fmt.Sprintf("Ваш уровень: %s", levelRu), defaultKeyboard)
}
