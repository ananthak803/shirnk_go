package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"shrink/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// LocationData represents the response from IP geolocation API
type LocationData struct {
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Region    string  `json:"region"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// GetLocationFromIP retrieves location data based on IP address
func GetLocationFromIP(ip string) LocationData {
	location := LocationData{}

	// Use ipapi.co free API for geolocation
	url := "https://ipapi.co/" + ip + "/json/"

	// Make request with timeout
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return location
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return location
	}

	// Parse JSON response
	err = json.Unmarshal(body, &location)
	if err != nil {
		return location
	}

	return location
}

// GenerateShortURL generates a random short URL code
func GenerateShortURL(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// CreateURL creates a new shortened URL and saves it to MongoDB
func CreateURL(mongoClient *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateURLRequest

		// Bind JSON request
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request: " + err.Error(),
			})
			return
		}

		// Get MongoDB collection
		collection := mongoClient.Database("shrink").Collection("urls")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Generate short URL
		var shortURL string
		if req.CustomAlias != "" {
			// Check if custom alias already exists
			count, err := collection.CountDocuments(ctx, bson.M{"short_url": req.CustomAlias})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Database error",
				})
				return
			}
			if count > 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Custom alias already exists",
				})
				return
			}
			shortURL = req.CustomAlias
		} else {
			// Generate random short URL
			for {
				shortURL = GenerateShortURL(6)
				count, err := collection.CountDocuments(ctx, bson.M{"short_url": shortURL})
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Database error",
					})
					return
				}
				if count == 0 {
					break
				}
			}
		}

		// Create URL document
		url := models.URL{
			ID:          primitive.NilObjectID,
			OriginalURL: req.OriginalURL,
			ShortURL:    shortURL,
			CustomAlias: req.CustomAlias,
			TotalClicks: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsActive:    true,
			Clicks:      []models.Click{},
		}

		// Insert into MongoDB
		result, err := collection.InsertOne(ctx, url)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save URL",
			})
			return
		}

		// Return response
		c.JSON(http.StatusCreated, gin.H{
			"id":           result.InsertedID,
			"original_url": url.OriginalURL,
			"short_url":    url.ShortURL,
			"total_clicks": url.TotalClicks,
			"created_at":   url.CreatedAt,
			"is_active":    url.IsActive,
		})
	}
}

// GetURL retrieves a shortened URL by short code
func GetURL(mongoClient *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortURL := c.Param("shortUrl")

		collection := mongoClient.Database("shrink").Collection("urls")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var url models.URL
		err := collection.FindOne(ctx, bson.M{"short_url": shortURL}).Decode(&url)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Short URL not found",
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error",
			})
			return
		}

		// Capture click details
		click := models.Click{
			ID:        primitive.NewObjectID(),
			Timestamp: time.Now(),
			IP:        c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
			Referrer:  c.GetHeader("Referer"),
		}

		// Get location from query parameters (sent from browser Geolocation API)
		latitude := c.Query("lat")
		longitude := c.Query("lng")
		country := c.Query("country")
		city := c.Query("city")
		region := c.Query("region")

		if latitude != "" && longitude != "" {
			// Browser geolocation was provided
			var lat, lng float64
			fmt.Sscanf(latitude, "%f", &lat)
			fmt.Sscanf(longitude, "%f", &lng)
			click.Latitude = lat
			click.Longitude = lng
			click.Country = country
			click.City = city
			click.Region = region
		} else if click.IP != "127.0.0.1" && click.IP != "::1" {
			// Fallback to IP geolocation for non-localhost
			location := GetLocationFromIP(click.IP)
			click.Country = location.Country
			click.City = location.City
			click.Region = location.Region
			click.Latitude = location.Latitude
			click.Longitude = location.Longitude
		} else {
			// For localhost/testing
			click.Country = "Local"
			click.City = "Testing"
			click.Region = "Dev"
		}

		// Update URL document with new click and increment total clicks
		updateResult, err := collection.UpdateOne(
			ctx,
			bson.M{"short_url": shortURL},
			bson.M{
				"$push": bson.M{"clicks": click},
				"$inc":  bson.M{"total_clicks": 1},
				"$set":  bson.M{"updated_at": time.Now()},
			},
		)

		if err != nil {
			// Still redirect even if click tracking fails
			c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
			return
		}

		if updateResult.ModifiedCount == 0 {
			// URL was deleted or something went wrong, still redirect
			c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
			return
		}

		// Redirect to original URL
		c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
	}
}

// GetURLInfo retrieves all details of a shortened URL including clicks
func GetURLInfo(mongoClient *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortURL := c.Param("shortUrl")

		collection := mongoClient.Database("shrink").Collection("urls")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var url models.URL
		err := collection.FindOne(ctx, bson.M{"short_url": shortURL}).Decode(&url)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Short URL not found",
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":           url.ID.Hex(),
			"original_url": url.OriginalURL,
			"short_url":    url.ShortURL,
			"custom_alias": url.CustomAlias,
			"total_clicks": url.TotalClicks,
			"created_at":   url.CreatedAt,
			"updated_at":   url.UpdatedAt,

			"is_active":    url.IsActive,
			"clicks_count": len(url.Clicks),
			"clicks":       url.Clicks,
		})
	}
}

// GetURLStats retrieves statistics for a shortened URL
func GetURLStats(mongoClient *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortURL := c.Param("shortUrl")

		collection := mongoClient.Database("shrink").Collection("urls")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var url models.URL
		err := collection.FindOne(ctx, bson.M{"short_url": shortURL}).Decode(&url)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Short URL not found",
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":           url.ID,
			"original_url": url.OriginalURL,
			"short_url":    url.ShortURL,
			"total_clicks": url.TotalClicks,
			"created_at":   url.CreatedAt,
			"updated_at":   url.UpdatedAt,
			"is_active":    url.IsActive,
			"clicks":       url.Clicks,
		})
	}
}
