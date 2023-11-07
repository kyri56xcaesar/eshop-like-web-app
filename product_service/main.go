package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pass123"
	dbname   = "products_db"
)

func main() {
	fmt.Print("Welcome to the Product Service.\n")

	// Load environment variables using godotenv
	godotenv.Load(".env")

	// Setup database connection
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)
	// db, err := sql.Open("postgres", psqlInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Print("--> Successfully connected to the Database!\n\n")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	fmt.Printf("--> Listening on PORT =: %v\n", portString)

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/hello", getHello)

	log.Fatal(http.ListenAndServe(":"+portString, mux))

}
