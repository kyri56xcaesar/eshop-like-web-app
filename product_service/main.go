package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Print("Welcome to the Product Service.\n")

	// Load environment variables using godotenv
	godotenv.Load(".env")

	// Setup database connection
	// host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	//connectionStr := fmt.Sprintf("postgres://postgres:%v@%v:%d/%v?sslmode=disable", user, host, port)
	connectionStr := fmt.Sprintf("user=%v password=%v dbname=%v port=%d sslmode=disable",
		user, password, dbname, port)

	db, err := sql.Open("postgres", connectionStr)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Print("--> Successfully connected to the Database!\n\n")

	// Setup service

	// portString := os.Getenv("PORT")
	// if portString == "" {
	// 	log.Fatal("PORT is not found in the environment")
	// }

	// fmt.Printf("--> Listening on PORT =: %v\n", portString)

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", getRoot)
	// mux.HandleFunc("/hello", getHello)

	// log.Fatal(http.ListenAndServe(":"+portString, mux))

}
