package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/didsqq/user_api/docs"
	"github.com/didsqq/user_api/internal/domain"
	"github.com/didsqq/user_api/internal/handler/middleware"
	"github.com/didsqq/user_api/internal/handler/validate"
	"github.com/didsqq/user_api/internal/service"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	r.Route("/api/users", func(r chi.Router) {
		r.Post("/", h.createUser)
		r.Get("/{id}", h.getUser)
		r.Delete("/{id}", h.deleteUser)
		r.Put("/{id}", h.updateUser)
		r.Get("/", h.getUsers)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return r
}

// @Summary     Get user by ID
// @Tags        user
// @Description Returns user data by their ID
// @Accept      json
// @Produce     json
// @Param       id   path     int    true  "User ID"
// @Success     200  {object} domain.User "User details"
// @Failure     400  {string} string "Bad Request - Invalid ID format"
// @Failure     404  {string} string "Not Found - User not found"
// @Failure     500  {string} string "Internal Server Error"
// @Router      /api/users/{id} [get]
func (h *Handler) getUser(w http.ResponseWriter, req *http.Request) {
	idParam := chi.URLParam(req, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Ошибка преобразования id", http.StatusBadRequest)
		log.Printf("Ошибка преобразования id: %s", err)
		return
	}

	ctx := req.Context()

	user, err := h.services.GetByID(ctx, int64(id))
	if err != nil {
		http.Error(w, "Ошибка получения пользователя", http.StatusBadRequest)
		log.Printf("Ошибка получения пользователя: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(&user); err != nil {
		http.Error(w, "Ошибка кодирования объекта", http.StatusInternalServerError)
		log.Printf("Ошибка кодирования объекта: %s", err)
		return
	}
}

// @Summary     Create a new user
// @Tags        user
// @Description Creates a new user by receiving user details in the request body
// @Accept      json
// @Produce     json
// @Param       req  body  domain.User  true  "User information"
// @Success     201  {string} string "User created successfully"
// @Failure     400  {string} string "Bad Request - Invalid input data"
// @Failure     500  {string} string "Internal Server Error"
// @Router      /api/users [post]
func (h *Handler) createUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var user domain.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка парсинга json", http.StatusBadRequest)
		log.Printf("Ошибка парсинга json: %s", err)
		return
	}

	if err := validate.ValidateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Ошибка валидации: %s", err)
		return
	}

	ctx := req.Context()

	err = h.services.Create(ctx, user)
	if err != nil {
		http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
		log.Printf("Ошибка создания пользователя: %s", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary     Delete a user by ID
// @Tags        user
// @Description Deletes a user from the system by the given ID
// @Accept      json
// @Produce     json
// @Param       id  path  int64  true  "User ID"
// @Success     200  {string} string "User deleted successfully"
// @Failure     400  {string} string "Bad Request - Invalid ID or deletion failure"
// @Failure     404  {string} string "Not Found - User not found"
// @Failure     500  {string} string "Internal Server Error"
// @Router      /api/users/{id} [delete]
func (h *Handler) deleteUser(w http.ResponseWriter, req *http.Request) {
	idParam := chi.URLParam(req, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Ошибка преобразования id", http.StatusBadRequest)
		log.Printf("Ошибка преобразования id: %s", err)
		return
	}

	ctx := req.Context()

	if err := h.services.Delete(ctx, int64(id)); err != nil {
		http.Error(w, "Ошибка получения пользователя", http.StatusBadRequest)
		log.Printf("Ошибка получения пользователя: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

// @Summary     Update a user by ID
// @Tags        user
// @Description Updates the user information based on the provided ID and new data
// @Accept      json
// @Produce     json
// @Param       id   path  int64   true  "User ID"
// @Param       user body  domain.UpdateUserInput true  "User data to update"
// @Success     200  {string} string "User updated successfully"
// @Failure     400  {string} string "Bad Request - Invalid data or validation failure"
// @Failure     404  {string} string "Not Found - User not found"
// @Failure     500  {string} string "Internal Server Error"
// @Router      /api/users/{id} [put]
func (h *Handler) updateUser(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	idParam := chi.URLParam(req, "id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Ошибка преобразования id", http.StatusBadRequest)
		log.Printf("Ошибка преобразования id: %s", err)
		return
	}

	var user domain.UpdateUserInput

	if err = json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(w, "Ошибка парсинга json", http.StatusBadRequest)
		log.Printf("Ошибка парсинга json: %s", err)
		return
	}

	user.ID = int64(id)

	if err := validate.ValidateUpdateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Ошибка валидации: %s", err)
		return
	}

	ctx := req.Context()

	if err := h.services.Update(ctx, user); err != nil {
		http.Error(w, "Ошибка получения пользователя", http.StatusBadRequest)
		log.Printf("Ошибка получения пользователя: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

// @Summary     Get users with optional filtering and pagination
// @Tags        user
// @Description Returns a list of users based on provided conditions, including optional pagination and filtering parameters
// @Accept      json
// @Produce     json
// @Param       limit      query   int    false "Limit number of users per page"
// @Param       offset     query   int    false "Offset for pagination"
// @Success     200        {array} domain.User "List of users"
// @Failure     400        {string} string "Bad Request - Invalid data or parsing failure"
// @Failure     500        {string} string "Internal Server Error"
// @Router      /api/users [get]  // Маршрут для получения пользователей
func (h *Handler) getUsers(w http.ResponseWriter, req *http.Request) {
	limit := req.URL.Query().Get("limit") // Например: /api/users?limit=10&offset=0
	offset := req.URL.Query().Get("offset")

	limitI, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Ошибка преобразования limit", http.StatusBadRequest)
		log.Printf("Ошибка преобразования limit: %s", err)
		return
	}

	offsetI, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Ошибка преобразования offset", http.StatusBadRequest)
		log.Printf("Ошибка преобразования offset: %s", err)
		return
	}

	c := domain.Conditions{
		Limit:  limitI,
		Offset: offsetI,
	}

	ctx := req.Context()

	users, err := h.services.List(ctx, c)
	if err != nil {
		http.Error(w, "Ошибка получения пользователей", http.StatusBadRequest)
		log.Printf("Ошибка получения пользователей: %s", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(&users); err != nil {
		http.Error(w, "Ошибка кодирования объекта", http.StatusInternalServerError)
		log.Printf("Ошибка кодирования объекта: %s", err)
		return
	}
}
