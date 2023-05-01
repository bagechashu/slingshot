package user

import (
	"net/http"
	"slingshot/config"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	echo "github.com/labstack/echo/v4"
)

type CustomClaims struct {
	Uid string `json:"uid"`
	jwt.StandardClaims
}

func JwtMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			// skip jwt middleware
			if c.Get("skip_jwt") != nil {
				return err
			}

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

			return err
		}
	}
}

func CreateToken(uid string, td time.Duration) (string, error) {
	claims := &CustomClaims{
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(td).Unix(),
			// TODO: Issuer learn and set
			Issuer: "Slingshot",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.Server.JwtSecretKey))
}

// userRole, err := e.GetRoleForUser(strconv.Itoa(int(claims.Uid)))
// if err != nil {
// 	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 		"success": false,
// 		"message": err.Error(),
// 	})
// }
// obj := c.Request().URL.Path
// act := c.Request().Method
// allowed, err := e.Enforce(userRole, obj, act)
// if err != nil {
// 	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 		"success": false,
// 		"message": err.Error(),
// 	})
// }
// if !allowed {
// 	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
// 		"success": false,
// 		"message": "Unauthorized access attempt",
// 	})
// }
