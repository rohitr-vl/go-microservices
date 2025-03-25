// Handlers for the router
package main

import (
	"net/http"
)

var Counter = make(map[string]int)

type ProductHandler struct {
}

func (ph ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	_, ok := Counter["get"]
	if ok {
		Counter["get"] += 1
	} else {
		Counter["get"] = 1
	}
	mux.Unlock()
	handleProductGet(w, r)
}
func (ph ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	_, ok := Counter["get"]
	if ok {
		Counter["get"] += 1
	} else {
		Counter["get"] = 1
	}
	mux.Unlock()
	handleProductGet(w, r)
}
func (ph ProductHandler) CreateProducts(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	_, ok := Counter["post"]
	if ok {
		Counter["post"] += 1
	} else {
		Counter["post"] = 1
	}
	mux.Unlock()
	handleProductPost(w, r)
}
func (ph ProductHandler) UpdateProducts(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	_, ok := Counter["put"]
	if ok {
		Counter["put"] += 1
	} else {
		Counter["put"] = 1
	}
	mux.Unlock()
	handleProductUpdate(w, r)
}
func (ph ProductHandler) DeleteProducts(w http.ResponseWriter, r *http.Request) {
	mux.Lock()
	_, ok := Counter["delete"]
	if ok {
		Counter["delete"] += 1
	} else {
		Counter["delete"] = 1
	}
	mux.Unlock()
	handleProductDelete(w, r)
}
