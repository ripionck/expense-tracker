package helpers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SignedDetails struct {
	Email    string
	Username string
	Name     string
	Uid      string
}

var secretKey string = os.Getenv("SECRET_KEY")

// Encoding JWT
func JWTEncode(payload map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	for key, val := range payload {
		claims[key] = val
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// Decoding JWT
func JWTDecode(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("unable to parse claims")
	}

	return claims, nil
}

// Validate JWT
func JWTValidate(tokenStr string) (bool, error) {
	claims, err := JWTDecode(tokenStr)
	if err != nil {
		return false, err
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, errors.New("invalid expiration time")
	}

	if time.Now().Unix() > int64(exp) {
		return false, errors.New("token has expired")
	}

	return true, nil
}

// Generate tokens
func GenerateTokens(email, username, name, uid string) (string, string, error) {
	now := time.Now().Unix()

	// Access token payload
	accessPayload := map[string]interface{}{
		"email":    email,
		"username": username,
		"name":     name,
		"uid":      uid,
		"exp":      time.Now().Add(24 * time.Hour).Unix(), // Expires in 24 hours
		"iat":      now,
	}

	// Refresh token payload
	refreshPayload := map[string]interface{}{
		"exp": time.Now().Add(168 * time.Hour).Unix(), // Expires in 7 days
		"iat": now,
	}

	// Generate access Token
	accessToken, err := JWTEncode(accessPayload)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := JWTEncode(refreshPayload)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

// Validate the given JWT token
func ValidateToken(signedToken string) (*SignedDetails, string) {
	claims, err := JWTDecode(signedToken)
	if err != nil {
		return nil, fmt.Sprintf("failed to parse token: %v", err)
	}

	// Validate expiration
	isValid, err := JWTValidate(signedToken)
	if err != nil || !isValid {
		return nil, "Token is expired or invalid"
	}

	// Convert claims to SignedDetails
	details := &SignedDetails{
		Email:    claims["email"].(string),
		Username: claims["username"].(string),
		Name:     claims["name"].(string),
		Uid:      claims["uid"].(string),
	}

	return details, ""
}
