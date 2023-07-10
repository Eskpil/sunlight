package ca

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/eskpil/sunlight/internal/ccontroller/mycontext"
	"github.com/eskpil/sunlight/pkg/models"
	"math/big"
	"time"
)

type Entry struct {
	Id     string
	Parent string

	Children []*Entry

	Certificate *x509.Certificate
	PrivateKey  *rsa.PrivateKey

	Certificates []*x509.Certificate
}

type Chain struct {
	Id       string
	Children []Entry

	Context *mycontext.Context
}

func (e *Entry) Find(ca string) (*Entry, error) {
	if e.Certificate.IsCA && e.Certificate.Subject.CommonName == ca {
		return e, nil
	}

	for _, child := range e.Children {
		entry, _ := child.Find(ca)
		if entry != nil {
			return entry, nil
		}
	}

	return nil, nil
}

func (c *Chain) Find(ca string) (*Entry, error) {
	for _, child := range c.Children {
		entry, _ := child.Find(ca)
		if entry != nil {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("ca: %s not found", ca)
}

func (c *Chain) Sign(ca string, csr *x509.CertificateRequest) (*x509.Certificate, error) {
	entry, err := c.Find(ca)
	if err != nil {
		return nil, err
	}

	cert, err := entry.Sign(csr)
	if err != nil {
		return nil, err
	}

	if err := c.Save(); err != nil {
		return nil, err
	}

	return cert, nil
}

func (e *Entry) Sign(csr *x509.CertificateRequest) (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		Signature:          csr.Signature,
		SignatureAlgorithm: csr.SignatureAlgorithm,

		PublicKeyAlgorithm: csr.PublicKeyAlgorithm,
		PublicKey:          csr.PublicKey,

		SerialNumber: serialNumber,
		Issuer:       e.Certificate.Subject,
		Subject:      csr.Subject,

		// TODO: Get this from the domain controllers configuration living in the core controller
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365),

		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	// Create certificate from template and root certificate, signed by the RootCA's private key.
	certData, err := x509.CreateCertificate(rand.Reader, &template, e.Certificate, template.PublicKey, e.PrivateKey)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		return nil, err
	}

	e.Certificates = append(e.Certificates, cert)

	return cert, nil
}

func (e *Entry) AsDbEntry() (*models.CAEntry, error) {
	dbEntry := new(models.CAEntry)

	dbEntry.Id = e.Id
	dbEntry.Parent = e.Parent

	for _, rawCert := range e.Certificates {
		certPem := &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: rawCert.Raw,
		}

		pemBytes := pem.EncodeToMemory(certPem)

		dbEntry.Certificates = append(dbEntry.Certificates, string(pemBytes))
	}

	certPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: e.Certificate.Raw,
	}

	dbEntry.Certificate = string(pem.EncodeToMemory(certPem))

	privKeyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(e.PrivateKey),
	}
	dbEntry.PrivateKey = string(pem.EncodeToMemory(privKeyPem))

	for _, rawChild := range e.Children {
		child, err := rawChild.AsDbEntry()
		if err != nil {
			return nil, err
		}

		dbEntry.Children = append(dbEntry.Children, *child)
	}

	return dbEntry, nil
}

func EntryFromDbEntry(dbEntry *models.CAEntry) (*Entry, error) {
	entry := new(Entry)

	entry.Id = dbEntry.Id
	entry.Parent = dbEntry.Parent

	privkeyBlock, _ := pem.Decode([]byte(dbEntry.PrivateKey))
	certBlock, _ := pem.Decode([]byte(dbEntry.Certificate))

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, err
	}

	privkey, err := x509.ParsePKCS1PrivateKey(privkeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	entry.Certificate = cert
	entry.PrivateKey = privkey

	for _, dbChild := range dbEntry.Children {
		child, err := EntryFromDbEntry(&dbChild)
		if err != nil {
			return nil, err
		}

		entry.Children = append(entry.Children, child)
	}

	for _, rawCert := range dbEntry.Certificates {
		certBlock, _ := pem.Decode([]byte(rawCert))
		cert, err := x509.ParseCertificate(certBlock.Bytes)
		if err != nil {
			return nil, err
		}

		entry.Certificates = append(entry.Certificates, cert)
	}

	return entry, nil
}

// Load loads the certificate authority tree from the database.
func (c *Chain) Load() error {
	rawChain, err := models.GetCAChain(context.Background(), c.Context.Db)
	if err != nil {
		return err
	}

	for _, rawChild := range rawChain.Children {
		child, err := EntryFromDbEntry(&rawChild)
		if err != nil {
			return err
		}

		c.Children = append(c.Children, *child)
	}

	c.Id = rawChain.Id

	return nil
}

func (c *Chain) Save() error {
	var children []models.CAEntry

	for _, rawChild := range c.Children {
		child, err := rawChild.AsDbEntry()
		if err != nil {
			return err
		}

		children = append(children, *child)
	}

	dbChain := models.CAChain{
		Id:       c.Id,
		Children: children,
	}

	if err := c.Context.Db.Save(&dbChain).Error; err != nil {
		return err
	}

	return nil
}

func New(m *mycontext.Context) (*Chain, error) {
	chain := new(Chain)

	chain.Context = m

	if err := chain.Load(); err != nil {
		return nil, err
	}

	return chain, nil
}
