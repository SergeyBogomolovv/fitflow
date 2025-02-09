package telegram

import (
	tele "gopkg.in/telebot.v4"
)

const (
	cmdStart       = "/start"
	cmdAbout       = "/about"
	cmdSubscribe   = "/subscribe"
	cmdUnsubscribe = "/unsubscribe"
	cmdChangeLvl   = "/change_lvl"
	cmdTest        = "/test"
)

const (
	startMessage = ""
	aboutMessage = ""
)

type Question struct {
	Question string
	Buttons  []tele.ReplyButton
	Scores   map[string]int
}

var (
	questions = []Question{
		{
			Question: "Как часто вы тренируетесь в зале?",
			Buttons:  []tele.ReplyButton{{Text: "1 Раз в неделю"}, {Text: "2 Раза в неделю"}, {Text: "3 и более раз в неделю"}},
			Scores:   map[string]int{"1 Раз в неделю": 1, "2 Раза в неделю": 2, "3 и более раз в неделю": 3},
		},
		{
			Question: "Сколько лет вы занимаетесь силовыми тренировками?",
			Buttons:  []tele.ReplyButton{{Text: "Меньше 6 месяцев"}, {Text: "1-2 Года"}, {Text: "3 и более лет"}},
			Scores:   map[string]int{"Меньше 6 месяцев": 1, "1-2 Года": 2, "3 и более лет": 3},
		},
		{
			Question: "Какие веса используете в базовых упражнениях?",
			Buttons: []tele.ReplyButton{
				{Text: "Только с собственным весом или легкими гантелями"},
				{Text: "Средние веса (50-70% от собственного веса)"},
				{Text: "Тяжелые веса (больше 100% собственного веса)"},
			},
			Scores: map[string]int{
				"Только с собственным весом или легкими гантелями": 1,
				"Средние веса (50-70% от собственного веса)":       2,
				"Тяжелые веса (больше 100% собственного веса)":     3,
			},
		},
	}
)
