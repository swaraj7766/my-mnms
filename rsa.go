package mnms

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/qeof/q"
)

// hardcode rsa private key
var mnmsOwnPrivateKeyPEM = `
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

/*
  This file consis of bunch of functions related to crypto stuff.
*/

// GenerateOwnPublickey generate own public key
func GenerateOwnPublickey() ([]byte, error) {
	pri, err := GetPrivateKeyFromPEM([]byte(mnmsOwnPrivateKeyPEM))
	if err != nil {
		q.Q(err)
		return nil, err
	}
	pub, err := GenerateRSAPublickey(pri)
	if err != nil {
		q.Q(err)
		return nil, err
	}
	return pub, nil
}

func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

// CheckPublicKeyPath check publickey file path, if not exist, create it
func CheckPublicKeyPath() (string, error) {
	mnmsDir, err := CheckMNMSFolder()
	if err != nil {
		return "", err
	}

	return path.Join(mnmsDir, "user_public_key.pub"), nil
}

// EndcodePrivateKeyToPEM encode private key to bytes
func EndcodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	b := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return pem.EncodeToMemory(&b)
}

// GetPrivateKeyFromPEM get private key from bytes
func GetPrivateKeyFromPEM(privateKeyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// GenerateRSAPublickey generate public key from private key
func GenerateRSAPublickey(privateKey *rsa.PrivateKey) ([]byte, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		q.Q(err)
		return nil, err
	}
	b := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	return pem.EncodeToMemory(&b), nil
}

// SavePublicKeyToFile save public key to file
func SavePublicKeyToFile(publicKey []byte) error {
	filename, err := CheckPublicKeyPath()
	if err != nil {
		q.Q(err)
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		q.Q(err)
		return err
	}
	defer file.Close()
	_, err = file.Write(publicKey)
	if err != nil {
		q.Q(err)
		return err
	}
	return nil
}

// EncryptWithPublicKey encrypt data with private key PEM
func EncryptWithPublicKey(data []byte, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		q.Q("public key error")
		return nil, fmt.Errorf("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		q.Q(err)
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	// Encrypt PKCS1v15
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)

	// hash := sha512.New()
	// ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, data, nil)
	// if err != nil {
	// 	q.Q(err)
	// 	return nil, err
	// }
	// return ciphertext, nil

}

// DecryptWithPrivateKeyPEM decrypt data with private key PEM
func DecryptWithPrivateKeyPEM(data []byte, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		q.Q("private key error")
		return nil, fmt.Errorf("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		q.Q(err)
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, data)
	// hash := sha512.New()
	// plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, data, nil)
	// if err != nil {
	// 	q.Q(err)
	// 	return nil, err
	// }
	// return plaintext, nil
}
