package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/gradusp/crispy/internal/secrets"
)

const (
	apiKeyHeader = "LBOS_API_KEY" // FIXME: not the best place for that
)

// API key auth middleware
func AuthAPIKey(secretId string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get(apiKeyHeader)

		secret, err := secrets.GetSecret(secretId)
		if err != nil {
			//log.Println("failed to get secret") // TODO: implement in core logger
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		//log.Println("comparing secret with provided key", secret, key) // TODO: implement in core logger

		if secret != key {
			//log.Println("key doesn't match!") // TODO: implement in core logger
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		//log.Println("no error during check") // TODO: implement in core logger
		c.Next()
	}
}
