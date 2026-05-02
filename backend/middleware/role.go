package middleware

import (
	"net/http"

	"be/config"
	"be/enums"
	"be/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RequireAdmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get(config.JWT_CLAIMS_KEY_USER_ID)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
			c.Abort()
			return
		}

		var userID uint
		switch v := userIDRaw.(type) {
		case float64:
			userID = uint(v)
		case uint:
			userID = v
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if user.Role != enums.Admin.String() && user.Role != enums.SuperAdmin.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
