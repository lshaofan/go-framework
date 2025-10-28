package helper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTResponse struct {
	AccessToken  string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredAt    time.Time `json:"expired_at"`
}

// JwtOption 定义配置选项函数类型
type JwtOption func(*JWTUtil)

// WithOAuthJWTConfigIssuer 设置 JWT 的签发者
func WithOAuthJWTConfigIssuer(issuer string) JwtOption {
	return func(config *JWTUtil) {
		config.Issuer = issuer
	}
}

// WithOAuthJWTConfigExpired 设置 JWT 的过期时间
func WithOAuthJWTConfigExpired(expired time.Duration) JwtOption {
	return func(config *JWTUtil) {
		config.Expired = expired
	}
}

// WithOAuthJWTConfigRefreshExpired 设置 JWT 的刷新时间
func WithOAuthJWTConfigRefreshExpired(refreshExpired time.Duration) JwtOption {
	return func(config *JWTUtil) {
		config.RefreshExpired = refreshExpired
	}
}

// WithOAuthJWTSecretKey 设置 JWT 的密钥
func WithOAuthJWTSecretKey(secretKey string) JwtOption {
	return func(config *JWTUtil) {
		config.SecretKey = secretKey
	}
}

// JWTUtil 是一个通用的 JWT 配置结构体
type JWTUtil struct {
	SecretKey      string        // SecretKey 是用于签名 JWT 的密钥
	Issuer         string        // Issuer 是 JWT 的签发者
	Expired        time.Duration // Expired 是 JWT 的过期时
	RefreshExpired time.Duration // Refresh 是 JWT 的刷新时间刷新时间是加在过期时间的基础上的
}

func NewJWTUtil(secretKey string, options ...JwtOption) *JWTUtil {
	config := &JWTUtil{
		SecretKey:      secretKey,
		Expired:        time.Hour * 24,
		RefreshExpired: time.Hour * 24 * 30,
	}

	for _, option := range options {
		option(config)
	}

	return config
}

// 生成JWT Token
func (c *JWTUtil) GenerateToken(data any, expired time.Duration) (string, error) {
	// 解析data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("解析data失败: %w", err)
	}
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expired)),
		Issuer:    c.Issuer,
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
		Subject:   string(jsonData),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.SecretKey))
}

// GetToken 获取JWT Token
func (c *JWTUtil) GetToken(data any) (*JWTResponse, error) {
	token, err := c.GenerateToken(data, c.Expired)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	refreshToken, err := c.GenerateToken(data, c.Expired+c.RefreshExpired)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return &JWTResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(c.Expired),
	}, nil
}

// ParseToken 解析JWT Token
func (c *JWTUtil) ParseToken(token string) (*jwt.RegisteredClaims, error) {
	claims, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return claims.Claims.(*jwt.RegisteredClaims), nil
}

// ValidateToken 验证token是否有效
func (c *JWTUtil) ValidateToken(token string) bool {
	_, err := c.ParseToken(token)
	return err == nil
}
