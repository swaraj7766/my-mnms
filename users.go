package mnms

import (
	"errors"
	"fmt"
	"time"

	"github.com/qeof/q"
)

// validUserPassword returns true if the password is valid for the user.
func validUserPassword(user, password string) bool {
	// get mnms config
	config, err := GetMNMSConfig()
	if err != nil {
		q.Q(err)
		return false
	}
	// find user
	for _, u := range config.Users {
		if u.Name == user {

			return u.Password == password
		}
	}
	return false
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
	// get mnms config
	c, err := GetMNMSConfig()
	if err != nil {
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

// AddUserConfig add user to mnms config
func AddUserConfig(user, role, password, email string) error {
	// get mnms config
	c, err := GetMNMSConfig()
	if err != nil {
		q.Q(err)
		return err
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

// UpdateUserConfig add user to mnms config
func UpdateUserConfig(user, role, password, email string) error {
	// get mnms config
	c, err := GetMNMSConfig()
	if err != nil {
		q.Q(err)
		return err
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
