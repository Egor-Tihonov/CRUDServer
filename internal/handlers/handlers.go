package handlers

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/service"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var validate = validator.New()

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
	err := ValidateValueID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	err = json.NewDecoder(c.Request().Body).Decode(&person)
	if err != nil {
		log.Errorf("failed parse json, %e", err)
		return err
	}
	err = h.s.UpdateUser(c.Request().Context(), id, &person)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
	}
	return c.JSON(http.StatusOK, "Ok")
}

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
	err = h.s.DeleteUserFromCache(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, "")
}

func (h *Handler) GetAllUsers(c echo.Context) error {
	users, found, err := h.s.GetAllUsersFromCache(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if !found {
		p, err := h.s.SelectAllUsers(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, p)
	}
	return c.JSON(http.StatusOK, users)
}

// GetUserById godoc
// @Summary     GetUserById
// @Description GetOrderByID is echo handler(GET) which returns json structure of User object
// @Accept json
// @Produce json
// @Tags        orders
// @Param       id  path     string true "Account ID"
// @Success     200
// @Router      /users/{id} [get]
// @Security    ApiKeyAuth
func (h *Handler) GetUserById(c echo.Context) error {
	id := c.Param("id")
	err := ValidateValueID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	user, found, err := h.s.GetUserFromCache(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if !found {
		person, err := h.s.GetUserById(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, person)
	}
	return c.JSON(http.StatusOK, user)

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
	fileByte, err := ioutil.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fileType = http.DetectContentType(fileByte)
	fileName = "uploads/" + strconv.FormatInt(time.Now().Unix(), 10) + ".jpg"
	err = ioutil.WriteFile(fileName, fileByte, 0777)
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
