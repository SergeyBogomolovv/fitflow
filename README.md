# FitFlow Telegram bot

Telegram-бот, который помогает получать тренировочные советы и планы в зависимости от уровня подготовки.

## Что представляет из себя проект

Данный проект - это telegram бот, в котором пользователи могут пройти небольшой тест для определения уровня подготовки, чтобы получать актуальные и полезные для них посты. Посты создают администраторы бота. Используется gemini api для генерации контента для постов.

У пользователей есть следующие уровни: начинающий, средний и продвинутый.
Уровень определяется либо самим пользователем, либо через прохождение теста с вариантами ответов.
Посты могут быть как для всех, так и для определенного уровня пользователей

## Реализовано в ходе проекта

### CLI утилита для управления администраторами

- [x] Создание учетной записи администратора
- [x] Изменение пароля для администратора
- [x] Удаление администратора

### REST API для управления постами

- [x] JWT авторизация администраторов
- [x] Получение сгенерированного контента для поста
- [x] Создание постов, указывается аудитория, контент и изображения
- [ ] Изменение контента поста
- [x] Удаление поста
- [x] Отображение всех постов с фильтрами (опубликованные, неопубликованные, сортировка)

### Телеграм бот

- [x] Публикация запланированных постов
- [x] Прохождение теста для определения уровня пользователя
- [x] Подписка/Отписка от рассылки
