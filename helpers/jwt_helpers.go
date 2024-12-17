package helpers

import (
	"encoding/base64"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"auth-rest-api/constants"
	"auth-rest-api/resources"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	Token     string
	SourceApp string
	Subject   string
}

type TokenHeaders struct {
	Exp     int    `json:"exp"`
	Iat     int    `json:"iat"`
	Subject string `json:"subject"`
}

type Auth struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
}

func ValidateToken(token string) (string, error) { //TODO: use claims.OmneManagerID for gm1, gm2, gm3. gm4

	env := os.Getenv("GO_ENV")
	secretKey := resources.GetConfig().GetString("config." + env + ".SECRET_KEY")

	tokens, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		expiryErr, _ := err.(*jwt.ValidationError)
		if expiryErr.Errors == jwt.ValidationErrorExpired {
			return "", err
		}
		return "", err
	}

	claims := tokens.Claims.(*Claims)

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) <= 0 {
		return "", err
	}

	if claims.Subject == "" {
		return "", err
	}

	return claims.Subject, nil
}

func GenerateJWT(userId string) (string, error) {

	iat := time.Now().Unix()
	exp := time.Now().Add((7 * 24) * time.Hour)
	atClaims := jwt.MapClaims{}
	atClaims["iat"] = iat
	atClaims["subject"] = userId
	atClaims["exp"] = exp.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, atClaims)
	env := os.Getenv("GO_ENV")
	secretKey := resources.GetConfig().GetString("config." + env + ".SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("GenerateJWT Error in JWT token generation: ", err)
		return "", err
	}
	return tokenString, nil
}

func ExtractTokenHeader(tokenString string) ([]byte, error) {

	// Split the token into its parts (header, payload, signature)
	parts := strings.Split(tokenString, ".")

	// var tokenHeaders TokenHeaders
	var tokenHeadersBytes []byte

	if len(parts) < 2 {
		log.Println("ExtractTokenHeader invalid authtoken size: ", len(parts))
		return tokenHeadersBytes, errors.New(constants.ErrorCodeMap[constants.InvalidToken])
	}

	// Decode the payload (claims)
	tokenHeadersBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Println("ExtractTokenHeader Error decoding token err:", err)
		return tokenHeadersBytes, err
	}

	return tokenHeadersBytes, nil

}
