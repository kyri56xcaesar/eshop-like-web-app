// Package main is the runnable code /executable scripts.
// Every other package is a library.

// Go build -> creates an executable
// Go run -> creates an executable and runs it
// Go install -> creates an executable and installs it in the machine, making it global

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Print("Welcome.\n")

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("URL request: %s\n", r.URL)
		if r.Method == "GET" {
			t, err := template.ParseFiles("templates/index.html")
			if err != nil {
				log.Fatal(err)
			}
			t.Execute(w, nil)
		} else {
			respondWithError(w, 400, "Something went wrong")
		}
	})

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server starting on port %v", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
