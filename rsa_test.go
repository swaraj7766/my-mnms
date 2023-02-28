package mnms

import "testing"

// TestPKI test rsa key pair generation and save to file
func TestPKI(t *testing.T) {
	prikey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		t.Fatal("generate rsa key pair fail", err)
	}
	// encode private key to pem
	pemPriKey := EndcodePrivateKeyToPEM(prikey)
	// generate public key from private key
	pubKey, err := GenerateRSAPublickey(prikey)
	if err != nil {
		t.Fatal("generate rsa public key fail", err)
	}

	data := []byte("hello world")
	// encrypt data with public key
	encryptedData, err := EncryptWithPublicKey(data, pubKey)
	if err != nil {
		t.Fatal("encrypt data fail", err)
	}
	t.Log(string(encryptedData))
	// decrypt data with private key
	decryptedData, err := DecryptWithPrivateKeyPEM(encryptedData, pemPriKey)
	if err != nil {
		t.Fatal("decrypt data fail", err)
	}
	if string(decryptedData) != string(data) {
		t.Fatal("expect", string(data), "but got", string(decryptedData))
	}
}

// TestMnmsOwnPrivateKey test mnms own private key
func TestMnmsOwnPrivateKey(t *testing.T) {
	// get private key
	ownPri, err := GetPrivateKeyFromPEM([]byte(mnmsOwnPrivateKeyPEM))
	if err != nil {
		t.Fatal("get private key fail", err)
	}
	// generate public key from private key
	pubKey, err := GenerateRSAPublickey(ownPri)
	if err != nil {
		t.Fatal("generate rsa public key fail", err)
	}

	data := []byte("hello world")
	// encrypt data with public key
	encryptedData, err := EncryptWithPublicKey(data, pubKey)
	if err != nil {
		t.Fatal("encrypt data fail", err)
	}

	// decrypt data with private key
	decryptedData, err := DecryptWithPrivateKeyPEM(encryptedData, []byte(mnmsOwnPrivateKeyPEM))
	if err != nil {
		t.Fatal("decrypt data fail", err)
	}
	if string(decryptedData) != string(data) {
		t.Fatal("expect", string(data), "but got", string(decryptedData))
	}
}
