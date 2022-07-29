package handlers

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/service"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var validate = validator.New()

type Handler struct { // handler
	s *service.Service
}

// NewHandler :define new handlers
func NewHandler(newS *service.Service) *Handler {
	return &Handler{s: newS}
}

// UpdateUser godoc
// @Summary     UpdateUser
// @Description UpdateUser is echo handler which delete user from cache and db
// @Param       id  path string true "Account ID"
// @Produce     string
// @Tags        User
// @Router      /users/{id} [delete]
// @Failure     500 string
// @Success     200 string
func (h *Handler) UpdateUser(c echo.Context) error {
	person := model.Person{}
	id := c.Param("id")
	err := ValidateValueID(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	err = json.NewDecoder(c.Request().Body).Decode(&person)
	if err != nil {
		log.Errorf("failed parse json, %e", err)
		return err
	}
	err = h.s.UpdateUser(c.Request().Context(), id, &person)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "Ok")
}

// DeleteUser godoc
// @Summary     DeleteUser
// @Description DeleteUser is echo handler which delete user from cache and db
// @Param       id path string true "Account ID"
// @Produce     string
// @Tags        User
// @Router      /users/{id} [delete]
// @Failure     500 json
// @Success     200 string
func (h *Handler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := ValidateValueID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = h.s.DeleteUser(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.String(http.StatusOK, "delete")
}

// GetAllUsers godoc
// @Summary     GetAllUsers
// @Description GetAllUsers is echo handler which returns json structure of Users objects
// @Produce     json
// @Tags        User
// @Router      /users [get]
// @Failure     500 json
// @Success     200 json
func (h *Handler) GetAllUsers(c echo.Context) error {
	p, err := h.s.SelectAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, p)
}

// GetUserByID godoc
// @Summary     GetUserByID
// @Description GetUserByID is echo handler which returns json structure of User object
// @Produce     json
// @Tags        User
// @Param       id path string true "Account ID"
// @Success     200 json
// @Failure     500 json
// @Router      /users/{id} [get]
// @Security    ApiKeyAuth
func (h *Handler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	err := ValidateValueID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	person, err := h.s.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, person)
}

func (h *Handler) DownloadFile(c echo.Context) error {
	filename := c.QueryString()
	err := c.Attachment(filename, "new_txt_file.txt")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, nil)
}
func (h *Handler) Upload(c echo.Context) error {
	var fileName, fileType string
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fileByte, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fileType = http.DetectContentType(fileByte)
	a := 5
	fileName = "uploads/" + strconv.FormatInt(time.Now().Unix(), a) + ".jpg"
	err = os.WriteFile(fileName, fileByte, 0600)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, model.Response{
		Message:  "Success",
		FileType: fileType,
		FileSize: file.Size,
	})
}
func ValidateValueID(id string) error {
	err := validate.Var(id, "required")
	if err != nil {
		return fmt.Errorf("id length couldnt be less then 36,~%v", err)
	}
	return nil
}
