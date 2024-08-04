package authentication

import (
	"errors"
	"net/http"
	"projectgo/model"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"

	//"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var cjwtkey = []byte("crekey")

func GenerateCreToken(username string) (string, error) {
	//setting token expiration time

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &model.CrediterClaims{
		Username:       username,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(ujwtkey)
}

// verify Admin Token
func CreAuthentication(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.CrediterClaims{}, func(token *jwt.Token) (interface{}, error) {
		return ujwtkey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*model.CrediterClaims); ok && token.Valid {
		return claims.Username, nil
	}
	return "", errors.New("invalid token")
}

//Admin Auth middleware

func CreAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing the authorization header"})
			return
		}

		authHeader := strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer"))

		username, err := CreAuthentication(authHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("username", username)
		c.Next()
	}
}