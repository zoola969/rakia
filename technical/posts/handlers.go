package posts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetAllPosts(w, r)
		case http.MethodPost:
			h.CreatePost(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/posts/" || r.URL.Path == "/posts" {
			return
		}

		idStr := strings.TrimPrefix(r.URL.Path, "/posts/")

		switch r.Method {
		case http.MethodGet:
			h.GetPostByID(w, r, idStr)
		case http.MethodPut:
			h.UpdatePost(w, r, idStr)
		case http.MethodDelete:
			h.DeletePost(w, r, idStr)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// GetAllPosts handles GET /posts
// @Summary Get all posts
// @Description Get a list of all blog posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} PostRead
// @Failure 500 {object} string "Internal Server Error"
// @Router /posts [get]
func (h *Handler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetAllPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, posts)
}

// GetPostByID handles GET /posts/{id}
// @Summary Get a post by ID
// @Description Get a single blog post by its ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} PostRead
// @Failure 400 {object} string "Invalid post ID"
// @Failure 404 {object} string "Post not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /posts/{id} [get]
func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.service.GetPostByID(id)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if errors.Is(err, InvalidPostIDError) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, post)
}

// CreatePost handles POST /posts
// @Summary Create a new post
// @Description Create a new blog post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body PostCreateUpdate true "Post data"
// @Success 201 {object} PostRead
// @Failure 400 {object} string "Invalid request body or validation error"
// @Router /posts [post]
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req PostCreateUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.service.CreatePost(req)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorMessages := make([]string, len(validationErrors))
			for i, fieldError := range validationErrors {
				errorMessages[i] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fieldError.Field(), fieldError.Tag())
			}
			http.Error(w, fmt.Sprintf("Validation failed: %s", strings.Join(errorMessages, "; ")), http.StatusBadRequest)
			return
		}

		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			http.Error(w, fmt.Sprintf("Invalid validation error: %s", err.Error()), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, http.StatusCreated, post)
}

// UpdatePost handles PUT /posts/{id}
// @Summary Update a post
// @Description Update an existing blog post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body PostCreateUpdate true "Updated post data"
// @Success 200 {object} PostRead
// @Failure 400 {object} string "Invalid post ID or request body"
// @Failure 404 {object} string "Post not found"
// @Router /posts/{id} [put]
func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var req PostCreateUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := h.service.UpdatePost(id, req)
	if err != nil {
		if errors.Is(err, ErrPostNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorMessages := make([]string, len(validationErrors))
			for i, fieldError := range validationErrors {
				errorMessages[i] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fieldError.Field(), fieldError.Tag())
			}
			http.Error(w, fmt.Sprintf("Validation failed: %s", strings.Join(errorMessages, "; ")), http.StatusBadRequest)
			return
		}

		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			http.Error(w, fmt.Sprintf("Invalid validation error: %s", err.Error()), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondWithJSON(w, http.StatusOK, post)
}

// DeletePost handles DELETE /posts/{id}
// @Summary Delete a post
// @Description Delete a blog post by its ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 204 "No Content"
// @Failure 400 {object} string "Invalid post ID"
// @Failure 404 {object} string "Post not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /posts/{id} [delete]
func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeletePost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}
