package handlers

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/service"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strconv"
	_ "strconv"
	"time"
)

type Handler struct { //handler
	s *service.Service
}

//NewHandler :define new handlers
func NewHandler(NewS *service.Service) *Handler {
	return &Handler{s: NewS}
}

type Response struct {
	Message  string
	FileName string
	FileType string
	FileSize int64
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
	return c.JSON(http.StatusOK, Response{
		Message:  "Success",
		FileName: fileName,
		FileType: fileType,
		FileSize: file.Size,
	})
}
