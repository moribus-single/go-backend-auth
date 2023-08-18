package services

import (
	"app/models"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtService struct {
	SecretKey       string
	AccessLifetime  int
	RefreshLifetime int
}

func (s JwtService) Generate(guid string, tokenCounter int) (string, string, error) {
	// expiry time
	accessExp := time.Now().Add(time.Minute * time.Duration(s.AccessLifetime)).Unix()
	refreshExp := time.Now().Add(time.Hour * time.Duration(s.RefreshLifetime)).Unix()

	// getting token instances
	access := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_guid":     guid,
		"token_counter": tokenCounter,
		"exp":           accessExp,
		"type":          "access",
	})

	combined := guid + string(rune(tokenCounter)) + string(rune(refreshExp))
	refreshToken := sign([]byte(combined), []byte(s.SecretKey))

	// signing jwt tokens
	accessToken, err := access.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "Access err", "", err
	}

	return accessToken, refreshToken, nil
}

func sign(msg, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(msg)

	return hex.EncodeToString(h.Sum(nil))
}

func GetJwtService(accessLifetime, refreshLifetime int, secret string) JwtService {
	return JwtService{
		SecretKey:       secret,
		AccessLifetime:  accessLifetime,
		RefreshLifetime: refreshLifetime,
	}
}
