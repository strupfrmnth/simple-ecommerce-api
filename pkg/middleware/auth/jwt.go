package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// reference: https://golang-jwt.github.io/jwt/usage/create/
// reference: https://pkg.go.dev/github.com/golang-jwt/jwt#MapClaims
// reference: https://pkg.go.dev/github.com/golang-jwt/jwt#example-Parse-Hmac

func CreateToken(id uint, name string) string {
	secretkey := []byte("secret") // load from somewhere, ex. environment

	payload := jwt.MapClaims{
		"exp":    time.Now().Add(time.Hour * 3).Unix(),
		"iat":    time.Now().Unix(),
		"userID": id,
		"name":   name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// create hard decoded token
	signature, err := token.SignedString(secretkey)
	if err != nil {
		log.Fatalln(err, "sign error")
	}
	return signature
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	// parse
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})
	return token, err
}

// client how to show jwt token?
// 1. put in http header
// 2. post request body
// 3. uri query parameter(not recommend)

// we assume the token will be put in http header
// create gin handler func for authorization
func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizeValue := c.GetHeader("Authorization")
		if authorizeValue == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not found"})
		} else {
			authorizeValueSlice := strings.Fields(authorizeValue)
			// why bearer? https://medium.com/@arunchaitanya/wtf-is-bearer-token-an-in-depth-explanation-60695b581928
			bearerToken := authorizeValueSlice[0]
			if bearerToken != "Bearer" || len(authorizeValueSlice) != 2 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorize header value"})
			} else {
				stringToken := authorizeValueSlice[1]
				if jwtToken, err := verifyToken(stringToken); err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				} else {
					if claims, ok := jwtToken.Claims.(jwt.MapClaims); !ok {
						c.AbortWithStatus(http.StatusUnauthorized)
					} else {
						if jwtToken.Valid {
							// add userid to request context
							c.Set("userID", claims["userID"])
							fmt.Println("Authorize", claims["name"])
							fmt.Println("Authorize", claims["userID"])
						} else {
							c.AbortWithStatus(http.StatusUnauthorized)
						}
					}
				}
			}
		}
	}
}
