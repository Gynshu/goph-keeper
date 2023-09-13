package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"io"
	"math/big"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

var instance *config

type config struct {
	// Server is the server configuration
	MongoURI       string `json:"mongo_uri"`
	HttpServerPort string `json:"http_server_port"`
	CertFilePath   string `json:"cert_file_path"`
	KeyFilePath    string `json:"key_file_path"`
}

// NewConfig creates a new configuration struct
func init() {
	// Initialize the config struct
	instance = &config{}

	// Read for flags
	flag.StringVar(&instance.MongoURI, "mongo_uri",
		"mongodb://admin:password@localhost:27017",
		"MongoDB URI default: mongodb://admin:password@localhost:27017")
	flag.StringVar(&instance.HttpServerPort, "port", "8080", "HTTP server port default: 8080")
	flag.StringVar(&instance.CertFilePath, "cert", "", "Certificate file path default: empty")
	flag.StringVar(&instance.KeyFilePath, "key", "", "Key file path default: empty")

	// Parse the flags and ignore the rest
	flag.CommandLine.SetOutput(io.Discard)

	// If no certificate is provided, generate a self-signed one
	if instance.CertFilePath == "" || instance.KeyFilePath == "" {
		log.Info().Msg("Generating self-signed certificate")
		osTempDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get user home dir")
		}
		err = os.MkdirAll(osTempDir+string(os.PathSeparator)+".goph-keeper", 0700)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create .goph-keeper dir")
		}
		instance.CertFilePath = osTempDir + string(os.PathSeparator) + ".goph-keeper/cert.pem"
		instance.KeyFilePath = osTempDir + string(os.PathSeparator) + ".goph-keeper/key.pem"

		err = generateCertificate(instance.CertFilePath, instance.KeyFilePath)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to generate certificate")
		}
		log.Info().Msgf("Generated certificate path : %s", instance.CertFilePath)
		log.Info().Msgf("Generated key path : %s", instance.KeyFilePath)
		return
	}
	log.Info().Msgf("Using user provided certificate: %s", instance.CertFilePath)
	log.Info().Msgf("Using user provided key: %s", instance.KeyFilePath)
}

// GetConfig returns the configuration initialized by newConfig
func GetConfig() *config {
	return instance
}

func generateCertificate(certFile, keyFile string) error {
	// Generate a new private key
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// Create a new certificate template
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Add a subject alternative name to the certificate template
	template.DNSNames = append(template.DNSNames, "localhost")

	// Create a new self-signed certificate using the private key and certificate template
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return err
	}

	// Write the private key to a file
	keyFileOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer keyFileOut.Close()
	err = pem.Encode(keyFileOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	if err != nil {
		return err
	}

	// Write the certificate to a file
	certFileOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer certFileOut.Close()
	err = pem.Encode(certFileOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return err
	}

	instance.CertFilePath = certFile
	instance.KeyFilePath = keyFile
	return nil
}
