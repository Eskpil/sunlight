package users

import (
	"errors"
	"github.com/eskpil/sunlight/internal/ccontroller/mycontext"
	"github.com/eskpil/sunlight/internal/ccontroller/utils"
	"github.com/eskpil/sunlight/pkg/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Create(m *mycontext.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		var body CreateInput
		if err := c.Bind(&body); err != nil {
			return err
		}

		if err := body.Validate(); err != nil {
			return err
		}

		user := new(models.User)

		user.Id = uuid.New().String()

		user.Username = body.Username

		user.FirstName = body.FirstName
		user.LastName = body.LastName

		user.Roles = make([]models.Role, 0)
		user.Keys = make([]models.Key, 0)

		user.Email = body.Email

		if err := user.EncryptPassword(body.Password); err != nil {
			return err
		}

		if err := user.Save(c.Request().Context(), m.Db); err != nil {
			if errors.Is(utils.TranslateErrors(m.Db.Dialector.Name(), err), utils.ErrUniqueKeyViolation) {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": "user with that specific username and email combination already exists",
				})
			}
		}

		return c.JSON(http.StatusOK, user)
	}
}
