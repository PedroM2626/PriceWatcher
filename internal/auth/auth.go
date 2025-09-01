package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

var (
	// ErrInvalidCredentials is returned when the credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// JWT signing key (in production, this should be loaded from environment variables)
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

// User represents an authenticated user
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Hashed password
	CreatedAt time.Time `json:"created_at"`
}

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token for the given user
func GenerateJWT(user *User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pricewatcher",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// HashPassword hashes a password using Argon2id
func HashPassword(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Generate the hash
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		3,      // time
		32*1024, // memory
		4,      // threads
		32,     // key length
	)

	// Encode the hash and salt
	encodedHash := base64.RawStdEncoding.EncodeToString(append(salt, hash...))
	return encodedHash, nil
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, encodedHash string) (bool, error) {
	// Decode the hash
	decoded, err := base64.RawStdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Extract the salt and hash
	if len(decoded) < 16 {
		return false, errors.New("invalid hash format")
	}
	salt := decoded[:16]
	hash := decoded[16:]

	// Generate the comparison hash
	comparisonHash := argon2.IDKey(
		[]byte(password),
		salt,
		3,      // time
		32*1024, // memory
		4,      // threads
		32,     // key length
	)

	// Compare the hashes in constant time
	if subtle.ConstantTimeCompare(hash, comparisonHash) == 1 {
		return true, nil
	}

	return false, nil
}

// ExtractToken extracts the JWT token from the Authorization header
func ExtractToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("authorization header format must be 'Bearer {token}'")
	}

	return parts[1], nil
}
