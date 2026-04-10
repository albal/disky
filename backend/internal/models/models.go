package models

import (
	"time"

	"github.com/google/uuid"
)

// ─── User ───────────────────────────────────────────────────────────────────

type User struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Email         string     `json:"email" db:"email"`
	Name          string     `json:"name" db:"name"`
	AvatarURL     string     `json:"avatar_url" db:"avatar_url"`
	GoogleID      *string    `json:"-" db:"google_id"`
	MicrosoftID   *string    `json:"-" db:"microsoft_id"`
	AppleID       *string    `json:"-" db:"apple_id"`
	EmailVerified bool       `json:"email_verified" db:"email_verified"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// ─── Locale ──────────────────────────────────────────────────────────────────

type Locale struct {
	Code           string `json:"code" db:"code"`
	Name           string `json:"name" db:"name"`
	CurrencyCode   string `json:"currency_code" db:"currency_code"`
	CurrencySymbol string `json:"currency_symbol" db:"currency_symbol"`
	AmazonDomain   string `json:"amazon_domain" db:"amazon_domain"`
	AmazonRegion   string `json:"amazon_region" db:"amazon_region"`
	PartnerTag     string `json:"partner_tag" db:"partner_tag"`
}

// ─── Category ────────────────────────────────────────────────────────────────

type Category struct {
	ID          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name"`
	Slug        string  `json:"slug" db:"slug"`
	ParentID    *int    `json:"parent_id,omitempty" db:"parent_id"`
	ProductType string  `json:"product_type" db:"product_type"` // "storage" | "ram"
}

// ─── Product ─────────────────────────────────────────────────────────────────

type Product struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ASIN       string    `json:"asin" db:"asin"`
	Locale     string    `json:"locale" db:"locale"`
	Title      string    `json:"title" db:"title"`
	Brand      string    `json:"brand" db:"brand"`
	Model      string    `json:"model" db:"model"`
	CategoryID *int      `json:"category_id,omitempty" db:"category_id"`

	ImageURL  string `json:"image_url" db:"image_url"`
	AmazonURL string `json:"amazon_url" db:"amazon_url"`

	// Storage
	CapacityGB       *float64 `json:"capacity_gb,omitempty" db:"capacity_gb"`
	FormFactor       *string  `json:"form_factor,omitempty" db:"form_factor"`
	StorageInterface *string  `json:"storage_interface,omitempty" db:"storage_interface"`
	RPM              *int     `json:"rpm,omitempty" db:"rpm"`
	CacheMB          *int     `json:"cache_mb,omitempty" db:"cache_mb"`
	ReadSpeedMBPS    *int     `json:"read_speed_mbps,omitempty" db:"read_speed_mbps"`
	WriteSpeedMBPS   *int     `json:"write_speed_mbps,omitempty" db:"write_speed_mbps"`
	TBW              *int     `json:"tbw,omitempty" db:"tbw"`

	// RAM
	RAMType       *string  `json:"ram_type,omitempty" db:"ram_type"`
	RAMSpeedMHz   *int     `json:"ram_speed_mhz,omitempty" db:"ram_speed_mhz"`
	RAMCapacityGB *float64 `json:"ram_capacity_gb,omitempty" db:"ram_capacity_gb"`
	RAMModules    *int     `json:"ram_modules,omitempty" db:"ram_modules"`
	RAMFormFactor *string  `json:"ram_form_factor,omitempty" db:"ram_form_factor"`
	IsECC         bool     `json:"is_ecc" db:"is_ecc"`
	BufferType    *string  `json:"buffer_type,omitempty" db:"buffer_type"`
	CASLatency    *int     `json:"cas_latency,omitempty" db:"cas_latency"`
	Timing        *string  `json:"timing,omitempty" db:"timing"`
	Voltage       *float64 `json:"voltage,omitempty" db:"voltage"`

	LastFetchedAt *time.Time `json:"last_fetched_at,omitempty" db:"last_fetched_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// ─── ProductWithPrice ────────────────────────────────────────────────────────

// ProductWithPrice combines a product with its latest price for API responses.
type ProductWithPrice struct {
	Product
	Price        *float64 `json:"price,omitempty"`
	Currency     string   `json:"currency,omitempty"`
	Availability *string  `json:"availability,omitempty"`
	Condition    *string  `json:"condition,omitempty"`
	Merchant     *string  `json:"merchant,omitempty"`
	IsPrime      bool     `json:"is_prime"`
	PricePerGB   *float64 `json:"price_per_gb,omitempty"`
	PriceRecordedAt *time.Time `json:"price_recorded_at,omitempty"`

	Category *Category `json:"category,omitempty"`
}

// ─── Price History ────────────────────────────────────────────────────────────

type PriceHistory struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	ProductID    uuid.UUID  `json:"product_id" db:"product_id"`
	Price        float64    `json:"price" db:"price"`
	Currency     string     `json:"currency" db:"currency"`
	Availability string     `json:"availability" db:"availability"`
	Condition    string     `json:"condition" db:"condition"`
	Merchant     string     `json:"merchant" db:"merchant"`
	IsPrime      bool       `json:"is_prime" db:"is_prime"`
	RecordedAt   time.Time  `json:"recorded_at" db:"recorded_at"`
}

// ─── Price Alert ─────────────────────────────────────────────────────────────

type AlertType string

const (
	AlertTypeProductPrice         AlertType = "product_price"
	AlertTypeCategoryPricePerGB   AlertType = "category_price_per_gb"
	AlertTypeCategoryPricePerUnit AlertType = "category_price_per_unit"
)

type PriceAlert struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	AlertType  AlertType  `json:"alert_type" db:"alert_type"`

	ProductID       *uuid.UUID   `json:"product_id,omitempty" db:"product_id"`
	CategoryID      *int         `json:"category_id,omitempty" db:"category_id"`
	CategoryFilters *interface{} `json:"category_filters,omitempty" db:"category_filters"`

	ThresholdPrice      *float64 `json:"threshold_price,omitempty" db:"threshold_price"`
	ThresholdPricePerGB *float64 `json:"threshold_price_per_gb,omitempty" db:"threshold_price_per_gb"`

	IsActive        bool       `json:"is_active" db:"is_active"`
	NotifyEmail     bool       `json:"notify_email" db:"notify_email"`
	LastTriggeredAt *time.Time `json:"last_triggered_at,omitempty" db:"last_triggered_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

