package handlers

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/service"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	_ "strconv"
)

type Handler struct {
	s *service.Service
}
type Id struct {
	Id string `json,bson:"id"`
}
type Authentication struct {
	Password string `json,bson:"password"`
}

//NewHandler :define new handlers
func NewHandler(NewS *service.Service) *Handler {
	return &Handler{s: NewS}
}

//Registration : create new model.person and read information about it from JSON
func (h *Handler) Registration(c echo.Context) error {
	person := model.Person{}
	err := json.NewDecoder(c.Request().Body).Decode(&person)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err, newId := h.s.Registration(c.Request().Context(), &person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, Id{Id: newId})
}

//UpdateUser handler:
func (h *Handler) UpdateUser(c echo.Context) error {
	refreshTokenString := c.QueryString()

	newAccessTokenString, newRefreshTokenString, err := h.s.RefreshToken(c.Request().Context(), refreshTokenString)
	if err != nil {
		log.Errorf("handler: token refresh failed - %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error while creating tokens")
	}
	person := model.Person{}
	id := c.Param("id")
	err = json.NewDecoder(c.Request().Body).Decode(&person)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = h.s.UpdateUser(c.Request().Context(), id, &person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf(`{
			"accessToken" : %v,
			"refreshToken" : %v}`, newAccessTokenString, newRefreshTokenString),
		),
	)
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

func (h *Handler) Authentication(c echo.Context) error {
	auth := Authentication{}
	id := c.Param("id")
	err := json.NewDecoder(c.Request().Body).Decode(&auth)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("handlers: cannot decode json file"))
	}
	accessToken, refreshToken, err := h.s.Authentication(c.Request().Context(), id, auth.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err) //return c.JSON(http.StatusOk, err)
	}
	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf(`{
			"refreshToken":%v,
			"accessToken" : %v}`, refreshToken, accessToken),
		),
	)

}
func (h *Handler) RefreshToken(c echo.Context) error {
	refreshTokenString := c.QueryString()

	newAccessTokenString, newRefreshTokenString, err := h.s.RefreshToken(c.Request().Context(), refreshTokenString)
	if err != nil {
		log.Errorf("handler: token refresh failed - %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error while creating tokens")
	}
	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf(`{
			"accessToken" : %v,
			"refreshToken" : %v}`, newAccessTokenString, newRefreshTokenString),
		),
	)
}
