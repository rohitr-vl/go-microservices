package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleProduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleProductGet(w, r)
	case "POST":
		handleProductPost(w, r)
	default:
		http.Error(w, "Method not supported.", http.StatusMethodNotAllowed)
	}
}

func handleProductPost(w http.ResponseWriter, r *http.Request) {
	var product Product
	// fmt.Printf("\nTypeof Request: %T & Request:%v\n", r.Body, r.Body)
	/*	log.Printf("\n s: %s \n", r.Body)
		log.Printf("\n v: %v \n", r.Body)
		log.Printf("\n +v: %+v \n", r.Body)
		log.Printf("\n #v: %#v \n", r.Body)
		log.Printf("\n T: %T \n", r.Body)
		log.Printf("\n t: %t \n", r.Body)
		log.Printf("\n q: %q \n", r.Body)
		log.Printf("\n ptr: %p \n", &r.Body)
		stringReq := fmt.Sprintf("\n sprintf: %s", r.Body)
		fmt.Println(stringReq)
		fmt.Fprintf(os.Stdin, "\n stdIn: %s \n", r.Body)
	*/
	/*	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
		myslog := slog.New(jsonHandler)
		myslog.Info(string(r.Body))
	*/
	err := json.NewDecoder(r.Body).Decode(&product)
	/*
		log.Printf("\n s: %s \n", err)
		log.Printf("\n v: %v \n", err)
		log.Printf("\n +v: %+v \n", err)
		log.Printf("\n #v: %#v \n", err)
		log.Printf("\n T: %T \n", err)
		log.Printf("\n t: %t \n", err)
		log.Printf("\n q: %q \n", err)
		log.Printf("\n ptr: %p \n", &err)
		jsonReq := fmt.Sprintf("\n sprintf: %s", err)
		fmt.Println(jsonReq)
		fmt.Fprintf(os.Stdin, "\n stdIn: %s \n", err)
	*/
	// fmt.Printf("\nTypeof err: %T & err:%v\n", err, err)
	// var reqJson map[string]interface{}
	// myslog.Info(json.Unmarshal(r.Body, reqJson))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ch := w.Header()
	ch.Set("Content-type", "application/json")
	mux.Lock()
	fmt.Printf("\nTypeof product: %T & product:%v\n", product, product)
	products[product.ID] = product
	mux.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func handleProductGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
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
