package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pdkonovalov/auth-server/pkg/config"
)

type JwtGenerator struct {
	secret []byte
}

func Init(config *config.Config) (*JwtGenerator, error) {
	if config.JwtSecret == "" {
		return nil, fmt.Errorf("jwt secret is empty")
	}
	return &JwtGenerator{[]byte(config.JwtSecret)}, nil
}

type accessTokenClaims struct {
	Ip string `json:"ip"`
	jwt.RegisteredClaims
}

type refreshTokenClaims struct {
	jwt.RegisteredClaims
}

func (gen *JwtGenerator) GenerateAccessToken(ip, jti string) (string, error) {
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, accessTokenClaims{
		ip,
		jwt.RegisteredClaims{
			ID: jti,
		}}).SignedString(gen.secret)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (gen *JwtGenerator) GenerateRefreshToken(jti string) (string, error) {
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshTokenClaims{
		jwt.RegisteredClaims{
			ID: jti,
		}}).SignedString(gen.secret)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (gen *JwtGenerator) ValidateAccessToken(tokenStr string) (string, string, bool) {
	accessToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return gen.secret, nil
	})
	if err != nil {
		return "", "", false
	}
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	if len(accessClaims) != 2 {
		return "", "", false
	}
	for k := range accessClaims {
		if k != "jti" && k != "ip" {
			return "", "", false
		}
	}
	jti := accessClaims["jti"].(string)
	ip := accessClaims["ip"].(string)
	return ip, jti, true
}

func (gen *JwtGenerator) ValidateRefreshToken(tokenStr string) (string, bool) {
	refreshToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return gen.secret, nil
	})
	if err != nil {
		return "", false
	}
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	if len(refreshClaims) != 1 {
		return "", false
	}
	if _, ok := refreshClaims["jti"]; !ok {
		return "", false
	}
	jti := refreshClaims["jti"].(string)
	return jti, true
}
