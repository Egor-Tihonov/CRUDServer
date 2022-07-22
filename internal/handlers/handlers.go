package handlers

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/service"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	_ "strconv"
)

type Handler struct { //handler
	s *service.Service
}

//NewHandler :define new handlers
func NewHandler(NewS *service.Service) *Handler {
	return &Handler{s: NewS}
}

//UpdateUser handler:
func (h *Handler) UpdateUser(c echo.Context) error {
	person := model.Person{}
	id := c.Param("id")
	err := json.NewDecoder(c.Request().Body).Decode(&person)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = h.s.UpdateUser(c.Request().Context(), id, &person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}
	return c.JSON(http.StatusOK, "Ok")
}

func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.s.DeleteUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "")
}

func (h *Handler) GetAllUsers(c echo.Context) error {
	p, err := h.s.SelectAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}
func (h *Handler) GetUserById(c echo.Context) error {
	id := c.Param("id")
	user, err := h.s.GetUserById(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) DownloadImage(c echo.Context) error {
	return c.File("https://ichef.bbci.co.uk/news/976/cpsprodpb/79F2/production/_123381213_06.jpg")
}
