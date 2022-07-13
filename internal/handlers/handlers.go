package handlers

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	rps *repository.Repository
}

//NewHandler :define new handlers
func NewHandler(rps *repository.Repository) *Handler {
	return &Handler{rps: rps}
}

//CreateUser handler: create new model.person and read information about it from JSON
func (h *Handler) CreateUser(c echo.Context) error {
	person := model.Person{}
	err := json.NewDecoder(c.Request().Body).Decode(&person)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = h.rps.Create(c.Request().Context(), &person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Successfully create")
}

//UpdateUser handler:
func (h *Handler) UpdateUser(c echo.Context) error {
	person := model.Person{}
	id, _ := strconv.Atoi(c.Param("id"))
	err := json.NewDecoder(c.Request().Body).Decode(&person)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = h.rps.Update(c.Request().Context(), id, &person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Successfully update")
}

func (h *Handler) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.rps.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Successfully delete")
}

func (h *Handler) GetAllUsers(c echo.Context) error {
	p, err := h.rps.SelectAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

func (h *Handler) GetUserById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	p, err := h.rps.SelectById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)

}
