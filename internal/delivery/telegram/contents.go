package telegram

import (
	tele "gopkg.in/telebot.v4"
)

const (
	cmdStart       = "/start"
	cmdAbout       = "/about"
	cmdSubscribe   = "/subscribe"
	cmdUnsubscribe = "/unsubscribe"
	cmdTest        = "/test"
	cmdCancel      = "/cancel"
)

const (
	startMessage = "💥Дорогой пользователь, мы рады что вы решили менять свою жизнь выбирая работу над собой, мы же в свою очередь поможем вам с этим 🏋🏿‍♂️\n\n" +
		"📝 Пройди тест для определения уровня - /test\n\n" +
		"📢 Иногда будем делиться крутыми предложениями.\n\n" +
		"🔔 Чтобы получать наши посты, нажми 👉 /subscribe\n" +
		"❌ Чтобы отписаться в любой момент – /unsubscribe"
	aboutMessage   = "Используй /subscribe /unsubscribe /test"
	unknownMessage = "Команда не распознана. Если хотите пройти тест, используйте - /test"
)

type Question struct {
	Question string
	Answers  map[string]int
}

var questions = []Question{
	{
		Question: "Как часто вы тренируетесь в зале?",
		Answers:  map[string]int{"1 Раз в неделю": 1, "2 Раза в неделю": 2, "3 и более раз в неделю": 3},
	},
	{
		Question: "Сколько лет вы занимаетесь силовыми тренировками?",
		Answers:  map[string]int{"Меньше 6 месяцев": 1, "1-2 Года": 2, "3 и более лет": 3},
	},
	{
		Question: "Какие веса используете в базовых упражнениях?",
		Answers: map[string]int{
			"Только с собственным весом или легкими гантелями": 1,
			"Средние веса (50-70% от собственного веса)":       2,
			"Тяжелые веса (больше 100% собственного веса)":     3,
		},
	},
}

var (
	defaultKeyboard = &tele.ReplyMarkup{RemoveKeyboard: true}
)
