package machines

import (
	"github.com/eskpil/sunlight/internal/api/mycontext"
	"github.com/eskpil/sunlight/pkg/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetAll(m *mycontext.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		var machines []models.Machine
		if err := m.Db.WithContext(c.Request().Context()).Model(&models.Machine{}).Find(&machines).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, machines)
	}
}

func GetById(m *mycontext.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		machine := new(models.Machine)

		id := c.Param("id")

		if err := m.Db.WithContext(c.Request().Context()).Model(&models.Machine{}).Preload("Verification").Find(machine, "id = ?", id).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, machine)
	}
}
