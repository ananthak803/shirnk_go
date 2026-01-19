package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"shrink/routes"

	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func init() {
	// Load .env file from parent directory
	godotenv.Load("../.env")
	godotenv.Load(".env")
}

func connectMongoDB() error {
	// Get MongoDB URI from environment
	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongodbURI))
	if err != nil {
		return err
	}

	// Verify connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	mongoClient = client
	fmt.Println("Connected to MongoDB successfully")
	return nil
}

func main() {
	// Connect to MongoDB
	if err := connectMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server is running",
		})
	})

	// Test route: POST /test/shrink?url=<original_url>
	// Example: POST http://localhost:8080/test/shrink?url=https://google.com
	router.POST("/test/shrink", func(c *gin.Context) {
		originalURL := c.Query("url")
		if originalURL == "" {
			c.JSON(400, gin.H{
				"error": "Missing 'url' query parameter",
				"usage": "POST /test/shrink?url=https://example.com",
			})
			return
		}

		// Redirect to /api/shorten with JSON body
		c.JSON(201, gin.H{
			"message":      "URL shortened successfully!",
			"original_url": originalURL,
			"short_url":    "Check /api/stats/<shortUrl> for stats",
		})
	})

	// URL routes
	router.POST("/api/shrink", routes.CreateURL(mongoClient))
	router.GET("/:shortUrl", routes.GetURL(mongoClient))
	router.GET("/info/:shortUrl", routes.GetURLInfo(mongoClient))
	router.GET("/api/stats/:shortUrl", routes.GetURLStats(mongoClient))

	// Start server on configured port (default 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s\n", port)
	router.Run(":" + port)
}
