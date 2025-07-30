package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"alfredo/tabunganku/config"
	"alfredo/tabunganku/pkg/dtos"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"

	defaultAccessTokenExpiry  = 15
	defaultRefreshTokenExpiry = 60 * 24

	revokedTokenExpiration = 24 * 7
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrMalformedToken   = errors.New("malformed token")
	ErrTokenSigning     = errors.New("error signing token")
	ErrRefreshTokenSign = errors.New("error signing refresh token")
)

type jwtServiceImpl struct {
	redisService RedisService
}

// IsTokenExpired checks if a token has expired
func (j *jwtServiceImpl) IsTokenExpired(token string) bool {
	newToken, err := j.ValidateToken(token)
	if err != nil {
		return true
	}

	claims, ok := newToken.Claims.(jwt.MapClaims)
	if !ok {
		return true
	}

	expiredTime := time.Unix(int64(claims["exp"].(float64)), 0)
	return expiredTime.Before(time.Now())
}

// Revoke invalidates a token by storing it in Redis with an expiration time
func (j *jwtServiceImpl) Revoke(token string) error {
	// Store token in Redis with expiration to automatically clean up old tokens
	if err := j.redisService.Set(token, true); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

// IsTokenRevoked checks if a token has been revoked
func (j *jwtServiceImpl) IsTokenRevoked(token string) bool {
	res, err := j.redisService.Get(token)
	if err != nil {
		return false
	}
	return res != ""
}

// GenerateToken creates both access and refresh tokens for a user
func (j *jwtServiceImpl) GenerateToken(userUuid string, tokens string) (dtos.GenerateTokenResponse, error) {
	// Get JWT configuration
	secretKey := []byte(config.JwtSecret)

	// Get token expiration times from config or use defaults
	accessExpiry, refreshExpiry := j.getTokenExpiryTimes()

	// Generate access token
	accessToken, err := j.createToken(userUuid, tokens, accessExpiry, string(AccessToken), secretKey)
	if err != nil {
		return dtos.GenerateTokenResponse{}, fmt.Errorf("failed to create access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := j.createToken(userUuid, tokens, refreshExpiry, string(RefreshToken), secretKey)
	if err != nil {
		return dtos.GenerateTokenResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return dtos.GenerateTokenResponse{
		TokenType:    "Bearer",
		ExpiresIn:    int(accessExpiry) * 60, // Convert minutes to seconds
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken verifies a JWT token's signature and format
func (j *jwtServiceImpl) ValidateToken(token string) (*jwt.Token, error) {
	secretKey := []byte(config.JwtSecret)

	// Basic format validation
	segments := strings.Split(token, ".")
	if len(segments) != 3 {
		return nil, ErrMalformedToken
	}

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})
}

// GetUserIdFromToken extracts the user ID from a valid token
func (j *jwtServiceImpl) GetUserIdFromToken(token string) (string, error) {
	newToken, err := j.ValidateToken(token)
	if err != nil {
		return "", err
	}

	claims, ok := newToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid user_id in token")
	}

	return userID, nil
}

// Helper methods

// getTokenExpiryTimes returns the configured expiration times for access and refresh tokens
func (j *jwtServiceImpl) getTokenExpiryTimes() (accessExpiry, refreshExpiry int64) {
	accessTokenExpiredTime := config.JwtTokenAccessExpire
	refreshTokenExpiredTime := config.JwtTokenRefreshExpire

	if accessTokenExpiredTime == "" || refreshTokenExpiredTime == "" {
		return defaultAccessTokenExpiry, defaultRefreshTokenExpiry
	}

	//accessExpiry = int64(helper.ParseStringToInt(accessTokenExpiredTime))
	accessExpiry, _ = strconv.ParseInt(accessTokenExpiredTime, 10, 64)
	refreshExpiry, _ = strconv.ParseInt(refreshTokenExpiredTime, 10, 64)
	return
}

// createToken generates a signed JWT token with the given parameters
func (j *jwtServiceImpl) createToken(userUuid string, tokens string, expiry int64, tokenType string, secretKey []byte) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userUuid,
		"tokens":  tokens,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expiry))).Unix(),
		"iat":     time.Now().Unix(),
		"type":    tokenType,
	}

	// Add exp_in only for access tokens
	if tokenType == string(AccessToken) {
		claims["exp_in"] = jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expiry))).Second()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", ErrTokenSigning
	}

	return signedToken, nil
}

type JwtService interface {
	IsTokenExpired(token string) bool
	Revoke(token string) error
	IsTokenRevoked(token string) bool
	GenerateToken(userUuid string, tokens string) (dtos.GenerateTokenResponse, error)
	ValidateToken(token string) (*jwt.Token, error)
	GetUserIdFromToken(token string) (string, error)
}

func NewJwtService(redisService RedisService) JwtService {
	return &jwtServiceImpl{
		redisService: redisService,
	}
}
