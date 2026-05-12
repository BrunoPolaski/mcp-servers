package jwt

import (
	"github.com/BrunoPolaski/go-rest-err/rest_err"
	"github.com/golang-jwt/jwt/v5"
)

// JWT is the interface for JWT operations
type JWT interface {
	GenerateToken(tid, sub string) (string, *rest_err.RestErr)
	ParseToken(token string) (*jwt.Token, *rest_err.RestErr)
	TrimPrefix(auth string) (string, *rest_err.RestErr)
}
