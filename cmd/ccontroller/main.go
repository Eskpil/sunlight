package main

import (
	"fmt"
	"github.com/eskpil/sunlight/cmd/ccontroller/essentials"
	"github.com/eskpil/sunlight/internal/ca"
	caRoutes "github.com/eskpil/sunlight/internal/ccontroller/ca"
	machines2 "github.com/eskpil/sunlight/internal/ccontroller/machines"
	"github.com/eskpil/sunlight/internal/ccontroller/mycontext"
	users2 "github.com/eskpil/sunlight/internal/ccontroller/users"
	"github.com/eskpil/sunlight/internal/ccontroller/verifications"
	"github.com/eskpil/sunlight/pkg/models"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

func main() {
	os.Setenv("SUNLIGHT_PKI_DIR", "ccontroller1")
	essentials.Load()

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

	if err := m.Db.AutoMigrate(&models.Resource{}); err != nil {
		panic(err)
	}

	if err := m.Db.AutoMigrate(&models.Domain{}); err != nil {
		panic(err)
	}

	if err := m.Db.AutoMigrate(&models.CAEntry{}); err != nil {
		panic(err)
	}

	if err := m.Db.AutoMigrate(&models.CAChain{}); err != nil {
		panic(err)
	}

	chain, err := ca.New(m)
	if err != nil {
		log.Fatalf("failed to construct CA chain: %v", err)
	}

	fmt.Printf("chain: %v", chain)

	e := echo.New()

	e.POST("/v1/users/", users2.Create(m))
	e.GET("/v1/users/", users2.GetAll(m))
	e.GET("/v1/users/:id/", users2.GetById(m))

	e.POST("/v1/machines/", machines2.Create(m))
	e.GET("/v1/machines/", machines2.GetAll(m))
	e.GET("/v1/machines/:id/", machines2.GetById(m))

	e.GET("/v1/ca/", caRoutes.GetAll(m))

	e.GET("/v1/verifications/", verifications.GetAll(m))
	e.GET("/v1/verifications/:id/", verifications.GetById(m))

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func(w *sync.WaitGroup) {
		defer wg.Done()
		e.Logger.Info(e.Start("0.0.0.0:9000"))
	}(wg)

	//cs, err := coreserver.NewServer()
	//if err != nil {
	//	log.Fatalf("failed to create a new coreserver: %v", err)
	//}

	//wg.Add(1)
	//go func(w *sync.WaitGroup) {
	//	defer wg.Done()
	//	if err := cs.Start(); err != nil {
	//		log.Fatalf("failed to run coreserver: %v", err)
	//	}
	//}(wg)

	wg.Wait()
}
