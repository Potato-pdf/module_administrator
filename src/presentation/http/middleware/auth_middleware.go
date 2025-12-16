package middleware

import (
	"net/http"

	"luminaMO-orq-modAI/src/domain/services"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService services.AuthService
}

func NewAuthMiddleware(authService services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// ValidateGatewayEntrada valida requests del Gateway Entrada
func (m *AuthMiddleware) ValidateGatewayEntrada() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API Key"})
			c.Abort()
			return
		}

		if err := m.authService.ValidateAPIKey(apiKey, "gateway-entrada"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateModule valida responses de módulos
func (m *AuthMiddleware) ValidateModule() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing API Key"})
			c.Abort()
			return
		}

		// Extraer módulo del header
		moduleID := c.GetHeader("X-Module-ID")
		if moduleID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Module ID"})
			c.Abort()
			return
		}

		// Validar API key del módulo
		source := "module-" + moduleID
		if err := m.authService.ValidateAPIKey(apiKey, source); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Module API Key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
