package mnms

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/qeof/q"
)

func SetPrivateKey(key string) {
	mnmsOwnPrivateKeyPEM = key
}

func GetDefaultPrivateKey() string {
	return mnmsOwnPrivateKeyPEM
}

// hardcode rsa private key
var mnmsOwnPrivateKeyPEM = `
-----BEGIN PRIVATE KEY-----
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
74TouNiy6+gL9VZC4Q9NBYMbH9Yp4wDQRQ==
-----END PRIVATE KEY-----
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

// EndcodePrivateKeyToPEM encode private key to bytes
func EndcodePrivateKeyToPEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	// privatekey to pkcs8 pem

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	b := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return pem.EncodeToMemory(&b), nil
}

// GetPrivateKeyFromPEM get private key from bytes
func GetPrivateKeyFromPEM(privateKeyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("failed to parse RSA private key")
	}
	return rsaPrivateKey, nil
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
	hash := sha256.New()
	setp := pub.Size() - 2*hash.Size() - 2
	mlen := len(data)
	var ciphertext []byte
	for i := 0; i < mlen; i += setp {
		if i+setp > mlen {
			chunk, err := rsa.EncryptOAEP(hash, rand.Reader, pub, data[i:], nil)
			if err != nil {
				q.Q(err)
				return nil, err
			}
			ciphertext = append(ciphertext, chunk...)
		} else {
			chunk, err := rsa.EncryptOAEP(hash, rand.Reader, pub, data[i:i+setp], nil)
			if err != nil {
				q.Q(err)
				return nil, err
			}
			ciphertext = append(ciphertext, chunk...)
		}
	}
	base64Text := base64.StdEncoding.EncodeToString(ciphertext)
	return []byte(base64Text), err

}

// DecryptWithOwnPrivateKey decrypt data with own private key
func DecryptWithOwnPrivateKey(data []byte) ([]byte, error) {
	// get private key
	return DecryptWithPrivateKeyPEM(data, []byte(mnmsOwnPrivateKeyPEM))
}

// DecryptWithPrivateKeyPEM decrypt data with private key PEM
func DecryptWithPrivateKeyPEM(txtdata []byte, privateKey []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(txtdata))
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(privateKey)
	if block == nil {
		q.Q("private key error")
		return nil, fmt.Errorf("private key error")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		q.Q(err)
		return nil, err
	}

	rsaPrivateKey, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("failed to parse RSA private key")
	}

	mlen := len(data)
	hash := sha256.New()
	setp := rsaPrivateKey.PublicKey.Size()
	var plaintext []byte
	for i := 0; i < mlen; i += setp {
		if i+setp > mlen {
			chunk, err := rsa.DecryptOAEP(hash, rand.Reader, rsaPrivateKey, data[i:], nil)
			if err != nil {
				q.Q(err)
				return nil, err
			}
			plaintext = append(plaintext, chunk...)
		} else {
			chunk, err := rsa.DecryptOAEP(hash, rand.Reader, rsaPrivateKey, data[i:i+setp], nil)
			if err != nil {
				q.Q(err)
				return nil, err
			}
			plaintext = append(plaintext, chunk...)
		}
	}

	return plaintext, err
}
