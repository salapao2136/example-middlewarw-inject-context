package middleware

import (
	"fmt"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type middleware struct {
}

type claims struct {
	ChannelID string `json:"channelId"`
	UserID    string `json:"userId"`
	jwt.StandardClaims
}

// NewMiddleware is init Middleware
func NewMiddleware() *middleware {
	return &middleware{}
}

func (m *middleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.DefaultQuery("token", "")
		if token == "" {
			respondWithError(401, "invalid token", c)
			return
		}
		m.decodeJWT(token, c)
	}
}

func respondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}
	c.JSON(code, resp)
	c.Abort()
}

func (m middleware) decodeJWT(tokenString string, c *gin.Context) {

	secretKey := os.Getenv("SECRET_KEY")
	var secret []byte
	if secretKey == "" {
		secret = []byte("secret")
	} else {
		secret = []byte(secretKey)
	}

	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		c.Set("claims", claims)
		c.Next()
	} else {
		respondWithError(401, err.Error(), c)
		return
	}
}