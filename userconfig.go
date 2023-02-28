package mnms

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/qeof/q"
)

/*
	This file consis of bunch of functions related to mnms config file.
*/

// the mnms config mutex for writing file avoid race condition
var mnmsconfigMutex = sync.Mutex{}

var MNMSAdminRole = "admin"
var MNMSSuperUserRole = "superuser"
var MNMSUserRole = "user"

// UserConfig is the configuration for a user.
type UserConfig struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

// MNMSConfig is the configuration for the MNMS.
type MNMSConfig struct {
	Users []UserConfig `json:"users"`
}

// GetMNMSConfig returns the MNMS configuration
func GetMNMSConfig() (*MNMSConfig, error) {
	configFullPath, err := checkMNMSConfigPath()
	if err != nil {
		return nil, err
	}
	var config MNMSConfig

	// read all from configFullPath
	configEncryptedBytes, err := ioutil.ReadFile(configFullPath)
	if err != nil {
		q.Q(err)
		return nil, err
	}
	// decrypt configEncryptedBytes with mnmsOwnPrivatekeyPEM
	decryptedConfig, err := DecryptWithPrivateKeyPEM(configEncryptedBytes, []byte(mnmsOwnPrivateKeyPEM))
	if err != nil {
		q.Q(err)
		return nil, err
	}

	err = json.Unmarshal(decryptedConfig, &config)

	if err != nil {
		q.Q(err)
		return nil, err
	}
	return &config, err
}

// WriteMNMSConfig writes the MNMS configuration
func WriteMNMSConfig(c *MNMSConfig) error {
	mnmsconfigMutex.Lock()
	defer mnmsconfigMutex.Unlock()
	configFullPath, err := checkMNMSConfigPath()
	if err != nil {
		q.Q(err)
		return err
	}

	configJSON, err := json.Marshal(c)
	if err != nil {
		q.Q(err)
		return err
	}
	publickey := QC.OwnPublicKeys
	if len(QC.OwnPublicKeys) <= 0 {
		publickey, err = GenerateOwnPublickey()
		if err != nil {
			q.Q(err)
			return err
		}
	}
	// encrypt configJSON with QC.OwnPublicKey
	encryptedConfig, err := EncryptWithPublicKey(configJSON, publickey)
	if err != nil {
		q.Q(err)
		return err
	}

	// write encryptedConfig to configFullPath
	err = ioutil.WriteFile(configFullPath, encryptedConfig, 0644)
	if err != nil {
		q.Q(err)
		return err
	}

	return nil
}

// cleanMNMSConfig remove users list file
func cleanMNMSConfig() error {
	fullPath, err := checkMNMSConfigPath()
	if err != nil {
		return err
	}
	return os.Remove(fullPath)
}

// CheckMNMSFolder check mnms folder, if not exist, create it
func CheckMNMSFolder() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	folder := path.Join(homeDir, ".mnms")
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, 0o700)
		if err != nil {
			return "", err
		}
	}
	return folder, nil
}

// checkMNMSConfigPath check users file path, if not exist, create it
func checkMNMSConfigPath() (string, error) {
	mnmsDir, err := CheckMNMSFolder()
	if err != nil {
		return "", err
	}

	return path.Join(mnmsDir, "config.json"), nil
}

// InitDefaultMNMSConfigIfNotExist generate default users, if already there, do nothing
func InitDefaultMNMSConfigIfNotExist() error {
	_, err := GetMNMSConfig()
	if err != nil {
		// create default config
		adminPass, err := GenPassword(QC.Name, "admin")
		if err != nil {
			q.Q(err)
			return err
		}

		config := &MNMSConfig{
			Users: []UserConfig{
				{
					Name:     "admin",
					Role:     MNMSAdminRole,
					Password: adminPass,
				},
			},
		}
		err = WriteMNMSConfig(config)
		if err != nil {
			q.Q(err)
			return err
		}
	}
	// exist, update admin password with QC.Name
	pass, err := GenPassword(QC.Name, "admin")
	if err != nil {
		q.Q(err)
		return err
	}
	c, err := GetUserConfig("admin")
	if err != nil {
		q.Q(err)
		return err
	}
	UpdateUserConfig(c.Name, c.Role, pass, c.Email)
	return nil
}
