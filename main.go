// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	Unit        string  `json:"unit"`
	Rating      float64 `json:"rating"`
	Reviews     int     `json:"reviews"`
}

var db *gorm.DB
var store = sessions.NewCookieStore([]byte("your-secret-key")) // Replace with a secure key

func main() {
	// Open SQLite database
	var err error
	db, err = gorm.Open(sqlite.Open("farm_inputs.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Successfully connected to SQLite database")

	// Auto migrate the schema
	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	// Initialize router
	r := mux.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("."))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	r.PathPrefix("/styles.css").Handler(fs)
	r.PathPrefix("/app.js").Handler(fs)

	// Authentication routes
	r.HandleFunc("/", redirectToLogin).Methods("GET")
	r.HandleFunc("/login", serveLoginPage).Methods("GET")
	r.HandleFunc("/auth/login", handleLogin).Methods("POST")
	r.HandleFunc("/auth/google", handleGoogleAuth).Methods("GET")
	r.HandleFunc("/auth/google/callback", handleGoogleCallback).Methods("GET")
	
	// Protected routes (require authentication)
	r.HandleFunc("/home", authMiddleware(serveHomePage)).Methods("GET")
	r.HandleFunc("/api/products", authMiddleware(getProducts)).Methods("GET")
	r.HandleFunc("/api/products/category/{category}", authMiddleware(getProductsByCategory)).Methods("GET")
	r.HandleFunc("/api/products/{id}", authMiddleware(getProduct)).Methods("GET")
	r.HandleFunc("/api/products/search", authMiddleware(searchProducts)).Methods("GET")

	// Enable CORS
	r.Use(corsMiddleware)

	// Start server
	log.Println("Server starting on port http://localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func serveLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "login.html")
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Add your login logic here
	// For now, just redirect to home page
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = true
	session.Save(r, w)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func handleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	// Implement Google OAuth2 initialization
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Implement Google OAuth2 callback handling
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session-name")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// ... (rest of the existing handlers remain the same)
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
	result := db.Where("name LIKE ? OR description LIKE ?", "%"+query+"%", "%"+query+"%").Find(&products)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}