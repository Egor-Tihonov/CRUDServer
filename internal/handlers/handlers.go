package handlers

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
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
	p, err := h.rps.SelectAll()
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