// ─── Request / Response DTOs ─────────────────────────────────────────────────

type ProductListResponse struct {
	Products []ProductWithPrice `json:"products"`
	Total    int                `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

type ProductFilters struct {
	Locale       string   `form:"locale" json:"locale"`
	ProductType  string   `form:"product_type" json:"product_type"` // "storage" | "ram"
	CategorySlug string   `form:"category" json:"category"`
	Search       string   `form:"q" json:"q"`
	MinGB        *float64 `form:"min_gb" json:"min_gb"`
	MaxGB        *float64 `form:"max_gb" json:"max_gb"`
	MinPrice     *float64 `form:"min_price" json:"min_price"`
	MaxPrice     *float64 `form:"max_price" json:"max_price"`
	FormFactor   string   `form:"form_factor" json:"form_factor"`
	Interface    string   `form:"interface" json:"interface"`

	// RAM filters
	RAMType     string `form:"ram_type" json:"ram_type"`
	RAMFormFactor string `form:"ram_form_factor" json:"ram_form_factor"`
	IsECC       *bool  `form:"is_ecc" json:"is_ecc"`
	BufferType  string `form:"buffer_type" json:"buffer_type"`
	RAMSpeed    *int   `form:"ram_speed" json:"ram_speed"`
	CASLatency  *int   `form:"cas_latency" json:"cas_latency"`

	SortBy   string `form:"sort_by" json:"sort_by"` // "price", "price_per_gb", "capacity", "price_asc"
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
}

type CreateAlertRequest struct {
	AlertType           AlertType `json:"alert_type" binding:"required"`
	ProductID           *string   `json:"product_id"`
	CategoryID          *int      `json:"category_id"`
	CategoryFilters     *map[string]string `json:"category_filters"`
	ThresholdPrice      *float64  `json:"threshold_price"`
	ThresholdPricePerGB *float64  `json:"threshold_price_per_gb"`
	NotifyEmail         bool      `json:"notify_email"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
