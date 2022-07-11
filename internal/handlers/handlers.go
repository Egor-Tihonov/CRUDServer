package handlers

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	rps *repository.Repository
}

func NewHandler(rps *repository.Repository) *Handler {
	return &Handler{rps: rps}
}

func (h *Handler) CreateUser(c echo.Context) error {
	person := model.Person{}
	err := h.rps.Create(c.Request().Context(), &person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "Successfully create")
}

func (h *Handler) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.rps.Update(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK /*???*/, h.rps.SelectAll() /*???*/)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.rps.Delete(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, h.rps.SelectAll())
}

func (h *Handler) GetAllUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, h.rps.SelectAll())
}

func (h *Handler) GetUserById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, h.rps.SelectById(id))
}
