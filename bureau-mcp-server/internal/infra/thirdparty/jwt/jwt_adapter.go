package jwt

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/golang-jwt/jwt/v5"
)

type jwtAdapter struct{}

func NewJWTAdapter() *jwtAdapter {
	return &jwtAdapter{}
}

func (ja *jwtAdapter) GenerateToken(tid, sub string) (string, *rest_err.RestErr) {
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		return "", rest_err.NewInternalServerError("token secret is not set")
	}

	expTime, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_TIME"))
	if err != nil {
		expTime = 3600
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"tid": tid,
			"sub": sub,
			"iss": os.Getenv("APP_URL"),
			"exp": time.Now().Add(time.Second * time.Duration(expTime)).Unix(),
			"iat": time.Now().Unix(),
		},
	)

	tokenWithSignature, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", rest_err.NewInternalServerError(fmt.Sprintf("failed to sign JWT: %s", err.Error()))
	}

	return tokenWithSignature, nil
}

func (ja *jwtAdapter) ParseToken(token string) (*jwt.Token, *rest_err.RestErr) {
	secret := os.Getenv("TOKEN_SECRET")

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(secret), nil
		}

		return nil, rest_err.NewBadRequestError("invalid token")
	})

	if err != nil {
		return nil, rest_err.NewBadRequestError("invalid token")
	}

	return parsedToken, nil
}

func (ja *jwtAdapter) TrimPrefix(auth string) (string, *rest_err.RestErr) {
	if !strings.Contains(auth, "Bearer ") {
		return "", rest_err.NewBadRequestError("Unauthorized")
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 {
		return "", rest_err.NewUnauthorizedError("malformed authorization header")
	}

	token := parts[1]

	return token, nil
}
