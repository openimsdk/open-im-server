package tls

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
)


func decryptPEM(data []byte, passphrase []byte) ([]byte, error) {
	if len(passphrase) == 0 {
		return data, nil
	}
	b, _ := pem.Decode(data)
	d, err := x509.DecryptPEMBlock(b, passphrase)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  b.Type,
		Bytes: d,
	}), nil
}

func readEncryptablePEMBlock(path string, pwd []byte) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return decryptPEM(data, pwd)
}

// NewTLSConfig setup the TLS config from general config file.
func NewTLSConfig(clientCertFile, clientKeyFile, caCertFile string, keyPwd []byte) *tls.Config {
	tlsConfig := tls.Config{}

	if clientCertFile != "" && clientKeyFile != "" {
		certPEMBlock, err := os.ReadFile(clientCertFile)
		if err != nil {
			panic(err)
		}
		keyPEMBlock, err := readEncryptablePEMBlock(clientKeyFile, keyPwd)
		if err != nil {
			panic(err)
		}
		cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			panic(err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		panic(errors.New("not a valid CA cert"))
	}
	tlsConfig.RootCAs = caCertPool

	tlsConfig.InsecureSkipVerify = config.Config.Kafka.TLS.InsecureSkipVerify

	return &tlsConfig
}
