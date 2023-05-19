package user

import (
	"net/http"
	"slingshot/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

type CustomClaims struct {
	Uid string `json:"uid"`
	jwt.RegisteredClaims
}

func CheckJwtToken(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "No token provided",
		})
	}

	if !strings.Contains(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Invalid token format",
		})
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", -1)

	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.Server.JwtSecretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"message": "Invalid signature",
			})
		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Error parsing token",
		})
	}

	if !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Invalid token",
		})
	}

	c.Set("uid", claims.Uid)
	c.Set("response", map[string]interface{}{
		"success": true,
	})
	return nil
}

func CreateJwtToken(uid string, timeduration string) (string, error) {
	t, err := time.ParseDuration(timeduration)
	if err != nil {
		return "", err
	}

	expiresAt := jwt.NumericDate{Time: time.Now().Add(t)}

	claims := &CustomClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &expiresAt,
			Issuer:    config.Cfg.Server.JwtIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.Server.JwtSecretKey))
}
