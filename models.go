package main

import "sync"

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

var (
	products = make(map[string]Product)
	mux      sync.Mutex
)

func defaultProducts() {
	products = map[string]Product{
		"1": {
			ID:    "1",
			Name:  "Laptop",
			Price: "2000 USD",
		},
		"2": {
			ID:    "2",
			Name:  "Mobile",
			Price: "1000 USD",
		},
		"3": {
			ID:    "3",
			Name:  "Tablet",
			Price: "1500 USD",
		},
	}
}
