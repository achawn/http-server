package auth

import "golang.org/x/crypto/bcrypt"
import "github.com/google/uuid"
import "time"
import "github.com/golang-jwt/jwt/v5"
import "fmt"
import "strings"
import "net/http"

func HashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pw), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer: "chirpy",
			IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject: userID.String(),
		},
	)
	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token claims")
	}

	id := claims.Subject
	vjwt, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}

	return vjwt, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token_string := headers.Get("Authorization")
	if token_string == "" {
		return "", fmt.Errorf("no token found")
	}

	if !strings.HasPrefix(token_string, "Bearer ") {
		return "", fmt.Errorf("malformed token")
	}

	trimmed := strings.TrimPrefix(token_string, "Bearer ")
	if trimmed == "" {
		return "", fmt.Errorf("empty token")
	}

	return trimmed, nil

}
