package users

import (
	"github.com/eskpil/sunlight/internal/ccontroller/mycontext"
	"github.com/eskpil/sunlight/pkg/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetAll(m *mycontext.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		var users []models.User
		if err := m.Db.WithContext(c.Request().Context()).Model(&models.User{}).Find(&users).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, users)
	}
}

func GetById(m *mycontext.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(models.User)

		id := c.Param("id")

		if err := m.Db.WithContext(c.Request().Context()).Model(&models.User{}).Preload("Keys").Preload("Roles.Group").Find(user, "id = ?", id).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, user)
	}
}
