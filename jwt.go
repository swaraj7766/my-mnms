package mnms

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"github.com/cloudflare/gokey"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/qeof/q"
	"github.com/xdg-go/pbkdf2"
)

func GenerateEncryptedKeySeed(pass string) ([]byte, error) {
	seed := make([]byte, 256)

	for i := range seed {
		seed[i] = 0xA5
	}

	masterkey := pbkdf2.Key([]byte(pass), []byte(seed[:12]), 4096, 32, sha256.New)

	aes, err := aes.NewCipher(masterkey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}

	pt := seed[12 : len(seed)-16]

	// encrypt in place
	gcm.Seal(pt[:0], seed[:12], pt, nil)

	return seed, nil
}

// GenPassword generates password for user
func GenPassword(hostname, account string) (string, error) {

	seedBytes, err := GenerateEncryptedKeySeed(hostname)
	if err != nil {
		return "", err
	}

	password, err := gokey.GetPass(hostname, account, seedBytes,
		&gokey.PasswordSpec{
			Length:         16,
			Upper:          3,
			Lower:          3,
			Digits:         2,
			Special:        0,
			AllowedSpecial: ""})
	if err != nil {
		return "", err
	}
	return password, nil
}

// GetToken get a internal token that can be used by CLI and node
func GetToken(name string) (string, error) {
	var token string
	var err error

	if len(token) == 0 {
		// generate special token for CLI and node
		_, token, err = jwtTokenAuth.Encode(map[string]any{
			"user":      name,
			"timestamp": time.Now().Format(time.RFC3339),
		})
		if err != nil {
			return "", err
		}
	}
	return token, nil
}

var jwtSecret = []byte("mjnwmtssecret")

// jwtTokenAuth is a global variable for JWT authentication
var jwtTokenAuth = jwtauth.New("HS256", jwtSecret, nil)

// temprary token
var tempararyUrlToken = jwtauth.New("HS256", []byte("mnmstemparayurl"), nil)

func generateJWT(user, password string) (string, error) {

	// Validate password
	if !validUserPassword(user, password) {
		return "", fmt.Errorf("wrong password or user not existed")
	}

	// token will expire in 30 days
	_, token, err := jwtTokenAuth.Encode(map[string]any{
		"user":      user,
		"timestamp": time.Now().Format(time.RFC3339),
		"exp":       time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	if err != nil {
		return "", err
	}

	return token, nil
}

func parseJWT(tokenString string) (map[string]any, error) {
	t, err := jwtTokenAuth.Decode(tokenString)
	if err != nil {
		return nil, err
	}
	claims, err := t.AsMap(context.Background())
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// JWTVerifyToken verifies the token
func JWTVerifyToken(ja *jwtauth.JWTAuth, tokenString string) (jwt.Token, error) {
	// Decode & verify the token
	token, err := ja.Decode(tokenString)
	if err != nil {
		return token, err
	}

	if token == nil {
		return nil, fmt.Errorf("no token")
	}

	if err := jwt.Validate(token); err != nil {
		q.Q("jwt.Validate fail", err)
		return token, err
	}

	// Valid!
	return token, nil
}

// JWTAuthenticatorRole is a authentication middleware to enforce access from the
// Verifier middleware request context values. The JWTAuthenticatorRole sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through. It's just fine
// until you decide to write something similar and customize your client response.
func JWTAuthenticatorRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// check is admin
		emailRaw, ok := token.Get("user")
		if !ok {
			q.Q("no user", token)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		userString, ok := emailRaw.(string)
		if !ok {
			q.Q("user is not string")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		u, err := GetUserConfig(userString)
		// q.Q(u)

		if err != nil {
			q.Q("GetUserConfig fail", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		switch role {
		case MNMSAdminRole:
			if u.Role != MNMSAdminRole {
				q.Q("user is not a admin")
				http.Error(w, fmt.Sprintf("user %s is not admin", userString), http.StatusUnauthorized)
				return
			}
		case MNMSSuperUserRole:
			if u.Role != MNMSSuperUserRole && u.Role != MNMSAdminRole {
				q.Q("user is not a super user")
				http.Error(w, fmt.Sprintf("user %s is not admin", userString), http.StatusUnauthorized)
				return
			}
		case MNMSUserRole:
			if u.Role != MNMSUserRole && u.Role != MNMSSuperUserRole && u.Role != MNMSAdminRole {
				q.Q("user is not a user")
				http.Error(w, fmt.Sprintf("user %s is not admin", userString), http.StatusUnauthorized)
				return
			}
		default:
			q.Q("unknown role")
			http.Error(w, "unknown role", http.StatusUnauthorized)
			return

		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

// SendGMail sends email using gmail
func SendGMail(to, subject, body string) error {

	from := "mistest@atop.com.tw"
	pass := "atop0130"
	msg := "Subject: " + subject + "\r\n" + "\r\n" + body + "\r\n"
	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}
