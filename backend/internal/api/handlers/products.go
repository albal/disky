package handlers

import (
	"context"
	"net/http"

	"github.com/albal/disky/backend/internal/database"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	DB *database.Pool
}

func New(db *database.Pool) *Handlers {
	return &Handlers{DB: db}
}

// GetProducts returns the list of products and their latest prices.
func (h *Handlers) GetProducts(c *gin.Context) {
	// Let's implement a simple direct SQL query to the database.
	// In reality we would use a more robust repository layer.
	query := `
		SELECT 
			p.id, p.asin, p.title, p.brand, p.model, p.capacity_gb, 
			p.form_factor, p.storage_interface, p.amazon_url, p.image_url,
			p.ram_type, p.ram_capacity_gb, p.ram_speed_mhz, p.is_ecc, p.cas_latency,
			lp.price, lp.currency, lp.condition
		FROM products p
		LEFT JOIN latest_prices lp ON lp.product_id = p.id
		ORDER BY lp.price ASC NULLS LAST
		LIMIT 100
	`

	rows, err := h.DB.Query(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer rows.Close()

	var products []map[string]interface{}
	for rows.Next() {
		var (
			id, asin, title             string
			brand, model                *string
			capacityGB                  *float64
			formFactor, storageIntf     *string
			amazonUrl, imageUrl         *string
			ramType                     *string
			ramCapGB                    *float64
			ramSpeed                    *int
			isECC                       *bool
			casL                        *int
			price                       *float64
			currency, condition         *string
		)

		err := rows.Scan(
			&id, &asin, &title, &brand, &model, &capacityGB,
			&formFactor, &storageIntf, &amazonUrl, &imageUrl,
			&ramType, &ramCapGB, &ramSpeed, &isECC, &casL,
			&price, &currency, &condition,
		)
		if err != nil {
			continue
		}

		p := map[string]interface{}{
			"id": id, "asin": asin, "title": title,
		}
		if brand != nil { p["brand"] = *brand }
		if capacityGB != nil { p["capacity_gb"] = *capacityGB }
		if formFactor != nil { p["form_factor"] = *formFactor }
		if storageIntf != nil { p["storage_interface"] = *storageIntf }
		if amazonUrl != nil { p["amazon_url"] = *amazonUrl }
		if imageUrl != nil { p["image_url"] = *imageUrl }
		if ramType != nil { p["ram_type"] = *ramType }
		if ramCapGB != nil { p["ram_capacity_gb"] = *ramCapGB }
		if ramSpeed != nil { p["ram_speed_mhz"] = *ramSpeed }
		if isECC != nil { p["is_ecc"] = *isECC }
		if casL != nil { p["cas_latency"] = *casL }
		if price != nil { p["price"] = *price }
		if currency != nil { p["currency"] = *currency }
		if condition != nil { p["condition"] = *condition }

		// Mock calculating price per GB
		if price != nil {
			if capacityGB != nil && *capacityGB > 0 {
				p["price_per_gb"] = *price / *capacityGB
			} else if ramCapGB != nil && *ramCapGB > 0 {
				p["price_per_gb"] = *price / *ramCapGB
			}
		}

		products = append(products, p)
	}

	// Always return some mocked data if the database is completely empty 
	// (for demo/development purposes)
	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{"products": getMockProducts()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func getMockProducts() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"id": "mock-1", "asin": "B089C51D2L", "title": "Samsung 870 QVO 8 TB SATA 2.5 Inch Internal Solid State Drive (SSD)",
			"capacity_gb": 8000.0, "form_factor": "2.5\"", "storage_interface": "SATA III", 
			"price": 319.99, "currency": "GBP", "condition": "New", "price_per_gb": 0.03999875,
			"amazon_url": "https://www.amazon.co.uk/dp/B089C51D2L?tag=prbox-21",
		},
		{
			"id": "mock-2", "asin": "B08QBJ2YMG", "title": "Crucial P2 1TB M.2 PCIe Gen3 NVMe Internal SSD",
			"capacity_gb": 1000.0, "form_factor": "M.2 2280", "storage_interface": "NVMe PCIe 3.0",
			"price": 39.99, "currency": "GBP", "condition": "New", "price_per_gb": 0.03999,
			"amazon_url": "https://www.amazon.co.uk/dp/B08QBJ2YMG?tag=prbox-21",
		},
		{
			"id": "mock-3", "asin": "B083TS2GG2", "title": "Corsair Vengeance LPX 32GB (2 x 16GB) DDR4 3200MHz",
			"ram_capacity_gb": 32.0, "ram_type": "DDR4", "ram_form_factor": "DIMM", "ram_speed_mhz": 3200,
			"price": 54.99, "currency": "GBP", "condition": "New", "price_per_gb": 1.71,
			"amazon_url": "https://www.amazon.co.uk/dp/B083TS2GG2?tag=prbox-21",
		},
	}
}

// Auth provider stubs
func (h *Handlers) GoogleLogin(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"url": "https://accounts.google.com/o/oauth2/v2/auth"}) }
func (h *Handlers) AppleLogin(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"url": "https://appleid.apple.com/auth/authorize"}) }
func (h *Handlers) MicrosoftLogin(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"url": "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"}) }
