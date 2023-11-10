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
	fmt.Print("Welcome.\n")

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
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Subrouter (for version control)
	v1Router := chi.NewRouter()

	// Setup FileServer for the static files.
	filesDir := http.Dir("./web/static")
	FileServer(v1Router, "/static", filesDir)

	// Setup routes.
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)

	v1Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\nMethod type: %v\n", r.Method)

		if r.Method == http.MethodGet {

			// Serve the HTML form
			http.ServeFile(w, r, "web/index.html")

		}

	})

	v1Router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\nMethod type: %v\n", r.Method)

		err := r.ParseForm()

		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Retrieve and process the form data
		// username := r.FormValue("username")
		// password := r.FormValue("password")
		// fmt.Printf("Username: %v\nPassword: %v\n\n\n", username, password)

		for key, value := range r.Form {
			fmt.Printf("%s = %s\n", key, value)
		}

		// Validate user credintials and redirect accordingly.
		// Perhaps use html/template to form the ServedFile

		// Send request to other service

		http.ServeFile(w, r, "web/index.html")
	})

	// Mount subrouter to the main router
	router.Mount("/v1", v1Router)

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
