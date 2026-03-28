package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/repo"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/utils/encrypt"
	"github.com/gin-gonic/gin"
)

func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Missing API Key",
			})
			return
		}

		settingRepo := repo.NewISettingRepo()
		storedHash, err := settingRepo.GetValueByKey("ApiKeyHash")
		if err != nil || storedHash == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "API Key not configured",
			})
			return
		}

		decryptedHash, err := encrypt.StringDecrypt(storedHash)
		if err != nil {
			decryptedHash = storedHash
		}

		incomingHash := hashApiKey(apiKey)
		if subtle.ConstantTimeCompare([]byte(incomingHash), []byte(decryptedHash)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid API Key",
			})
			return
		}

		// IP whitelist check (optional)
		allowedIPs, _ := settingRepo.GetValueByKey("ApiKeyAllowedIPs")
		if allowedIPs != "" {
			clientIP := c.ClientIP()
			if !isIPAllowed(clientIP, allowedIPs) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "IP not allowed",
				})
				return
			}
		}

		c.Next()
	}
}

func hashApiKey(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

func isIPAllowed(clientIP, allowedIPs string) bool {
	ips := strings.Split(allowedIPs, ",")
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			continue
		}
		if ip == clientIP {
			return true
		}
	}
	return false
}

// HashApiKey exports the hash function for use by installer/CLI
func HashApiKey(key string) string {
	return hashApiKey(key)
}
