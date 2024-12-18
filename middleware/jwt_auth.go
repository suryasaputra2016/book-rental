package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// custom jwt claim
type jwtCustomClaims struct {
	ID    int    `json:"id"`
	Email string `json:"name"`
	jwt.RegisteredClaims
}

// GetUserEmail accept echo context and return email data from jwt encoded token
func GetUserEmail(c interface{}) string {
	user := c.(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	return claims.Email
}

// GetUserID accept echo context and return id data from jwt encoded token
func GetUserID(c interface{}) int {
	user := c.(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	return claims.ID
}

// GenerateTokenString generate token string with custom claims and return encoded token and error
func GenerateTokenString(id int, email string) (string, error) {
	// Set custom claim
	claims := &jwtCustomClaims{
		id,
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Authorizition return middleware jwt authorization
func Authorization() echo.MiddlewareFunc {
	return echojwt.WithConfig(
		echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(jwtCustomClaims)
			},
			SigningKey: []byte(os.Getenv("JWT_SECRET")),
		},
	)
}
