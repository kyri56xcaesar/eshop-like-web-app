package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Order struct {
	ID          int       `json:"id"`
	Products    []Product `json:"products"`
	Total_price float64   `json:"total_price"`
	Status      string    `json:"status"`
	Username    string    `json:"username"`
}

type Product struct {
	ID       int     `json:"id"`
	Title    string  `json:"title"`
	Img      string  `json:"img"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Username string  `json:"username"`
}

var db *sql.DB

func main() {
	fmt.Print("Welcome to the Order Service.\n")

	// Load environment variables using godotenv
	godotenv.Load(".env")

	// Setup database connection
	// host := os.Getenv("HOST")
	dbport := os.Getenv("DB_PORT")

	cfg := mysql.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("localhost:%s", dbport),
		DBName:               os.Getenv("MYSQL_DATABASE"),
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Print("--> $usccesfully connected to the Database.\n\n")

	// Initialize Database

	createTable := `
		CREATE TABLE IF NOT EXISTS orders (
			id INT AUTO_INCREMENT PRIMARY KEY,
			products JSON,
			total_price DECIMAL(10, 2),
			status VARCHAR(100),
			username VARCHAR(255)
		);
        `
	rows, err := db.Query(createTable)

	if err != nil {
		panic(err)
	}
	rows.Close()

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

	Router.Route("/orders", func(Router chi.Router) {
		Router.Post("/", createOrder)

		Router.Route("/:{username}", func(r chi.Router) {
			//r.Use(ProductsCtx)
			r.Get("/", getOrdersByCUSTOMER)

		})

	})

	port := os.Getenv("PORT")

	srv := &http.Server{
		Handler: Router,
		Addr:    ":" + port,
	}

	// Launch Service
	log.Printf("Service starting on port %v", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		srv.Close()
	}

}

func createOrder(w http.ResponseWriter, r *http.Request) {

	fmt.Print("\n\n------------------------------------\n")
	fmt.Println("Post request at /orders")

	decoded, _ := io.ReadAll(r.Body)
	// fmt.Println(decoded)
	formatted, _ := formatJSON(decoded)
	// fmt.Println(formatJSON(decoded))

	//fmt.Println(formatted)

	ord := Order{}
	err := json.Unmarshal([]byte(fmt.Sprintf(`%s`, formatted)), &ord)
	if err != nil {
		respondWithError(w, 404, "Resource not inserted")
		return
	}
	if ord.Products == nil || ord.Username == "" {
		respondWithError(w, 404, "Resource not inserted")
		return
	}

	ord.Status = "pending"

	order_products, err := json.Marshal(ord.Products)
	if err != nil {
		respondWithError(w, 404, "Resource not inserted")
		return
	}

	res, err := db.Exec(`INSERT INTO orders (products, total_price, status, username) VALUES (?, ?, ?, ?)`, order_products, ord.Total_price, ord.Status, ord.Username)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Insertion not successful.\n%v", err))
		return

	}
	last_id, _ := res.LastInsertId()

	ord.ID = int(last_id)
	fmt.Printf("%+v\n\n", ord)

	// send to kafka
	var products_info []ProductsInfo

	for _, product := range ord.Products {
		var product_info ProductsInfo

		product_info.Product_id = product.ID
		product_info.Amount = product.Quantity

		products_info = append(products_info, product_info)
	}

	//fmt.Printf("Products info: %+v\n\n", products_info)

	kafka_msg := KafkaMsg{Order_id: int(last_id), Pr_info: products_info}

	//fmt.Printf("Kafka Msg: %+v\n\n", kafka_msg)

	jmsg, _ := json.Marshal(kafka_msg)
	msg, _ := formatJSON(jmsg)

	Produce(string(msg), 1)

	respondWithJSON(w, 200, fmt.Sprintf("ID:%v Inserted", last_id))

}

func getOrdersByCUSTOMER(w http.ResponseWriter, r *http.Request) {

	fmt.Print("\n\n------------------------------------\n")
	fmt.Println("Get request at /orders/:username")

	username := chi.URLParam(r, "username")

	rows, err := db.Query("SELECT * FROM orders WHERE username=?;", username)
	if err != nil {
		respondWithError(w, 404, "Something went wrong..")
		return
	}

	var orders []Order

	for rows.Next() {
		order := Order{}

		// Use temporary variables to store the JSON-encoded string
		var productsJSON string

		if err := rows.Scan(&order.ID, &productsJSON, &order.Total_price, &order.Status, &order.Username); err != nil {
			respondWithError(w, 500, fmt.Sprintf("Retrieval not successful.\n%v", err))
			rows.Close()
			return
		}

		// Unmarshal the JSON-encoded string into the Products field
		if err := json.Unmarshal([]byte(productsJSON), &order.Products); err != nil {
			respondWithError(w, 500, fmt.Sprintf("Failed to unmarshal Products field.\n%v", err))
			rows.Close()
			return
		}

		orders = append(orders, order)
	}

	rows.Close()

	//fmt.Println(products)
	//pp, _ := json.Marshal(products)
	//fmt.Println(string(pp))

	// if err != nil {
	// 	respondWithError(w, 500, "Internal server Error.")
	// 	return
	// }

	respondWithJSON(w, 200, orders)

}
