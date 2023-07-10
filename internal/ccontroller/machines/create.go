package machines

import (
	"errors"
	"fmt"
	"github.com/eskpil/sunlight/internal/ccontroller/mycontext"
	"github.com/eskpil/sunlight/internal/ccontroller/utils"
	"github.com/eskpil/sunlight/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Create(m *mycontext.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		var body CreateInput
		if err := c.Bind(&body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": err.Error(),
			})
		}

		if err := body.Validate(); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": err.Error(),
			})
		}

		var machine models.Machine

		machineId := uuid.New().String()
		verificationId := uuid.New().String()

		machine.Id = machineId
		machine.Hostname = body.Hostname

		// We can assume this since the machines mere seconds ago sent this request.
		machine.State = models.MachineStateOperational
		machine.Status = models.MachineStatusPending

		machine.VerificationId = verificationId

		machine.PublicKey = body.PublicKey

		if err := m.Db.WithContext(c.Request().Context()).Create(&machine).Error; err != nil {
			if errors.Is(utils.TranslateErrors(m.Db.Dialector.Name(), err), utils.ErrUniqueKeyViolation) {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": fmt.Sprintf("machine with hostname: \"%s\" is already registered", body.Hostname),
				})
			}

			log.Errorf("unknown error: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "please contact a system administrator",
			})
		}

		verification := models.Verification{
			Id:        verificationId,
			Status:    models.VerificationStatusPending,
			MachineId: machine.Id,
		}

		if err := m.Db.WithContext(c.Request().Context()).Create(&verification).Error; err != nil {
			log.Errorf("unknown error: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "please contact a system administrator",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"verification": verification,
			"machines":     machine,
		})
	}
}
