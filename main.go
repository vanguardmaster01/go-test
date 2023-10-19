package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/appengine"
)

var NoTooSmall = errors.New("the number is too small")

func ReturnPositive(no int) (int, error) {
	if no > 0 {
		return no, nil
	} else {
		return 0, NoTooSmall
	}
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

	if err != nil {
		panic((err.Error()))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/products", addProduct)
	http.HandleFunc("/", homePage)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	fmt.Println("Listening on 127.0.0.1:8080")
	err := http.ListenAndServe(":9000", nil) // setup listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	appengine.Main()
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

		name := req.FormValue("name")
		description := req.FormValue("description")
		price := req.FormValue("price")

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
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(res, nil)
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}
