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

func getPriceArray(tcode string) string {
	var result string = "["
	db, err := sql.Open("postgres", "user="+db_user+" dbname="+db_name+" password="+db_pass+" sslmode=disable host=localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT stock_id, target_date, open, high, low, close, volume FROM prices WHERE stock_id=$1 order by target_date", tcode)
	if err != nil {
		log.Fatal(err)
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

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	response := getPriceArray(r.URL.Path[1:])
	fmt.Fprint(w, response)
}

func main() {
	http.HandleFunc("/", viewHandler)
	http.ListenAndServe(":28080", nil)
}
