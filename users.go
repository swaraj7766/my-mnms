package mnms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/pquerna/otp/totp"

	"github.com/qeof/q"
)

// validUserPassword returns true if the password is valid for the user.
func validUserPassword(user, password string) error {
	// get mnms config
	config, err := GetMNMSConfig()
	if err != nil {
		q.Q(err)
		return err
	}
	// find user
	for _, u := range config.Users {
		if u.Name == user {

			if u.Password != password {
				return errors.New("password not match")
			}
			return nil
		}
	}
	return fmt.Errorf("user %s not found", user)
}

func GenerateRetrievePasswordToken(user, pass string) (string, error) {

	if !UserExist(user) {
		q.Q("user not exist", user)
		return "", errors.New("user not exist")
	}

	// generate token
	_, token, err := tempararyUrlToken.Encode(map[string]any{
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
		"user": user,
		"pass": pass,
	})
	if err != nil {
		q.Q(err)
		return "", err
	}

	return token, nil
}

// UserExit check user exist in UsersList
func UserExist(user string) bool {
	// get mnms config
	c, err := GetMNMSConfig()
	if err != nil {
		q.Q(err)
		return false
	}
	// check user exist
	for _, u := range c.Users {
		if u.Name == user {
			return true
		}
	}
	return false
}

// ValidateUserPassword validate user password
func ValidateUserPassword(user, password string) bool {
	// get mnms config
	c, err := GetMNMSConfig()
	if err != nil {
		q.Q(err)
		return false
	}
	// check user exist
	for _, u := range c.Users {
		if u.Name == user {
			// check password

			if u.Password == password {
				return true
			}
		}
	}
	return false
}

// GetUserConfig
func GetUserConfig(user string) (*UserConfig, error) {
	// admin doesn't need to depend on config.json

	c, err := GetMNMSConfig()
	if err != nil {
		// is user is admin and config.json not exist, return admin because admin do not need to depend on config.json
		if user == "admin" {

			return &UserConfig{
				Name:     "admin",
				Role:     MNMSAdminRole,
				Password: AdminDefaultPassword,
			}, nil
		}
		q.Q(err)
		return nil, err
	}
	// check user exist
	for _, u := range c.Users {
		if u.Name == user {
			return &u, nil
		}
	}
	q.Q("user not exist", user)
	return nil, errors.New("user not exist")
}
func checkUsersPassword(pw string) error {
	// check password with regex: /^(?=.\d)(?=.[A-Z])(?=.[a-z])(?=.[a-zA-Z!#@$%&? "])[a-zA-Z0-9!#$@%&?]{8,20}$/
	regex := regexp2.MustCompile(`(?=.*[A-Z])(?=.*[a-z])(?=.*\d)(?=.*[@$!%*#?&])[A-Za-z\d@$!%*#?&]{8,20}$`, 0)
	if isMatch, _ := regex.MatchString(pw); isMatch {
		return nil
	}
	return errors.New("password must have 8-20 characters, at least one uppercase one lowercase one digit one special character")
}

// AddUserConfig add user to mnms config
func AddUserConfig(user, role, password, email string) error {
	// get mnms config
	c, err := GetMNMSConfig()
	if err != nil {
		q.Q(err)
		return err
	}

	err = checkUsersPassword(password)
	if err != nil {
		q.Q(err)
		return err
	}

	// check email exist
	for _, u := range c.Users {
		if u.Email == email {
			q.Q("email exist", email)
			return fmt.Errorf("email %s exist", email)
		}
	}

	// check user exist
	for _, u := range c.Users {
		if u.Name == user {
			q.Q("user exist", user)
			return fmt.Errorf("user %s exist", user)
		}
	}
	// add user
	c.Users = append(c.Users, UserConfig{
		Name:     user,
		Role:     role,
		Email:    email,
		Password: password,
	})
	err = WriteMNMSConfig(c)
	if err != nil {
		q.Q(err)
		return err
	}
	return nil
}

// DeleteUserConfig delete user from mnms config
func DeleteUserConfig(user string) error {
	// get mnms config
	c, err := GetMNMSConfig()
	if err != nil {
		q.Q(err)
		return err
	}
	// check user exist
	for i, u := range c.Users {
		if u.Name == user {
			c.Users = append(c.Users[:i], c.Users[i+1:]...)
			err = WriteMNMSConfig(c)
			if err != nil {
				q.Q(err)
				return err
			}
			return nil
		}
	}
	q.Q("user not exist", user)
	return fmt.Errorf("user %s not exist", user)
}

