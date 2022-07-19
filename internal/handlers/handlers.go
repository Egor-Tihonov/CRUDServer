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

type Handler struct {
	s *service.Service
}
type Id struct {
	Id string `json,bson:"id"`
}

//NewHandler :define new handlers
func NewHandler(NewS *service.Service) *Handler {
	return &Handler{s: NewS}
}

//CreateUser handler: create new model.person and read information about it from JSON
func (h *Handler) CreateUser(c echo.Context) error {
	person := model.Person{}
	err := json.NewDecoder(c.Request().Body).Decode(&person)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err, newId := h.s.CreateUser(c.Request().Context(), &person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, Id{Id: newId})
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
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "")
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
	person := model.Person{}
	id := c.Param("id")
	err := json.NewDecoder(c.Request().Body).Decode(&person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("handlers: cannot decode json file"))
	}
	p, err := h.s.GetUserById(c.Request().Context(), id, person.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err) //return c.JSON(http.StatusOk, err)
	}
	return c.JSON(http.StatusOK, p)

}
