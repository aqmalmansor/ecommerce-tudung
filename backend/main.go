package main

import (
	"log"
	"strconv"

	"be/config"
	"be/handlers"
	"be/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := config.MigrateDatabase(db); err != nil {
		log.Fatal("Migration failed:", err)
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", config.CLIENT_URL)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	authHandler := handlers.NewAuthHandler(db)

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			sqlDB, err := db.DB()
			if err != nil {
				c.JSON(500, gin.H{
					"status":   "error",
					"database": "disconnected",
					"error":    err.Error(),
				})
				return
			}

			if err := sqlDB.Ping(); err != nil {
				c.JSON(500, gin.H{
					"status":   "error",
					"database": "ping_failed",
					"error":    err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{
				"status":   "ok",
				"database": "connected",
				"service":  "be",
			})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/user/profile", authHandler.GetProfile)
		}
	}

	port := config.PORT

	log.Printf("Server starting on port %d", port)
	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
