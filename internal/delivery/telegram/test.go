package telegram

import (
	"context"
	"fmt"
	"math"
	"slices"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	tele "gopkg.in/telebot.v4"
)

func (h *handler) handleStartTest(c tele.Context) error {
	userID := c.Sender().ID
	ctx := context.TODO()

	if err := h.users.EnsureUserExists(ctx, userID); err != nil {
		return c.Send("Произошла непредвиденная ошибка.")
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

	type answer struct {
		text string
		num  int
	}
	answers := make([]answer, 0, len(q.Answers))
	for ans, num := range q.Answers {
		answers = append(answers, answer{ans, num})
	}
	slices.SortFunc(answers, func(a, b answer) int {
		return a.num - b.num
	})

	markup := &tele.ReplyMarkup{}
	rows := make([]tele.Row, 0, len(q.Answers))
	for _, ans := range answers {
		rows = append(rows, tele.Row{markup.Text(ans.text)})
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
	ctx := context.TODO()
	userID := c.Sender().ID

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

	if err := h.users.UpdateUserLvl(ctx, userID, lvl); err != nil {
		return c.Send("Произошла ошибка при обновлении уровня.", defaultKeyboard)
	}

	h.state.Delete(userID)
	return c.Send(fmt.Sprintf("Ваш уровень: %s", levelRu), defaultKeyboard)
}
