package main

import (
	"github.com/eskpil/sunlight/cmd/node/essentials"
	"github.com/eskpil/sunlight/cmd/node/security"
	"github.com/eskpil/sunlight/pkg/api/adoption"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	os.Setenv("SUNLIGHT_VAR_DIR", "./node1")

	if err := essentials.Load(); err != nil {
		log.Fatalf("failed to load essentials: %v", err)
	}

	hints := new(adoption.AdoptionHints)

	hints.CommonName = "test.local.sunlight"
	hints.Country = "Norway"
	hints.Locality = "NO"
	hints.Province = "Trondelag"
	hints.Organization = "Eskpil"
	hints.OrganizationalUnit = "machines"
	hints.Email = "admin@local.sunlight"

	csr, privkey, err := security.FindOrCreateCSR(hints)
	if err != nil {
		log.Fatalf("Could not find our create a csr: %v", err)
	}

	_ = csr
	_ = privkey

	log.Infof("Hello, Node!")
}
