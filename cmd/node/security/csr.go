package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"github.com/eskpil/sunlight/pkg/api/adoption"
	"os"
)

var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

// FindOrCreateCSR TODO: Store private key in tpm2 with a persistent handle if the system supports it and the domainc controller requires
// it.
func FindOrCreateCSR(hints *adoption.AdoptionHints) (*x509.CertificateRequest, *rsa.PrivateKey, error) {
	if _, err := os.Stat(os.ExpandEnv("$SUNLIGHT_VAR_DIR/csr/request.pem")); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, nil, err
		}

		privkey, _ := rsa.GenerateKey(rand.Reader, 4096)

		subj := pkix.Name{
			CommonName:         hints.GetCommonName(),
			Country:            []string{hints.GetCountry()},
			Province:           []string{hints.GetProvince()},
			Locality:           []string{hints.GetLocality()},
			Organization:       []string{hints.GetOrganization()},
			OrganizationalUnit: []string{hints.GetOrganizationalUnit()},
			ExtraNames: []pkix.AttributeTypeAndValue{
				{
					Type: oidEmailAddress,
					Value: asn1.RawValue{
						Tag:   asn1.TagIA5String,
						Bytes: []byte(hints.GetEmail()),
					},
				},
			},
		}

		template := x509.CertificateRequest{
			Subject:            subj,
			SignatureAlgorithm: x509.SHA256WithRSA,
		}

		csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, privkey)
		csrPemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})

		if err := os.WriteFile(os.ExpandEnv("$SUNLIGHT_VAR_DIR/csr/request.pem"), csrPemBytes, 0700); err != nil {
			return nil, nil, err
		}

		privkeyPemBytes := pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(privkey),
			},
		)
		if err := os.WriteFile(os.ExpandEnv("$SUNLIGHT_VAR_DIR/keys/request.pem"), privkeyPemBytes, 0700); err != nil {
			return nil, nil, err
		}

		return &template, privkey, nil
	}

	csrPemBytes, err := os.ReadFile(os.ExpandEnv("$SUNLIGHT_VAR_DIR/csr/request.pem"))
	if err != nil {
		return nil, nil, err
	}

	privkeyPemBytes, err := os.ReadFile(os.ExpandEnv("$SUNLIGHT_VAR_DIR/keys/request.pem"))
	if err != nil {
		return nil, nil, err
	}

	privkeyBlock, _ := pem.Decode(privkeyPemBytes)
	csrBlock, _ := pem.Decode(csrPemBytes)

	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	privkey, err := x509.ParsePKCS1PrivateKey(privkeyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return csr, privkey, nil
}
