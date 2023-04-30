package main

import (
	machines2 "github.com/eskpil/sunlight/internal/api/machines"
	"github.com/eskpil/sunlight/internal/api/mycontext"
	users2 "github.com/eskpil/sunlight/internal/api/users"
	"github.com/eskpil/sunlight/internal/api/verifications"
	"github.com/eskpil/sunlight/pkg/models"
	"github.com/labstack/echo/v4"
)

func main() {
	m, err := mycontext.New()
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	if err := m.Db.AutoMigrate(&models.Group{}); err != nil {
		panic(err)
	}

	if err := m.Db.AutoMigrate(&models.Key{}); err != nil {
		panic(err)
	}

	if err := m.Db.AutoMigrate(&models.Role{}); err != nil {
		panic(err)
	}

	if err := m.Db.AutoMigrate(&models.User{}); err != nil {
		panic(err)
	}

	if err := m.Db.AutoMigrate(&models.Machine{}); err != nil {
		panic(err)
	}

	if err := m.Db.AutoMigrate(&models.Verification{}); err != nil {
		panic(err)
	}

	s := echo.New()

	s.POST("/v1/users/", users2.Create(m))
	s.GET("/v1/users/", users2.GetAll(m))
	s.GET("/v1/users/:id/", users2.GetById(m))

	s.POST("/v1/machines/", machines2.Create(m))
	s.GET("/v1/machines/", machines2.GetAll(m))
	s.GET("/v1/machines/:id/", machines2.GetById(m))

	s.GET("/v1/verifications/", verifications.GetAll(m))
	s.GET("/v1/verifications/:id/", verifications.GetById(m))

	s.Logger.Info(s.Start("0.0.0.0:9000"))
}
