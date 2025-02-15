package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

// Product represents farm inputs and seeds
type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"imageUrl"`
	Brand       string  `json:"brand"`
	Unit        string  `json:"unit"` // e.g., "kg", "liters", "pieces"
	Rating      float64 `json:"rating"`
	Reviews     int     `json:"reviews"`
}

var db *gorm.DB

func main() {
	// Database connection
	dsn := "host=localhost user=postgres password=postgres dbname=farm_inputs port=5432"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&Product{})

	// Seed initial data if database is empty
	var count int64
	db.Model(&Product{}).Count(&count)
	if count == 0 {
		seedProducts()
	}

	// Initialize router
	r := mux.NewRouter()

	// Product routes
	r.HandleFunc("/api/products", getProducts).Methods("GET")
	r.HandleFunc("/api/products/category/{category}", getProductsByCategory).Methods("GET")
	r.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/products/search", searchProducts).Methods("GET")

	// Farmers Market routes
	r.HandleFunc("/api/farmers-markets", getFarmersMarkets).Methods("GET")

	// Enable CORS
	r.Use(corsMiddleware)

	// Start server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Add the missing getProduct handler
func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var product Product
	result := db.First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	var products []Product
	result := db.Find(&products)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func getProductsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	var products []Product
	result := db.Where("category = ?", category).Find(&products)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func searchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	var products []Product
	result := db.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&products)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func seedProducts() {
	products := []Product{
		{
			Name:        "Premium Maize Seeds",
			Category:    "Seeds",
			Description: "High-yield hybrid maize seeds suitable for various climates",
			Price:       799.99,
			Stock:       1000,
			ImageURL:    "/images/products/maize-seeds.jpg",
			Brand:       "Pioneer",
			Unit:        "kg",
			Rating:      4.8,
			Reviews:     156,
		},
		{
			Name:        "NPK Fertilizer",
			Category:    "Fertilizers",
			Description: "Balanced NPK 17-17-17 fertilizer for all-purpose use",
			Price:       2499.99,
			Stock:       500,
			ImageURL:    "/images/products/npk-fertilizer.jpg",
			Brand:       "GrowMore",
			Unit:        "50kg",
			Rating:      4.9,
			Reviews:     203,
		},
		{
			Name:        "Organic Pesticide",
			Category:    "Pesticides",
			Description: "Natural pest control solution safe for organic farming",
			Price:       999.99,
			Stock:       300,
			ImageURL:    "/images/products/organic-pesticide.jpg",
			Brand:       "EcoGuard",
			Unit:        "5L",
			Rating:      4.6,
			Reviews:     89,
		},
		{
			Name:        "Garden Hand Tools Set",
			Category:    "Tools",
			Description: "Complete set of essential garden tools",
			Price:       1499.99,
			Stock:       150,
			ImageURL:    "/images/products/garden-tools.jpg",
			Brand:       "FarmPro",
			Unit:        "set",
			Rating:      4.7,
			Reviews:     167,
		},
		{
			Name:        "Tomato Seeds",
			Category:    "Seeds",
			Description: "High-yielding determinate tomato variety",
			Price:       299.99,
			Stock:       800,
			ImageURL:    "/images/products/tomato-seeds.jpg",
			Brand:       "SeedCo",
			Unit:        "100g",
			Rating:      4.5,
			Reviews:     94,
		},
	}

	for _, product := range products {
		db.Create(&product)
	}
}

func getFarmersMarkets(w http.ResponseWriter, r *http.Request) {
	zipCode := r.URL.Query().Get("zip")
	radius := r.URL.Query().Get("radius")

	if zipCode == "" || radius == "" {
		http.Error(w, "Missing zip or radius parameter", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("USDA_API_KEY")
	if apiKey == "" {
		http.Error(w, "Missing USDA API key in environment variables", http.StatusInternalServerError)
		return
	}

	url := "https://search.ams.usda.gov/farmersmarkets/v1/data.svc/zipSearch?zip=" + zipCode + "&radius=" + radius + "&key=" + apiKey
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch data from USDA API", http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}