package verifications

import (
	"github.com/eskpil/sunlight/internal/api/mycontext"
	"github.com/eskpil/sunlight/pkg/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetAll(m *mycontext.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		var verifications []models.Verification

		if err := m.Db.WithContext(c.Request().Context()).Model(&models.Verification{}).Find(&verifications).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, verifications)
	}
}

func GetById(m *mycontext.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		var verification models.Verification

		id := c.Param("id")

		if err := m.Db.WithContext(c.Request().Context()).Model(&models.Verification{}).Find(verification, "id = ?", id).Error; err != nil {
			return err
		}

		return c.JSON(http.StatusOK, verification)
	}
}
