package main

import (
	"github.com/eskpil/sunlight/cmd/dcontroller/essentials"
	"github.com/eskpil/sunlight/internal/dcontroller/adoption"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	os.Setenv("SUNLIGHT_PKI_DIR", "./dcontroller1")

	essentials.Load()

	adoptionServer, err := adoption.NewServer()
	if err != nil {
		log.Fatalf("could not initialize adoption server")
	}

	log.Infof("starting adoption server")
	if err := adoptionServer.Start(); err != nil {
		log.Fatalf("could not start adoption server")
	}
}
