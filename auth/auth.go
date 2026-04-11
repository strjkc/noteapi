package auth

import (
	"runtime"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(bearer string) (string, error) {
	// get from string
	parts := strings.Split(bearer, " ")
	tokenString := strings.Trim(parts[1], " ")
	// validate
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte("teststring"), nil
	})
	if err != nil {
		return "", err
	}
	// TODO: check if token expired
	subj, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	return subj, err
}

func CreateToken(userID string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "notesapi",
		IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
		ExpiresAt: &jwt.NumericDate{Time: time.Now().UTC().Add(time.Hour * 2)},
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("teststring"))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func HashPassword(pass string) (string, error) {
	params := argon2id.Params{
		Memory:      128 * 1024,
		Iterations:  10,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  16,
		KeyLength:   32,
	}
	hashedPassword, err := argon2id.CreateHash(pass, &params)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func ValidatePassword(pass, encryptedPass string) (bool, error) {
	return argon2id.ComparePasswordAndHash(pass, encryptedPass)
}
