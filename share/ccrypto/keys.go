package ccrypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// GenerateKey for use as an SSH private key
func GenerateKey(seed string) ([]byte, error) {
	var err error
	var priv *ecdsa.PrivateKey

	r := rand.Reader
	if seed != "" {
		r = NewDetermRand([]byte(seed))
	}

	if seed == "" {
		priv, err = ecdsa.GenerateKey(elliptic.P256(), r)
	} else {
		priv, err = GenerateKeyGo119(elliptic.P256(), r)
	}

	if err != nil {
		return nil, err
	}
	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal ECDSA private key: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b}), nil
}

// GenerateKeyFile generates an SSH private key file
func GenerateKeyFile(keyFilePath, seed string) error {
	keyBytes, err := GenerateKey(seed)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(keyFilePath, keyBytes, 0600)
}

// FingerprintKey calculates the SHA256 hash of an SSH public key
func FingerprintKey(k ssh.PublicKey) string {
	bytes := sha256.Sum256(k.Marshal())
	return base64.StdEncoding.EncodeToString(bytes[:])
}
