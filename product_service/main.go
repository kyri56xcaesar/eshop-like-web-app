package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
type Service struct {
	db *sql.DB
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := s.db
	switch r.URL.Path {
	default:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case "/healthz":
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		err := s.db.PingContext(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("db down: %v", err), http.StatusFailedDependency)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case "/quick-action":
		// This is a short SELECT. Use the request context as the base of
		// the context timeout.
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		product1 := Product{
			id:          1,
			title:       "test",
			img:         "test",
			price:       53.00,
			quantity:    5,
			seller_name: "JoNDo",
		}

		var product2 Product

		err := db.QueryRowContext(ctx, `
		SELECT 
			P.product 
		FROM 
			products P 
		WHERE 
			p.id = :id
		;`,
			sql.Named("id", product1.id)).Scan(&product2)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		io.WriteString(w, product2.seller_name)
		return

	case "/long-action":
		// This is a long SELECT. Use the request context as the base of
		// the context timeout, but give it some time to finish. If
		// the client cancels before the query is done the query will also
		// be canceled.
		ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
		defer cancel()

		var products []Product
		rows, err := db.QueryContext(ctx, "select * from products;")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var product_ Product
			err = rows.Scan(&product_)
			if err != nil {
				break
			}
			products = append(products, product_)
		}

		// Check for errors during rows "Close".
		// This may be more important if multiple statements are executed
		// in a single batch and rows were written as well as read.
		if closeErr := rows.Close(); closeErr != nil {
			http.Error(w, closeErr.Error(), http.StatusInternalServerError)
			return
		}

		// Check for row scan error.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check for errors during row iteration.
		if err = rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(products)
		return

	case "/async-action":
		// This action has side effects that we want to preserve
		// even if the client cancels the HTTP request part way through.
		// For this we do not use the http request context as a base for
		// the timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var orderRef = "AB123"
		tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		_, err = tx.ExecContext(ctx, "stored_proc_name", orderRef)

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tx.Commit()
		if err != nil {
			http.Error(w, "action in unknown state, check state before attempting again", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return

	}
}

func main() {
	fmt.Print("Welcome to the Product Service.\n")

	// Load environment variables using godotenv
	godotenv.Load(".env")

	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || user == "" || password == "" || dbname == "" || err != nil {
		log.Fatal("Environment variable not set.")
	}

	// Setup database connection
	psqlInfo := fmt.Sprintf("host=%s port=%v user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := &Service{db: db}

	log.Fatal(http.ListenAndServe(":8081", s))

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
