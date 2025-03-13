package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func main() {
	/*	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
		myslog := slog.New(jsonHandler)
		myslog.Info("hi there")
		myslog.Info("hello again", "key", "val", "age", 25)
	*/
	defaultProducts()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		oplog := httplog.LogEntry(r.Context())
		oplog.Info("info here")
		w.Write([]byte("Hello World!"))
	})
	r.Get("/product", handleProductGet)
	r.Post("/product", handleProductPost)
	// http.HandleFunc("/product", handleProduct)
	fmt.Println("Server running on port 8081")
	http.ListenAndServe(":8081", r)

	// DumpRequest not working
	const serverAddr = "http://localhost/8081"
	req, err := http.NewRequest(http.MethodGet, serverAddr, nil)
	req.Header.Add("test-header", "test-header-value")

	if err != nil {
		log.Fatal(err)
	}
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("REQUEST:\n%s", string(reqDump))
	// log.Fatal(http.ListenAndServe(":8081", nil))
}
