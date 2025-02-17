package content

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/SergeyBogomolovv/fitflow/pkg/httpx"
	"github.com/go-playground/validator/v10"
)

type ContentService interface {
	GenerateContent(ctx context.Context, theme string) (string, error)
	CreatePost(ctx context.Context, in domain.CreatePostDTO) (domain.Post, error)
	RemovePost(ctx context.Context, id int64) error
	Posts(ctx context.Context, audience domain.UserLvl, incoming bool) ([]domain.Post, error)
}

type handler struct {
	logger     *slog.Logger
	validate   *validator.Validate
	contentSvc ContentService
}

func New(logger *slog.Logger, contentSvc ContentService) *handler {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &handler{logger, validate, contentSvc}
}

func (h *handler) Init(r *http.ServeMux) {
	r.HandleFunc("GET /content/generate", h.HandleGenerateContent)
	r.HandleFunc("POST /content/post", h.HandleCreatePost)
	r.HandleFunc("DELETE /content/post/{id}", h.HandleRemovePost)
}

// @Summary      Генерация контента для поста
// @Description  Генерирует контент для телеграм поста на заданную тему с помощью AI
// @Tags         content
// @Accept       json
// @Produce      json
// @Param 			 theme  query     string true "Тема контента"
// @Success      200    {object}  GenerateContentResponse
// @Failure      400    {object}  httpx.Response  "Неверный формат запроса"
// @Failure      500    {object}  httpx.Response  "Внутренняя ошибка сервера"
// @Router       /content/generate [get]
func (h *handler) HandleGenerateContent(w http.ResponseWriter, r *http.Request) {
	theme := r.URL.Query().Get("theme")
	if theme == "" {
		httpx.WriteError(w, "theme is required", http.StatusBadRequest)
		return
	}

	content, err := h.contentSvc.GenerateContent(r.Context(), theme)
	if err != nil {
		h.logger.Error("failed to generate content", "error", err, "theme", theme)
		httpx.WriteError(w, "failed to generate content", http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, GenerateContentResponse{Content: content, Status: httpx.StatusSuccess}, http.StatusOK)
}

// @Summary      Создание нового поста
// @Description  Сохраняет пост в бд, сохраняет изображения в s3
// @Tags         content
// @Accept 			 multipart/form-data
// @Produce      json
// @Param images formData file true "Список изображений (можно несколько)"
// @Param content formData string true "Текст поста"
// @Param audience formData string true "Аудитория (beginner, intermediate, advanced)"
// @Success      200    {object}  domain.Post
// @Failure      400    {object}  httpx.Response  "Неверные данные в запросе"
// @Failure      500    {object}  httpx.Response  "Внутренняя ошибка сервера"
// @Router       /content/post [post]
func (h *handler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		httpx.WriteError(w, "invalid payload", http.StatusBadRequest)
		return
	}

	dto := domain.CreatePostDTO{
		Content:  r.FormValue("content"),
		Images:   r.MultipartForm.File["images"],
		Audience: domain.UserLvl(r.FormValue("audience")),
	}
	if err := h.validate.Struct(dto); err != nil {
		httpx.WriteError(w, "invalid payload", http.StatusBadRequest)
		return
	}

	post, err := h.contentSvc.CreatePost(r.Context(), dto)
	if err != nil {
		h.logger.Error("error creating post", "error", err)
		httpx.WriteError(w, "failed to create post", http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, post, http.StatusCreated)
}

// @Summary      Удаление поста
// @Tags         content
// @Produce      json
// @Param        id   path      int  true  "ID поста"
// @Success      200  {object}  httpx.Response  "Пост успешно удалён"
// @Failure      400  {object}  httpx.Response  "Некорректный ID"
// @Failure      404  {object}  httpx.Response  "Пост не найден"
// @Failure      500  {object}  httpx.Response  "Внутренняя ошибка сервера"
// @Router       /content/post/{id} [delete]
func (h *handler) HandleRemovePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		httpx.WriteError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.contentSvc.RemovePost(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrPostNotFound) {
			httpx.WriteError(w, "post not found", http.StatusNotFound)
			return
		}
		httpx.WriteError(w, "failed to delete post", http.StatusInternalServerError)
		return
	}

	httpx.WriteSuccess(w, "post deleted", http.StatusOK)
}

// @Summary      Получение постов
// @Tags         content
// @Produce      json
// @Param        audience   query     string  false  "Уровень пользователя (beginner, intermediate, advanced)" default(default)
// @Param        incoming   query     boolean false  "Фильтр по публикации (true - не опубликованные, false - все)"
// @Success      200  {array}   domain.Post   "Список постов"
// @Failure      500  {object}  httpx.Response  "Внутренняя ошибка сервера"
// @Router       /content/posts [get]
func (h *handler) HandleGetPosts(w http.ResponseWriter, r *http.Request) {
	audience := r.URL.Query().Get("audience")
	if audience != "beginner" && audience != "intermediate" && audience != "advanced" {
		audience = "default"
	}
	incoming := false
	if strings.EqualFold(r.URL.Query().Get("incoming"), "true") {
		incoming = true
	}

	posts, err := h.contentSvc.Posts(r.Context(), domain.UserLvl(audience), incoming)
	if err != nil {
		httpx.WriteError(w, "failed to get posts", http.StatusInternalServerError)
		return
	}

	httpx.WriteJSON(w, posts, http.StatusOK)
}
