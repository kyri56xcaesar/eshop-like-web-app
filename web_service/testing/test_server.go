package main

import (
	"fmt"
	"log"

	"net/http"
)

//var templates = template.Must(template.ParseFiles("./static/js", "./static/style"))

func serveFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	p := "." + r.URL.Path

	if p == "./" {
		p = "./static/index.html"
	}

	http.Handle("/", http.FileServer(http.Dir("./static")))

}

func main() {
	fmt.Println("Server starting on port 8080")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
