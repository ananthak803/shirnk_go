package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// URL represents a shortened URL with click analytics
type URL struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OriginalURL string             `bson:"original_url" json:"original_url"`
	ShortURL    string             `bson:"short_url" json:"short_url"`
	CustomAlias string             `bson:"custom_alias,omitempty" json:"custom_alias,omitempty"`
	TotalClicks int64              `bson:"total_clicks" json:"total_clicks"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	Clicks      []Click            `bson:"clicks,omitempty" json:"clicks,omitempty"`
}

// Click represents a single click event with all details
type Click struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Timestamp      time.Time          `bson:"timestamp" json:"timestamp"`
	IP             string             `bson:"ip,omitempty" json:"ip,omitempty"`
	Country        string             `bson:"country,omitempty" json:"country,omitempty"`
	City           string             `bson:"city,omitempty" json:"city,omitempty"`
	Region         string             `bson:"region,omitempty" json:"region,omitempty"`
	Latitude       float64            `bson:"latitude,omitempty" json:"latitude,omitempty"`
	Longitude      float64            `bson:"longitude,omitempty" json:"longitude,omitempty"`
	Browser        string             `bson:"browser,omitempty" json:"browser,omitempty"`
	BrowserVersion string             `bson:"browser_version,omitempty" json:"browser_version,omitempty"`
	OS             string             `bson:"os,omitempty" json:"os,omitempty"`
	OSVersion      string             `bson:"os_version,omitempty" json:"os_version,omitempty"`
	DeviceType     string             `bson:"device_type,omitempty" json:"device_type,omitempty"`
	UserAgent      string             `bson:"user_agent,omitempty" json:"user_agent,omitempty"`
	Referrer       string             `bson:"referrer,omitempty" json:"referrer,omitempty"`
}

// CreateURLRequest is the request body for creating a new shortened URL
type CreateURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	CustomAlias string `json:"custom_alias,omitempty"`
}
