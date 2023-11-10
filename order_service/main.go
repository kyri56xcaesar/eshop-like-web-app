package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Order struct {
	id            int
	product_title string
	quantity      int
	buyer_name    string
	seller_name   string
}

func main() {
	fmt.Print("Welcome to the Orders Service.\n")

	// Load environment variables using godotenv
	godotenv.Load(".env")

	// Setup database connection
	// host := os.Getenv("HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	connectionStr := fmt.Sprintf("user=%v password=%v dbname=%v port=%d sslmode=false",
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

	fmt.Print("--> $usccesfully connected to the Database.\n\n")
}
