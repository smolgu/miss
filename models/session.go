package models

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	pkgErrors "github.com/pkg/errors"

	"github.com/smolgu/miss/pkg/errors"
	"github.com/smolgu/miss/pkg/setting"
)

// Claims used for jwt
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

type sessionModel int

// Sessions api for session flow
var Sessions sessionModel

func (sessionModel) New(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(),
		},
		UserID: userID,
	})

	return token.SignedString(setting.App.Secret)
}

func (sessionModel) Check(sessionID string) (int64, error) {
	token, err := jwt.Parse(sessionID, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return setting.App.Secret, nil
	})
	if err != nil {
		return 0, errors.New(fmt.Errorf("Авторизация не действительна"), pkgErrors.Wrap(err, "parse jwt token"))
	}

	if claims, ok := token.Claims.(Claims); ok {
		if token.Valid {
			return claims.UserID, nil
		}
		return 0, errors.New(fmt.Errorf("Авторизация не действительна"), fmt.Errorf("jwt token invalid"))
	}
	return 0, errors.New(fmt.Errorf("Авторизация не действительна"), fmt.Errorf("cannot cast token.Claim to session.Claims"))
}

func (sessionModel) Delete(sessionID string) error {
	return fmt.Errorf("unimplemented")
}
