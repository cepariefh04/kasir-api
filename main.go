package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID int `json:"id"`
	Nama  string  `json:"nama"`
	Harga int `json:"harga"`
	Stok int `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Laptop", Harga: 15000000, Stok: 10},
	{ID: 2, Nama: "Smartphone", Harga: 5000000, Stok: 25},
	{ID: 3, Nama: "Tablet", Harga: 7000000, Stok: 15},
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
				id, err := strconv.Atoi(idStr)
				if err != nil {
					http.Error(w, "ID tidak valid", http.StatusBadRequest)
					return
				}

				for _, p := range produk {
					if p.ID == id {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(p)
						return
					}
				}

				http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

func updateProdukHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range produk {
		if produk[i].ID == id {
			produk[i].Nama = updateProduk.Nama
			produk[i].Harga = updateProduk.Harga
			produk[i].Stok = updateProduk.Stok
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk[i])
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

func deleteProdukHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}
	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}


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

func main() {
	// PRODUK SECTION

	//GET localhost:8080/api/produk - LIST PRODUK
	//POST localhost:8080/api/produk - CREATE PRODUK
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request){
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(produkBaru)
		}
	})

	// GET localhost:8080/api/produk/{id} - LIHAT DETAIL DAN UPDATE based on ID
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request){
		if r.Method == "GET" {
			getProdukByID(w, r)
		} else if r.Method == "PUT" {
			updateProdukHandler(w, r)
		} else if r.Method == "DELETE" {
			deleteProdukHandler(w, r)
		}
	})

	// KATEGORI SECTION
		
	//GET localhost:8080/api/categories - LIST PRODUK
	//POST localhost:8080/api/categories - CREATE PRODUK
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request){
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == "POST" {
			var newCategories Category
			err := json.NewDecoder(r.Body).Decode(&newCategories)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			newCategories.ID = len(categories) + 1
			categories = append(categories, newCategories)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(newCategories)
		}
	})

	// GET localhost:8080/api/categories/{id} - VIEW, UPDATE, DELETE KATEGORI based on ID
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request){
		if r.Method == "GET" {
			getCategoryByID(w, r)
		} else if r.Method == "PUT" {
			updateCategoryHandler(w, r)
		} else if r.Method == "DELETE" {
			deleteCategoryHandler(w, r)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
			"message": "Server is healthy",
		})
	}) // localhost:8000/health
	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Gagal menjalankan server:", err)
	}
}