package handlers

import (
	"awesomeProject/internal/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

//Registration : create new model.person and read information about it from JSON
func (h *Handler) Registration(c echo.Context) error {
	person := model.Person{}
	err := ValidateStruct(&person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = json.NewDecoder(c.Request().Body).Decode(&person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	newId, err := h.s.Registration(c.Request().Context(), &person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.String(http.StatusOK, fmt.Sprintf("You register with "+`{"ID":%v}`, newId))
}

func (h *Handler) Authentication(c echo.Context) error {
	auth := model.Authentication{}
	id := c.Param("id")
	err := ValidateValueID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "id cant be empty")
	}
	err = json.NewDecoder(c.Request().Body).Decode(&auth)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("handlers: cannot decode json file"))
	}
	accessToken, refreshToken, err := h.s.Authentication(c.Request().Context(), id, auth.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("error:%v", err))
	}
	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf("You_entry_with "+`{
			"refreshToken":%v,
			"accessToken" : %v}`, refreshToken, accessToken),
		),
	)

}
func (h *Handler) RefreshToken(c echo.Context) error {
	refreshToken := model.RefreshTokens{}
	err := json.NewDecoder(c.Request().Body).Decode(&refreshToken)
	if err != nil {
		return err
	}
	newAccessTokenString, newRefreshTokenString, err := h.s.RefreshToken(c.Request().Context(), refreshToken.RefreshToken)
	if err != nil {
		log.Errorf("handler: token refresh failed - %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error while creating tokens")
	}
	return c.JSONBlob(
		http.StatusOK,
		[]byte(
			fmt.Sprintf("Tokens_refresh: "+`{
			"accessToken" : %v,
			"refreshToken" : %v}`, newAccessTokenString, newRefreshTokenString),
		),
	)
}
func (h *Handler) Logout(c echo.Context) error {
	id := c.Param("id")
	err := ValidateValueID(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "id cant be empty")
	}
	err = h.s.DeleteUserFromCache(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed delete user from cache, %e", err))
	}
	err = h.s.UpdateUserAuth(c.Request().Context(), id, "")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	h.s.DeleteFromCache()
	return c.JSON(http.StatusOK, "logout")
}
func ValidateStruct(person *model.Person) error {
	err := validate.Struct(person)
	if err != nil {
		return fmt.Errorf("error with validate user, check your name(min length = 6),password(min length = 8) and age couldnt be less then 0 or greater than 200,~ %v", err)
	}
	return nil
}
