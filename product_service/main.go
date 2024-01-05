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
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Product struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Img      string  `json:"img"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Username string  `json:"username"`
}

var Db *sql.DB

func main() {

	time.Sleep(time.Second * 15)
	fmt.Print("Welcome to the Product Service.\n")

	godotenv.Load(".env")

	// Setup database connection
	// In Go database connection should be opened once.
	// host := os.Getenv("HOST")

	dbport := os.Getenv("DB_PORT")

	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", os.Getenv("HOST"), dbport),
		DBName:               os.Getenv("MYSQL_DATABASE"),
		AllowNativePasswords: true,
	}

	// Get database handle
	var err error
	Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	pingErr := Db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}

	fmt.Print("--> Successfully connected to the Database.\n\n")

	// Initialize database
	createTable := `
            CREATE TABLE IF NOT EXISTS products (
                id INT AUTO_INCREMENT PRIMARY KEY,
                title VARCHAR(255),
                img VARCHAR(1024),
                price DECIMAL(10, 2),
                quantity INT,
                username VARCHAR(255)
            );
        `
	rows, err := Db.Query(createTable)

	if err != nil {
		panic(err)
	}
	rows.Close()

	// Setup Kafka

	// Setup Service
	Router := chi.NewRouter()

	Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	Router.Get("/health", handlerReadiness)
	Router.Get("/err", handlerErr)
	Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	Router.Route("/products", func(Router chi.Router) {
		Router.Post("/", insertProduct)
		Router.Get("/", getProducts)

		Router.Route("/:{id}", func(r chi.Router) {
			//r.Use(ProductsCtx)
			r.Get("/", getProductsByID)
			r.Put("/", updateProductByID)
			r.Delete("/", deleteProductByID)
		})

		Router.Route("/name:{name}", func(r chi.Router) {
			r.Get("/", getProductsByNAME)
		})

		Router.Route("/username:{username}", func(r chi.Router) {
			r.Get("/", getProductsByUSERNAME)
		})
	})

	port := os.Getenv("PORT")

	srv := &http.Server{
		Handler: Router,
		Addr:    ":" + port,
	}

	// Setup Kafka Consumer
	time.Sleep(2 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	go Consume(ctx)
	defer cancel()

	// Launch Service
	log.Printf("Service starting on port %v", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		srv.Close()
	}

}

// func ProductsCtx(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
// 		productID := chi.URLParam(r, "id")

// 	})
// }

func getProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n------------------------------------\n")
	fmt.Println("Get request at /products")

	rows, err := Db.Query("SELECT * FROM products;")

	var products []Product

	for rows.Next() {
		pr := Product{}

		if err := rows.Scan(&pr.ID, &pr.Title, &pr.Img, &pr.Price, &pr.Quantity, &pr.Username); err != nil {
			respondWithError(w, 500, "Internal server Error.")
			rows.Close()
			return
		}
		products = append(products, pr)
		//fmt.Printf("%+v\n\n", pr)

	}

	rows.Close()

	//fmt.Println(products)
	//pp, _ := json.Marshal(products)
	//fmt.Println(string(pp))

	if err != nil {
		respondWithError(w, 500, "Internal server Error.")
		return
	}

	respondWithJSON(w, 200, products)
}

func insertProduct(w http.ResponseWriter, r *http.Request) {

	fmt.Print("\n\n------------------------------------\n")
	fmt.Println("Post request at /products")

	decoded, _ := io.ReadAll(r.Body)
	// fmt.Println(decoded)
	formatted, _ := formatJSON(decoded)
	// fmt.Println(formatJSON(decoded))

	fmt.Println(formatted)

	pr := Product{}
	err := json.Unmarshal([]byte(fmt.Sprintf(`%s`, formatted)), &pr)
	if err != nil {
		respondWithError(w, 404, "Resource not inserted")
		return
	}

	if pr.Title == "" || pr.Price < 0 || pr.Username == "" {
		respondWithError(w, 404, "Resource not inserted")
		return
	}

	fmt.Printf("%+v\n", pr)

	res, err := Db.Exec(`INSERT INTO products (title, img, price, quantity, username) VALUES (?, ?, ?, ?, ?)`, pr.Title, pr.Img, pr.Price, pr.Quantity, pr.Username)
	if err != nil {
		respondWithJSON(w, 404, "Insertion not successful.")
		return

	}
	last_id, _ := res.LastInsertId()
	respondWithJSON(w, 200, fmt.Sprintf("ID:%v Inserted", last_id))

}

func getProductsByID(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n------------------------------------\n")
	fmt.Println("Get request at /products/:id")

	id := chi.URLParam(r, "id")

	row := Db.QueryRow(`SELECT * FROM products WHERE id=?;`, id)

	pr := Product{}
	if err := row.Scan(&pr.ID, &pr.Title, &pr.Img, &pr.Price, &pr.Quantity, &pr.Username); err != nil {
		respondWithError(w, 404, "Resource not found")
		return
	}

	fmt.Printf("%+v\n", pr)

	respondWithJSON(w, 200, pr)

}

func updateProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n------------------------------------\n")
	fmt.Println("Put request at /products/:id")

	id := chi.URLParam(r, "id")

	decoded, _ := io.ReadAll(r.Body)
	// fmt.Println(decoded)
	formatted, _ := formatJSON(decoded)
	// fmt.Println(formatJSON(decoded))

	fmt.Println(formatted)

	pr := Product{}
	err := json.Unmarshal([]byte(fmt.Sprintf(`%s`, formatted)), &pr)
	if err != nil {
		http.Error(w, http.StatusText(405), 405)
	}

	fmt.Printf("%+v\n", pr)

	res, err := Db.Exec(`UPDATE products SET title=?, img=?, price=?, quantity=?, username=? WHERE id=?;`, pr.Title, pr.Img, pr.Price, pr.Quantity, pr.Username, id)

	if err != nil {
		respondWithJSON(w, 404, "Update not successful.")
		return

	}
	last_id, _ := res.RowsAffected()
	respondWithJSON(w, 200, fmt.Sprintf("Rows:%v Updated", last_id))

}
func deleteProductByID(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n------------------------------------\n")
	fmt.Println("Delete request at /products/:id")

	id := chi.URLParam(r, "id")

	res, err := Db.Exec(`DELETE FROM products WHERE id=?;`, id)

	if err != nil {
		respondWithJSON(w, 404, "Deletion not successful.")
		return

	}
	last_id, _ := res.RowsAffected()
	respondWithJSON(w, 200, fmt.Sprintf("Rows:%v Deleted", last_id))

}

func getProductsByNAME(w http.ResponseWriter, r *http.Request) {

	fmt.Print("\n\n------------------------------------\n")
	fmt.Println("Get request at /products/:name")

	name := chi.URLParam(r, "name")

	rows, err := Db.Query("SELECT * FROM products WHERE title=?;", name)

	var products []Product

	for rows.Next() {
		pr := Product{}

		if err := rows.Scan(&pr.ID, &pr.Title, &pr.Img, &pr.Price, &pr.Quantity, &pr.Username); err != nil {
			respondWithError(w, 500, "Internal server Error.")
			rows.Close()
			return
		}
		products = append(products, pr)
		//fmt.Printf("%+v\n\n", pr)

	}
	rows.Close()

	//fmt.Println(products)
	//pp, _ := json.Marshal(products)
	//fmt.Println(string(pp))

	if err != nil {
		respondWithError(w, 500, "Internal server Error.")
		return
	}

	respondWithJSON(w, 200, products)

}
func getProductsByUSERNAME(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n\n------------------------------------\n")
	fmt.Println("Get request at /products/:username")

	username := chi.URLParam(r, "username")

	rows, err := Db.Query("SELECT * FROM products WHERE username=?;", username)

	var products []Product

	for rows.Next() {
		pr := Product{}

		if err := rows.Scan(&pr.ID, &pr.Title, &pr.Img, &pr.Price, &pr.Quantity, &pr.Username); err != nil {
			respondWithError(w, 500, "Internal server Error.")
			rows.Close()
			return
		}
		products = append(products, pr)
		//fmt.Printf("%+v\n\n", pr)

	}

	rows.Close()

	//fmt.Println(products)
	//pp, _ := json.Marshal(products)
	//fmt.Println(string(pp))

	if err != nil {
		respondWithError(w, 500, "Internal server Error.")
		return
	}

	respondWithJSON(w, 200, products)
}
