package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Category struct {
	ID int `json:"id"`
	Name  string  `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Desktop", Description: "Komputer desktop untuk kebutuhan sehari-hari dan gaming."},
	{ID: 2, Name: "Mobile", Description: "Perangkat mobile seperti smartphone dan tablet."},
}


func getCategoryByID(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
				id, err := strconv.Atoi(idStr)
				if err != nil {
					http.Error(w, "ID tidak valid", http.StatusBadRequest)
					return
				}

				for _, c := range categories {
					if c.ID == id {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(c)
						return
					}
				}

				http.Error(w, "Category not found!", http.StatusNotFound)
}

func updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateCategory Category
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			categories[i].Name = updateCategory.Name
			categories[i].Description = updateCategory.Description
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories[i])
			return
		}
	}

	http.Error(w, "Category not found!", http.StatusNotFound)
}

func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Delete Successfully",
			})
			return
		}
	}
	http.Error(w, "Category not found!", http.StatusNotFound)
}

type Config struct {
	Port    string `mapstructure:"PORT"`
	DBConn  string `mapstructure:"DB_CONN"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}


	config := Config{
		Port: viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup Database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup Routes
	http.HandleFunc("/api/products", productHandler.HandleProducts) // localhost:8080/api/products
	http.HandleFunc("/api/products/", productHandler.HandleProductsById) // localhost:8080/api/products/{id}

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"message": "Server is healthy",
		})
	}) // localhost:8000/health
	fmt.Println("Server running di localhost:"+ config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
		if err != nil {
			fmt.Println("Gagal menjalankan server:", err)
	}
}