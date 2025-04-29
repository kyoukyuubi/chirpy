package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// hash generate
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// check hash
func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

// make JWT token
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	})

	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// validate JWT token
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		id, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, err
		}
		return id, nil
	}

	return uuid.Nil, fmt.Errorf("invalid token")
}

// get token from header
func GetBearerToken(headers http.Header) (string, error) {
	headerList := headers.Values("Authorization")
	token := ""
	for _, header := range headerList {
		words := strings.Fields(header)
		if len(words) == 0 {
			return "", fmt.Errorf("no token found")
		}

		if len(words) == 2 && words[0] == "Bearer" {
			token = words[1]
			break
		}
	}

	if token != "" {
		return token, nil
	}
	return "", fmt.Errorf("no token found")
}

// make refresh token
func MakeRefreshToken() (string, error) {
	// make a key which is used to generate a 256-bit of random data
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	// convert the random data into hex string
	hex := hex.EncodeToString(key)
	return hex, nil
}

// get polka get
func GetAPIKey(headers http.Header) (string, error) {
	headerList := headers.Values("Authorization")
	key := ""
	for _, header := range headerList {
		words := strings.Fields(header)
		if len(words) == 0 {
			return "", fmt.Errorf("no key found")
		}

		if len(words) == 2 && words[0] == "ApiKey" {
			key = words[1]
			break
		}
	}
	if key != "" {
		return key, nil
	}
	return "", fmt.Errorf("no token found")
}