// MergeUserConfig merge user config
func MergeUserConfig(user UserConfig) error {
	// get mnms config
	c, err := GetMNMSConfig()
	if err != nil {
		q.Q(err)
		return err
	}
	// check user exist
	for i, u := range c.Users {
		if u.Name == user.Name {
			// merge user if not empty

			if user.Email != "" {
				c.Users[i].Email = user.Email
			}
			if user.Password != "" {
				c.Users[i].Password = user.Password
			}
			if user.Role != "" {
				c.Users[i].Role = user.Role
			}

			c.Users[i].Enable2FA = user.Enable2FA
			if user.Secret != "" {
				c.Users[i].Secret = user.Secret
			}
			err = WriteMNMSConfig(c)
			if err != nil {
				q.Q(err)
				return err
			}
			return nil
		}
	}
	q.Q("user not exist", user.Name)
	return fmt.Errorf("user %s not exist", user.Name)
}

// UpdateUserConfig add user to mnms config
func UpdateUserConfig(user, role, password, email string) error {
	// get mnms config
	c, err := GetMNMSConfig()
	if err != nil {
		q.Q(err)
		return err
	}

	err = checkUsersPassword(password)
	if err != nil {
		q.Q(err)
		return err
	}
	// check email exist
	for _, u := range c.Users {
		if u.Email == email {
			if u.Name == user {
				continue
			}
			q.Q("email exist", email)
			return fmt.Errorf("email %s exist", email)
		}
	}

	// check user exist
	for i, u := range c.Users {
		if u.Name == user {
			c.Users[i].Role = role
			c.Users[i].Password = password
			c.Users[i].Email = email
			err = WriteMNMSConfig(c)
			if err != nil {
				q.Q(err)
				return err
			}
			return nil
		}
	}
	q.Q("user not exist", user)
	return fmt.Errorf("user %s not exist", user)
}

// the mnms config mutex for writing file avoid race condition
var mnmsconfigMutex = sync.Mutex{}

var MNMSAdminRole = "admin"
var MNMSSuperUserRole = "superuser"
var MNMSUserRole = "user"

// UserConfig is the configuration for a user.
type UserConfig struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Password  string `json:"password"`
	Enable2FA bool   `json:"enable2FA"`
	Secret    string `json:"secret"`
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

	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return workingDir, nil
}

// checkMNMSConfigPath check users file path, if not exist, create it
func checkMNMSConfigPath() (string, error) {
	mnmsDir, err := CheckMNMSFolder()
	if err != nil {
		return "", err
	}

	return path.Join(mnmsDir, "config.json"), nil
}

var AdminDefaultPassword = "default"

type LoginSession struct {
	User      UserConfig
	ExpiresAt time.Time
}

var loginsSssionStore = struct {
	sync.RWMutex
	m map[string]LoginSession
}{m: make(map[string]LoginSession)}

// createLoginSession create a login session
func createLoginSession(user UserConfig) string {
	sessionID := fmt.Sprintf("%s-%d", user.Name, time.Now().UnixNano())
	// create session
	session := LoginSession{
		User:      user,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	loginsSssionStore.Lock()
	defer loginsSssionStore.Unlock()
	loginsSssionStore.m[sessionID] = session
	return sessionID
}

// getLoginSession get a login session
func getLoginSession(sessionID string) (*UserConfig, error) {
	loginsSssionStore.RLock()
	defer loginsSssionStore.RUnlock()
	session, ok := loginsSssionStore.m[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not exist")
	}
	if session.ExpiresAt.Before(time.Now()) {
		delete(loginsSssionStore.m, sessionID)
		return nil, fmt.Errorf("session expired")
	}
	return &session.User, nil
}

// InitDefaultMNMSConfigIfNotExist generate default users, if already there, do nothing
func InitDefaultMNMSConfigIfNotExist() error {
	_, err := GetMNMSConfig()
	if err != nil {
		// create default config

		config := &MNMSConfig{
			Users: []UserConfig{
				{
					Name:      "admin",
					Role:      MNMSAdminRole,
					Password:  AdminDefaultPassword,
					Enable2FA: false,
				},
			},
		}
		err = WriteMNMSConfig(config)
		if err != nil {
			q.Q(err)
			return err
		}
	}
	return nil

}

const IssuerOf2FA = "Atop_MNMS"

func generate2FASecret(email string) (string, error) {
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      IssuerOf2FA,
		AccountName: email,
	})

	if err != nil {
		return "", err
	}

	return secret.Secret(), nil

}
