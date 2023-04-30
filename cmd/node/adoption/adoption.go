package adoption

import (
	"context"
	"crypto/tls"
	"github.com/eskpil/sunlight/pkg/api/adoption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
)

func discoverMachineId() (string, error) {
	if _, err := os.Stat("/etc/machineid"); err != nil {
		panic("/etc/machineid not found")
		return "", err
	}

	machineId, err := os.ReadFile("/etc/machineid")
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

	conn, err := grpc.Dial("local.sunlight.:2100", grpc.WithTransportCredentials(creds))
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

	return nil
}
