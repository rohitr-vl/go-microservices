package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getApiCallCount(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Counter);
}

func handleProductPost(w http.ResponseWriter, r *http.Request) {
	var product []Product
	fmt.Printf("\nTypeof Request: %T & Request:%v\n", r.Body, r.Body)
	err := json.NewDecoder(r.Body).Decode(&product)
	fmt.Printf("\nTypeof err: %T & err:%v\n", err, err)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ch := w.Header()
	ch.Set("Content-type", "application/json")
	mux.Lock()
	for _, p := range product {
		fmt.Printf("\nTypeof product: %T & product:%v\n", p, p)
		products[p.ID] = p
	}
	mux.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func handleProductGet(w http.ResponseWriter, r *http.Request) {
	// id := r.URL.Query().Get("id")
	id := chi.URLParam(r, "id")
	if id != "" {
		mux.Lock()
		product, exists := products[id]
		mux.Unlock()

		if !exists {
			http.Error(w, "Product not found.", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(product)
		return
	}

	var productList []Product
	ch := w.Header()
	ch.Set("Content-type", "application/json")
	mux.Lock()
	for _, product := range products {
		productList = append(productList, product)
	}
	mux.Unlock()

	json.NewEncoder(w).Encode(productList)
}

func handleProductUpdate(w http.ResponseWriter, r *http.Request) {
	var product []Product
	fmt.Printf("\nTypeof Request: %T & Request:%v\n", r.Body, r.Body)
	err := json.NewDecoder(r.Body).Decode(&product)
	fmt.Printf("\nTypeof err: %T & err:%v\n", err, err)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if product with given id exists
	id := product[0].ID
	fmt.Printf("\nID type: %T & ID:%v\n", id, id)
	if id != "" {
		mux.Lock()
		productOld, exists := products[id]
		mux.Unlock()
		if !exists {
			http.Error(w, "Product not found.", http.StatusNotFound)
			return
		} else {
			fmt.Printf("\nTypeof Product: %T & Product:%v\n", productOld, productOld)
		}
	}
	ch := w.Header()
	ch.Set("Content-type", "application/json")
	mux.Lock()
	for _, p := range product {
		fmt.Printf("\nTypeof product: %T & product:%v\n", p, p)
		products[p.ID] = p
	}
	mux.Unlock()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func handleProductDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Product ID is required.", http.StatusBadRequest)
		return
	} else {
		mux.Lock()
		product, exists := products[id]
		mux.Unlock()
		if !exists {
			http.Error(w, "Product not found."+product.Name, http.StatusNotFound)
			return
		}
	}

	mux.Lock()
	delete(products, id)
	mux.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func delete(products map[string]Product, id string) []Product {
	var newProducts []Product
	for _, product := range products {
		if product.ID != id {
			newProducts = append(newProducts, product)
		}
	}
	return newProducts

}