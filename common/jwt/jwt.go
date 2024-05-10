package jwt

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

type JWT struct {
	secretKey     string
	tokenDuration time.Duration
}

type CustomClaims struct {
	StandardClaims jwt.StandardClaims
	Nim            string `json:"nim"`
	Role           uint32 `json:"role"`
}

func NewJWT(secretKey string, tokenDuration time.Duration) *JWT {
	return &JWT{
		secretKey,
		tokenDuration,
	}
}

func (j *JWT) GenerateToken(nim string, role uint32) (string, error) {
	claims := &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(j.tokenDuration).Unix(),
		},
		Nim:  nim,
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWT) Verify(accessToken string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				log.Println("ERROR: [JWT - Verify] Unexpected signing method")
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(j.secretKey), nil
		},
	)

	if err != nil {
		log.Println("ERROR: [JWT - Verify] Error while parsing token:", err)
		// return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		log.Println("ERROR: [JWT - Verify] Invalid token claims")
		// return nil, status.Error(codes.Unauthenticated, "invalid token claims")
		return nil, fmt.Errorf("invalid token claims")
	}

	if err := claims.Valid(); err != nil {
		log.Println("ERROR: [JWT - Verify] Invalid token:", err)
		return nil, err
	}

	return claims, nil
}

func (c *CustomClaims) Valid() error {
	// check if the token has expired.
	if time.Now().Unix() > c.StandardClaims.ExpiresAt {
		// return status.Error(codes.Unauthenticated, "token has expired")
		return fmt.Errorf("token has expired")
	}

	return nil
}
