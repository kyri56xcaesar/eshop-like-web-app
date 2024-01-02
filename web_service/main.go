// Package main is the runnable code /executable scripts.
// Every other package is a library.

// Go build -> creates an executable
// Go run -> creates an executable and runs it
// Go install -> creates an executable and installs it in the machine, making it global

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("File Server does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Print("Welcome to the Web Service.\n")

	//Load environment variables
	godotenv.Load(".env")

	// Get environment variables.
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	// Create a router
	router := chi.NewRouter()

	// Insert cors options.
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Setup FileServer for the static files.
	filesDir := http.Dir("./web/static")
	FileServer(router, "/static", filesDir)

	// Initialize Client
	createClient()

	// Setup routes.
	router.Get("/health", handlerReadiness)
	router.Get("/err", handlerErr)

	router.Get("/", HandleRootGet)
	router.Post("/", HandleRootPost)

	// Setup the server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// Launch Server.
	log.Printf("Server starting on port %v", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		srv.Close()
	}

}
