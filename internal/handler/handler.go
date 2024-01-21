package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/begenov/effective-mobile-task/internal/logger"
	"github.com/begenov/effective-mobile-task/internal/model"
	"github.com/begenov/effective-mobile-task/internal/service"
	"github.com/go-chi/chi"
)

type Handler struct {
	userService *service.Service
}

func New(userService *service.Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/users", h.getUsersHandler)
	r.Post("/users", h.createUserHandler)
	r.Put("/users/{id}", h.updateUserHandler)
	r.Delete("/users/{id}", h.deleteUserHandler)

	return r
}

func (h *Handler) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(r.URL.Query().Get("userID"), 10, 64)
	gender, _ := strconv.Atoi(r.URL.Query().Get("gender"))
	nationality := r.URL.Query().Get("nationality")
	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	offset, _ := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)

	logger.Infof("getUsersHandler(): userID=%d, gender=%d, nationality=%s, limit=%d, offset=%d",
		userID, gender, nationality, limit, offset)

	users, err := h.userService.GetUsers(r.Context(), &userID, &limit, &offset, &gender, &nationality)
	if err != nil {
		logger.Error("h.userService.GetUsers(): ", err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			h.error(w, http.StatusNotFound, "Users not found")
		default:
			h.error(w, http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}

	h.respond(w, http.StatusOK, users)
}

func (h *Handler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser model.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		h.error(w, http.StatusBadRequest, model.ErrBadRequest.Error())
		return
	}

	err = h.userService.CreateUser(r.Context(), newUser)
	if err != nil {
		logger.Error("h.userService.CreateUser(): ", err)
		switch {
		case errors.Is(err, model.ErrBadRequest):
			h.error(w, http.StatusBadRequest, err.Error())
		default:
			h.error(w, http.StatusInternalServerError, model.ErrInternalServer.Error())
		}
		return
	}

	h.respond(w, http.StatusOK, nil)
}

func (h *Handler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.error(w, http.StatusBadRequest, model.ErrBadRequestUserID.Error())
		return
	}

	err = h.userService.DeleteUser(r.Context(), userID)
	if err != nil {
		logger.Error("h.userService.DeleteUser(): ", err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			h.error(w, http.StatusNotFound, model.ErrNotFound.Error())
		default:
			h.error(w, http.StatusInternalServerError, model.ErrInternalServer.Error())
		}
		return
	}

	h.respond(w, http.StatusOK, nil)
}

func (h *Handler) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		h.error(w, http.StatusBadRequest, model.ErrBadRequestUserID.Error())
		return
	}

	var updatedUser model.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		h.error(w, http.StatusBadRequest, model.ErrBadRequest.Error())
		return
	}

	updatedUser.ID = userID
	err = h.userService.UpdateUser(r.Context(), updatedUser)
	if err != nil {
		logger.Error("h.userService.UpdateUser(): ", err)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			h.error(w, http.StatusNotFound, model.ErrNotFound.Error())
		case errors.Is(err, model.ErrBadRequest):
			h.error(w, http.StatusBadRequest, err.Error())
		default:
			h.error(w, http.StatusInternalServerError, model.ErrInternalServer.Error())
		}
		return
	}

	h.respond(w, http.StatusOK, nil)
}
