package adoption

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"github.com/eskpil/sunlight/cmd/node/security"
	"github.com/eskpil/sunlight/pkg/api/adoption"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
)

func discoverMachineId() (string, error) {
	if _, err := os.Stat("/etc/machine-id"); err != nil {
		panic("/etc/machine-id not found")
		return "", err
	}

	machineId, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		return "", err
	}

	return string(machineId), nil
}

type Adopter struct {
	client adoption.AdoptionClient
}

func New() (*Adopter, error) {
	adopter := new(Adopter)

	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true,
	})

	conn, err := grpc.Dial("local.sunlight.:2001", grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	adopter.client = adoption.NewAdoptionClient(conn)

	return adopter, nil
}

func (a *Adopter) AttemptAdoption(ctx context.Context) error {
	machineId, err := discoverMachineId()
	if err != nil {
		return err
	}

	needsAdoptionRequest := &adoption.NeedsAdoptionRequest{
		MachineId: machineId,
	}

	res, err := a.client.NeedsAdoption(ctx, needsAdoptionRequest)
	if err != nil {
		return err
	}

	// we don't need to perform adoption
	if !res.Verdict {
		return nil
	}

	log.Infof("%v", res)

	csr, privkey, err := security.FindOrCreateCSR(res.GetHints())
	if err != nil {
		return err
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, csr, privkey)

	adoptRequest := &adoption.AdoptRequest{
		Machineid:             machineId,
		Csr:                   string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})),
		Metadata:              nil,
		FulfilledRequirements: []string{"tpm2"},
	}

	adoptRes, err := a.client.Adopt(ctx, adoptRequest)
	if err != nil {
		return err
	}

	log.Infof("%v", adoptRes)

	return nil
}
