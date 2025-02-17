package content

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
	"github.com/SergeyBogomolovv/fitflow/pkg/httpx"
	"github.com/go-playground/validator/v10"
)

type ContentService interface {
	GenerateContent(ctx context.Context, theme string) (string, error)
	CreatePost(ctx context.Context, in domain.CreatePostDTO) (domain.Post, error)
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
// @Success      200    {object}  GenerateContentResponse
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
