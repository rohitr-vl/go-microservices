package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)
var tokenAuth *jwtauth.JWTAuth
func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil) // replace with secret key
  
	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)
  }
func main() {
	/*	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
		myslog := slog.New(jsonHandler)
		myslog.Info("hi there")
		myslog.Info("hello again", "key", "val", "age", 25)
	*/
	defaultProducts()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})
	// PUBLIC ROUTES
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// oplog not working
/*		oplog := httplog.LogEntry(r.Context())
		oplog.Info("info here")
*/
		w.Write([]byte("Hello World!"))
	})
	
	productRouter := chi.NewRouter()
	productRouter.Get("/", handleProductGet)
	// sub route
	productRouter.Get("/{id}", handleProductGet)
	
	// PROTECTED ROUTES
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))
	
		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		  _, claims, _ := jwtauth.FromContext(r.Context())
		  w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		})
		productRouter.Post("/", handleProductPost)
		productRouter.Put("/", handleProductUpdate)
		productRouter.Delete("/{id}", handleProductDelete)
	})
	// add sub route to main router
	r.Mount("/product", productRouter)
	// http.HandleFunc("/product", handleProduct)
	fmt.Println("Server running on port 8081")
	http.ListenAndServe(":8081", r)

	// DumpRequest not working
/*	const serverAddr = "http://localhost/8081"
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
*/
	// log.Fatal(http.ListenAndServe(":8081", nil))
}
