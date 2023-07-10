package ca

import (
	"github.com/eskpil/sunlight/internal/ccontroller/mycontext"
	"github.com/eskpil/sunlight/pkg/models"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GetAll(m *mycontext.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		chain, err := models.GetCAChain(c.Request().Context(), m.Db)
		if err != nil {
			log.Errorf("could not get ca chain: %v", err)
		}

		return c.JSON(http.StatusOK, chain)
	}
}
