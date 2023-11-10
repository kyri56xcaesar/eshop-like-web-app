package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Product struct {
	id          int
	title       string
	img         string
	price       float64
	quantity    int
	seller_name string
}

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

	// _, err = db.Query("INSERT INTO products VALUES (1, 'a title', 'img_link', 74.3, 5, 'JonDo');")
	// if err != nil {
	// 	panic(err)
	// }

	rows, err := db.Query("SELECT * FROM products;")
	if err != nil {
		panic(err)
	}

	results := []Product{}

	for rows.Next() {
		product := Product{}
		err = rows.Scan(&product.id, &product.title, &product.img, &product.price, &product.quantity, &product.seller_name)
		if err != nil {
			log.Print(err)
		}

		results = append(results, product)

	}

	fmt.Println(results)

	rows.Close()

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
