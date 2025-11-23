package httpmw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anton1ks96/college-app-core/internal/domain"
	"github.com/anton1ks96/college-app-core/pkg/logger"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	client        *http.Client
	validationURL string
}

type ValidationRequest struct {
	Token string `json:"token"`
}

type ValidationResponse struct {
	Valid bool         `json:"valid"`
	User  *domain.User `json:"user,omitempty"`
	Error string       `json:"error,omitempty"`
}

func NewAuthMiddleware(authServiceURL string, timeout time.Duration) *AuthMiddleware {
	return &AuthMiddleware{
		client: &http.Client{
			Timeout: timeout,
		},
		validationURL: fmt.Sprintf("%s/api/v1/app/validate", strings.TrimSuffix(authServiceURL, "/")),
	}
}

func (m *AuthMiddleware) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		valid, userID, err := m.validateWithAuthService(token)
		if err != nil {
			logger.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token validation failed",
			})
			c.Abort()
			return
		}

		if !valid {
			logger.Warn("Invalid or expired token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func (m *AuthMiddleware) validateWithAuthService(token string) (bool, string, error) {
	reqBody := ValidationRequest{
		Token: token,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return false, "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", m.validationURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := m.client.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	var validationResp ValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&validationResp); err != nil {
		return false, "", fmt.Errorf("failed to decode response: %w", err)
	}

	if !validationResp.Valid {
		return false, "", nil
	}

	if validationResp.User == nil {
		return false, "", fmt.Errorf("valid response but user data is missing")
	}

	return true, validationResp.User.ID, nil
}

func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}
