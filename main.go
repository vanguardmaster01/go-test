package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Createdat   time.Time
	Updatedat   time.Time
}

type AllProducts struct {
	Products []*Product
}

var db *sql.DB
var err error

func main() {
	fmt.Println("hi")

	username := "root"
	password := ""
	host := "localhost"
	port := 3306
	database := "go"

	debugmode := "true"
	if debugmode == "true" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
		db, err = sql.Open("mysql", dsn)
	} else {
		// connectionName := mustGetenv("CLOUDSQL_CONNECTION_NAME")
		// user := mustGetenv("CLOUDSQL_USER")
		// password := os.Getenv("CLOUDSQL_PASSWORD") // NOTE: password may be empty
		// db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/GoJudge", user, password, connectionName))
	}

	// Maximum Idle Connections
	db.SetMaxIdleConns(5)
	// Maximum Open Connections
	db.SetMaxOpenConns(10)
	// Idle Connection Timeout
	db.SetConnMaxIdleTime(1 * time.Second)
	// Connection Lifetime
	db.SetConnMaxLifetime(30 * time.Second)

	if err != nil {
		panic((err.Error()))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	r := mux.NewRouter()
	r.HandleFunc("/login", loginPage)
	r.HandleFunc("/products", addProduct)
	r.HandleFunc("/products/{id}", handleProduct)
	r.HandleFunc("/products/{updateId}", updateProduct)
	r.HandleFunc("/products/{deleteId}", deleteProduct)
	r.HandleFunc("/", homePage)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.Handle("/", r)
	fmt.Println("Listening on 127.0.0.1:8080")
	err := http.ListenAndServe(":9000", nil) // setup listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	appengine.Main()
}

func handleProduct(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		fmt.Print("switch-getID")
		param := mux.Vars(req)["id"]

		var product Product

		query := "select id, name, description, price from products where id = ?"
		err = db.QueryRow(query, param).Scan(&product.ID, &product.Name, &product.Description,
			&product.Price)

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(product)

	case "POST":

		param := mux.Vars(req)["id"]

		if req.FormValue("_method") == "PUT" {
			fmt.Print("switch-put")
			fmt.Print("param =>", param)
			// id := req.FormValue("update_id")
			name := req.FormValue("update_name")
			description := req.FormValue("update_description")
			price := req.FormValue("update_price")
			query := "update products set name = ?, description = ?, price = ? where id = ?"

			_, err := db.Query(query, name, description, price, param)

			if err != nil {
				fmt.Println("update error", err)
			}

		} else if req.FormValue("_method") == "DELETE" {
			fmt.Print("switch-delete")
			// id := req.FormValue("delete_id")
			query := "delete from products where id = ?"

			_, err := db.Query(query, param)

			if err != nil {
				fmt.Println("delete error", err)
			}
		}
		http.Redirect(res, req, "/", http.StatusSeeOther)

	}

	// if req.Method == "GET" {
	// 	param := mux.Vars(req)["getId"]

	// 	var product Product

	// 	query := "select id, name, description, price from products where id = ?"
	// 	err = db.QueryRow(query, param).Scan(&product.ID, &product.Name, &product.Description,
	// 		&product.Price)

	// 	res.Header().Set("Content-Type", "application/json")
	// 	res.WriteHeader(http.StatusCreated)
	// 	json.NewEncoder(res).Encode(product)
	// }
}

func updateProduct(res http.ResponseWriter, req *http.Request) {
	fmt.Print("updaetID")
	if req.Method == "POST" {
		param := mux.Vars(req)["getId"]
		fmt.Print(param)
	}
}

func deleteProduct(res http.ResponseWriter, req *http.Request) {
	fmt.Print("deleteID")
	param := mux.Vars(req)["getId"]
	fmt.Print(param)
}

func addProduct(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		query := "select * from products"
		_, err := db.Query(query)
		if err != nil {
			log.Fatalf("impossible get products: %s", err)
		}

		fmt.Println("GET")
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() err: %v", err)
			return
		}

		name := req.FormValue("register_name")
		description := req.FormValue("register_description")
		price := req.FormValue("register_price")

		num, err := strconv.Atoi(price)
		if err != nil {
			fmt.Println("Conversion error: ", err)
			return
		}

		if name == "" || description == "" || num <= 0 {
			fmt.Fprintf(res, "ParseForm() err: %v", err)
			return
		}

		query := "INSERT INTO `products` (`name`, `description`, `price`, `createdat`, `updatedat`) VALUES (?, ?, ?, NOW(), NOW())"
		_, err = db.ExecContext(context.Background(), query, name, description, price)
		if err != nil {
			log.Fatalf("impossible insert products: %s", err)
		}

		fmt.Println("POST")
	default:
		fmt.Println("default")
	}

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func loginPage(res http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {

		message := "Enter username and password to login!"
		retry := req.URL.Query().Get("retry")
		checkRetry, _ := strconv.ParseBool(retry)
		varmap := map[string]interface{}{
			"Message": message,
			"Status":  "",
		}
		if checkRetry == true {
			message = "Invalid Username or Password!"
			varmap["Message"] = message
			varmap["Status"] = "error"
		}

		//http.ServeFile(res, req, "login.html")
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(res, varmap)
		return
	}
	//logging
	req.ParseForm()
	username := html.EscapeString(req.FormValue("username"))
	// password := html.EscapeString(req.FormValue("password"))
	log.Println(time.Now().Format(time.RFC850), "User Login Attempt by: ", username)
	var databaseUsername string
	// var databasePassword string

	// err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	// if err != nil {
	// 	http.Redirect(res, req, "/login?retry=1", 301)
	// 	return
	// }

	// err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	// if err != nil {
	// 	http.Redirect(res, req, "/login?retry=1", 301)
	// 	return
	// }

	res.Write([]byte("Hello " + databaseUsername))

}

func homePage(res http.ResponseWriter, req *http.Request) {
	query := "select id, name, description, price from products"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("impossible get products: %s", err)
	}

	var allProducts AllProducts
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description,
			&product.Price); err != nil {
			fmt.Print(err, "scan error")
		}

		allProducts.Products = append(allProducts.Products, &product)
	}
	if err = rows.Err(); err != nil {
		fmt.Print(err, "Err error")
	}

	// t, _ := template.ParseFiles("templates/index.html")
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Unable to parse from template:", err)
	}

	// s := Student{Name: "bug"}
	// t.Execute(res, s)

	t.Execute(res, &allProducts)
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}
