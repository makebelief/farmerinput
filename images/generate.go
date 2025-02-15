package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Product struct {
	Name        string
	Category    string
	Description string
	Price       float64
	Stock       int
	ImageURL    string
	Brand       string
	Unit        string
	Rating      float64
	Reviews     int
}

var (
	categories = []string{"Seeds", "Fertilizers", "Pesticides", "Tools"}
	brands     = []string{"Pioneer", "GrowMore", "EcoGuard", "FarmPro", "SeedCo"}
	units      = []string{"kg", "50kg", "5L", "set", "100g"}
)

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	fmt.Println("[")
	for i := 1; i <= 1000; i++ {
		product := Product{
			Name:        capitalizeEachWord(generateProductName()),
			Category:    capitalizeEachWord(categories[rand.Intn(len(categories))]),
			Description: capitalizeEachWord(generateDescription()),
			Price:       float64(rand.Intn(5000-100)+100) + rand.Float64(), // Random price between 100 and 5000
			Stock:       rand.Intn(1000) + 1,                              // Random stock between 1 and 1000
			ImageURL:    fmt.Sprintf("/images/products/product-%d.jpg", i),
			Brand:       capitalizeEachWord(brands[rand.Intn(len(brands))]),
			Unit:        capitalizeEachWord(units[rand.Intn(len(units))]),
			Rating:      1.0 + rand.Float64()*4.0, // Random rating between 1.0 and 5.0
			Reviews:     rand.Intn(500) + 1,       // Random reviews between 1 and 500
		}

		// Print the product in the desired format
		fmt.Printf(`  {
    Name:        "%s",
    Category:    "%s",
    Description: "%s",
    Price:       %.2f,
    Stock:       %d,
    ImageURL:    "%s",
    Brand:       "%s",
    Unit:        "%s",
    Rating:      %.1f,
    Reviews:     %d,
  }`, product.Name, product.Category, product.Description, product.Price, product.Stock, product.ImageURL, product.Brand, product.Unit, product.Rating, product.Reviews)

		if i < 1000 {
			fmt.Println(",")
		} else {
			fmt.Println()
		}
	}
	fmt.Println("]")
}

// capitalizeEachWord capitalizes the first letter of each word in a string
func capitalizeEachWord(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

func generateProductName() string {
	prefixes := []string{"Premium", "Organic", "High-Yield", "Eco-Friendly", "Advanced"}
	names := []string{"Maize Seeds", "Tomato Seeds", "Wheat Seeds", "NPK Fertilizer", "Organic Pesticide", "Garden Tools", "Potato Seeds", "Rice Seeds"}
	return fmt.Sprintf("%s %s", prefixes[rand.Intn(len(prefixes))], names[rand.Intn(len(names))])
}

func generateDescription() string {
	descriptions := []string{
		"High-yield hybrid seeds suitable for various climates.",
		"Balanced fertilizer for all-purpose use.",
		"Natural pest control solution safe for organic farming.",
		"Complete set of essential garden tools.",
		"High-yielding determinate variety for optimal growth.",
		"Eco-friendly solution for sustainable farming.",
		"Advanced formula for maximum crop yield.",
	}
	return descriptions[rand.Intn(len(descriptions))]
}