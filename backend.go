package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strings"
)

type Price struct {
	Code   string
	Date   string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}

var db_user = os.Getenv("DBUSER")
var db_pass = os.Getenv("DBPASSWORD")
var db_name = os.Getenv("DBNAME")

func withHeader(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func withData(db *sql.DB, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SetVar(r, "db", db)
		f(w, r)
	}
}

func withVars(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		f(w, r)
	}
}

func getPriceArray(db *sql.DB, tcode string) string {
	var result string = "["

	rows, err := db.Query("SELECT stock_id, target_date, open, high, low, close, volume FROM prices WHERE stock_id=$1 order by target_date", tcode)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	defer rows.Close()

	for rows.Next() {
		var code string
		var tdate string
		var open float64
		var high float64
		var low float64
		var close float64
		var volume int
		rows.Scan(&code, &tdate, &open, &high, &low, &close, &volume)
		result += fmt.Sprintf("['%v', %v, %v, %v, %v],", tdate, low, close, open, high)
	}
	result = strings.TrimRight(result, ",")
	result += "]"
	return result
}

func getCodeArray(db *sql.DB) string {
	var result string = "["
	rows, err := db.Query("SELECT id, name FROM stocks")
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var name string
		rows.Scan(&id, &name)
		result += fmt.Sprintf("['%v', '%v'],", id, name)
	}
	result = strings.TrimRight(result, ",")
	result += "]"
	return result
}

func priceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("priceHandler() : start, r.URL.Path=", r.URL.Path)
	response := getPriceArray(GetVar(r, "db").(*sql.DB), r.URL.Path[len("/price/"):])
	fmt.Fprint(w, response)
	log.Println("priceHandler() : end")
}

func codeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("codeHandler() : start")
	response := getCodeArray(GetVar(r, "db").(*sql.DB))
	fmt.Fprint(w, response)
	log.Println("codeHandler() : end")
}

func main() {
	f, err := os.OpenFile("log/backend.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Fatal("error open file :", err.Error())
	}

	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	db, err := sql.Open("postgres", "user="+db_user+" dbname="+db_name+" password="+db_pass+" sslmode=disable host=localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/price/", withHeader(withVars(withData(db, priceHandler))))
	http.HandleFunc("/code/", withHeader(withVars(withData(db, codeHandler))))
	http.ListenAndServe(":28080", nil)
}
