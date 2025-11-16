package jwt

import "github.com/golang-jwt/jwt/v4"

func GenerateToken(role, jwtSecret string) (string, error) {
	claims := jwt.MapClaims{
		"role": role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
