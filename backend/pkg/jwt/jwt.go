package jwt

import (
	"crypto/rsa"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTInterface interface {
	JWTReader
	JWTWriter
}

type JWTReader interface {
	IsValidToken(token string) bool
	GetClaims(token string) (*Auth, error)
}

type JWTWriter interface {
	Sign(data Auth) (string, error)
}

type JWT struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	Issuer        string
	TokenDuration time.Duration
}

type NewJWTOptions struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	Issuer        string
	TokenDuration time.Duration
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type Auth struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func New(opts *NewJWTOptions) *JWT {
	return &JWT{
		PrivateKey:    opts.PrivateKey,
		PublicKey:     opts.PublicKey,
		Issuer:        opts.Issuer,
		TokenDuration: opts.TokenDuration,
	}
}

func (j *JWT) Sign(data Auth) (string, error) {
	now := time.Now()
	jti := uuid.New().String() // Generate unique JWT ID for revocation capability

	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    j.Issuer,
			Subject:   data.UserID,
			Audience:  []string{j.Issuer},
			ExpiresAt: jwt.NewNumericDate(now.Add(j.TokenDuration)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserID:   data.UserID,
		FullName: data.FullName,
		Email:    data.Email,
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		claims,
	)

	signedToken, err := token.SignedString(j.PrivateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JWT) IsValidToken(token string) bool {
	if token == "" {
		return false
	}

	token = strings.TrimPrefix(token, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return j.PublicKey, nil
	})
	if err != nil || !parsedToken.Valid {
		return false
	}

	return true
}

func (j *JWT) GetClaims(token string) (*Auth, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// verify the signing method to prevent algorithm confusion attacks
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if token.Method != jwt.SigningMethodRS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.PublicKey, nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token: %v", err.Error())
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token when getting claims: %v", err.Error())
	}

	return &Auth{
		UserID:   fmt.Sprintf("%s", claims["user_id"]),
		FullName: fmt.Sprintf("%s", claims["full_name"]),
		Email:    fmt.Sprintf("%s", claims["email"]),
	}, nil
}
