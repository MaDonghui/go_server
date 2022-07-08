package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
	_ "github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
)

// Phone just a general format of JSON/Struct
type Phone struct {
	ID         int    `json:"id"`
	Brand      string `json:"brand"`
	Model      string `json:"model"`
	OS         string `json:"os"`
	Image      string `json:"image"`
	Screensize int    `json:"screensize"`
}

var database *sql.DB

func main() {
	printTime()
	color.Cyan("─=≡Σ((( つ＞＜)つ\nFiring up the backend server")

	database, _ = sql.Open("sqlite3", "./phoneInventory.db")
	if rowCounts("phones") == 0 {
		dbReset()
	}

	http.HandleFunc("/all", allItems)
	http.HandleFunc("/new", newItem)
	http.HandleFunc("/retrieve", retrieveItem)
	http.HandleFunc("/update", updateItem)
	http.HandleFunc("/delete", deleteItem)
	http.HandleFunc("/reset", resetDB)

	color.Green("Ready! Listening to port 8000")
	http.ListenAndServe(":8000", nil)
}

// localhost:8000/all
// GET, 400 200, 500 application/json
// get all items
func allItems(w http.ResponseWriter, r *http.Request) {
	// check if GET method
	if r.Method != "GET" {
		http.Error(w, "400 Bad Request: Please use GET method for /all\n", 400)
		printTime()
		color.Yellow("/all Bad Request\n")
		return
	}

	// select all items from database
	rows, err := database.Query("SELECT * FROM phones")
	if err != nil {
		http.Error(w, "500 InternalServerError, Please Contact the Admin\n", 500)
		printTime()
		color.Yellow("/all query error\n")
		return
	}

	var items []Phone
	// compose JSON
	for rows.Next() {
		var id int
		var brand string
		var model string
		var os string
		var image string
		var screensize int

		rows.Scan(&id, &brand, &model, &os, &image, &screensize)
		item := Phone{id, brand, model, os, image, screensize}
		items = append(items, item)
	}

	// prepare response head and initial body
	// CROS to all may cause security breach, if you actually planning to use this in production environment, modify this
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	itemsJSON, err := json.Marshal(items)
	//handle jsonify error
	if err != nil {
		http.Error(w, "500 InternalServerError, Please Contact the Admin\n", 500)
		printTime()
		color.Red("/all JSON-ify error\n")
		return
	}
	w.Write(itemsJSON)

	// all done :D
	printTime()
	color.Green("Successfully Responded to GET /all\n")
}

// localhost:8000/new
// GET, 400 200, 500 plain/text
// add new item(s) in database
func newItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "400 Bad Request: Please use POST method for /new\n", 400)
		printTime()
		color.Yellow("/new Bad Request\n")
		return
	}

	if r.Body == nil {
		http.Error(w, "400 Bad Request: Please send a request body for /new\n", 400)
		printTime()
		color.Yellow("/new no request body\n")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var item Phone
	err := decoder.Decode(&item)
	if err != nil {
		http.Error(w, "400 Bad Request: request body invalid or not properly formatted\n", 400)
		printTime()
		color.Yellow("/new decode request to json failed\n")
		return
	}

	color.Blue("\nAdding Following Item into database\n")
	color.White("brand: %s\nmodel: %s\nos: %s\nimage: %s\nscreensize: %d\n",
		item.Brand, item.Model, item.OS, item.Image, item.Screensize)
	_, err = database.Exec(`INSERT INTO phones (brand, model, os, image, screensize) 
		VALUES (?, ?, ?, ?, ?)`, item.Brand, item.Model, item.OS, item.Image, item.Screensize)
	if err != nil {
		http.Error(w, "500 InternalServerError, Please Contact the Admin\n", 500)
		printTime()
		color.Red("/new query failed\n")
		return
	}

	// prepare response head and initial body
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "plain/text; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully added new Item\n"))
	w.Write([]byte(fmt.Sprintf("brand: %s\nmodel: %s\nos: %s\nimage: %s\nscreensize: %d\n",
		item.Brand, item.Model, item.OS, item.Image, item.Screensize)))

	// All done :D
	printTime()
	color.Green("Successfully Responded to POST /new\n\n")
}

// localhost:8000/reset
// DELETE 400 200, plain/text
// reset database, this can not be undone (nothing is recoverable so far hahahah)
func resetDB(w http.ResponseWriter, r *http.Request) {
	//NOTE in chrome, it will actually send a preflag request "OPTION", I can't be bothered to make this 100% goodie so, lazy check
	if r.Method != "DELETE" && r.Method != "OPTIONS" {
		http.Error(w, "400 Bad Request: Please use DELETE method for /reset\n", 400)
		printTime()
		color.Yellow("/reset Bad Request\n")
		return
	}

	dbReset()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Content-Type", "plain/text; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database reset!"))
	// who am I o.o
	printTime()
	color.Green("Successfully Responded to DELETE /reset\n")
}

// localhost:8000/retrieve
// GET, 404, 400 200, 500 application/json
// get new item in database based on ID
func retrieveItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "400 Bad Request: Please use GET method for /retrieve\n", 400)
		printTime()
		color.Yellow("/retrieve Bad Request\n")
		return
	}

	// if id is not a number, 400 the dumb user
	paramater, err := strconv.ParseInt(r.FormValue("id")[0:], 10, 32)
	if err != nil {
		http.Error(w, "400 Bad Request: id may be invalid (required to be a integer)", 400)
		printTime()
		color.Yellow("/retrieve query parameter invalid\n")
		return
	}

	var rows int
	rows = rowCounts("phones")
	if 1 > paramater || paramater > int64(rows) {
		http.Error(w, "404 Resource Not Found, item attempting to retrieve does not exist\n", 404)
		printTime()
		color.Yellow("/retrieve id not in range (item doesn't exist)\n")
		return
	}

	var item Phone
	err = database.QueryRow("SELECT * FROM phones WHERE id= ?", paramater).Scan(&item.ID, &item.Brand, &item.Model, &item.OS, &item.Image, &item.Screensize)
	if err != nil {
		http.Error(w, "500 InternalServerError, Please Contact the Admin\n", 500)
		printTime()
		color.Red("/retrieve query error\n")
		return
	}

	color.Blue("\nSending Following Item to client\n")
	color.White("brand: %s\nmodel: %s\nos: %s\nimage: %s\nscreensize: %d\n",
		item.Brand, item.Model, item.OS, item.Image, item.Screensize)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	itemJSON, err := json.Marshal(item)
	//handle jsonify error
	if err != nil {
		http.Error(w, "500 InternalServerError, Please Contact the Admin\n", 500)
		printTime()
		color.Red("/all JSON-ify error\n")
		return
	}
	w.Write(itemJSON)
	// all done :D
	printTime()
	color.Green("Successfully Responded to GET /retrieve\n")
}

// localhost:8000/update
// GET, 404 400 200, 500 application/json
// update an item by ID
func updateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "400 Bad Request: Please use PUT method for /update\n", 400)
		printTime()
		color.Yellow("/update Bad Request\n")
		return
	}

	if r.Body == nil {
		http.Error(w, "400 Bad Request: Please send a request body for /update\n", 400)
		printTime()
		color.Yellow("/update no request body\n")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var item Phone
	err := decoder.Decode(&item)
	if err != nil {
		http.Error(w, "400 Bad Request: JSON format invalid\n", 400)
		printTime()
		color.Yellow("/update decode request to json failed\n")
		return
	}

	if 1 > item.ID || item.ID > rowCounts("phones") {
		http.Error(w, "404 Resource Not Found, item attempting to update does not exist\n", 404)
		printTime()
		color.Yellow("/update no such item found\n")
		return
	}

	color.Blue("\nUpdating Item\n")
	color.White("id: %d\nbrand: %s\nmodel: %s\nos: %s\nimage: %s\nscreensize: %d\n",
		item.ID, item.Brand, item.Model, item.OS, item.Image, item.Screensize)
	_, err = database.Exec("UPDATE phones SET brand = ?, model = ?, os = ?, image = ?, screensize = ? WHERE id = ?",
		item.Brand, item.Model, item.OS, item.Image, item.Screensize, item.ID)
	if err != nil {
		http.Error(w, "500 InternalServerError, Please Contact the Admin\n", 500)
		printTime()
		color.Red("/update update query failed\n")
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "plain/text; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully updated an Item\n"))
	w.Write([]byte(fmt.Sprintf("id: %d\nbrand: %s\nmodel: %s\nos: %s\nimage: %s\nscreensize: %d\n",
		item.ID, item.Brand, item.Model, item.OS, item.Image, item.Screensize)))

	// All done :D
	printTime()
	color.Green("Successfully Responded to PUT /update\n\n")
}

// localhost:8000/delete
// GET, 404 400 200, 500 plain/text
// delete an item by ID
func deleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "400 Bad Request: Please use DELETE method for /delete\n", 400)
		printTime()
		color.Yellow("/delete Bad Request\n")
		return
	}

	// if id is not a number, 400 the dumb user
	paramater, err := strconv.ParseInt(r.FormValue("id")[0:], 10, 32)
	if err != nil {
		http.Error(w, "400 Bad Request: id may be invalid (required to be a integer)", 400)
		printTime()
		color.Yellow("/delete query parameter invalid\n")
		return
	}

	var rows int
	rows = rowCounts("phones")
	if 1 > paramater || paramater > int64(rows) {
		http.Error(w, "404 Resource Not Found, item attempting to delete does not exist\n", 404)
		printTime()
		color.Yellow("/delete id not in range (item doesn't exist)\n")
		return
	}

	_, err = database.Exec("DELETE FROM phones WHERE id = ?", paramater)
	if err != nil {
		http.Error(w, "500 InternalServerError, Please Contact the Admin\n", 500)
		printTime()
		color.Red("/delete delete query failed\n")
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "plain/text; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully deleted\n"))

	// All done :D
	printTime()
	color.Green("Successfully Responded to DELETE /delete\n")
}

func dbReset() {
	color.Magenta("Resetting Database\n")

	database.Exec(`DROP TABLE IF EXISTS phones`)
	database.Exec(`CREATE TABLE phones (id INTEGER PRIMARY KEY, brand	CHAR(100) NOT NULL, model CHAR(100) NOT NULL, os CHAR(10) NOT NULL, image CHAR(254) NOT NULL, screensize INTEGER NOT NULL)`)
	database.Exec(`INSERT INTO phones (brand, model, os, image, screensize) VALUES ("Default", "Default", "Default", "example.com", 0)`)

	color.Magenta("Database Reset\n")
}

func rowCounts(table string) int {
	result, err := database.Query(fmt.Sprintf("SELECT COUNT(*) AS \"totalProducts\" FROM %s", table))
	if err != nil {
		color.Red("table \"phone\" doesn't exist, database resetting\n")
		dbReset()
		return 1
	}

	var count int
	for result.Next() {
		result.Scan(&count)
	}

	return count
}

func printTime() {
	var t = time.Now()
	fmt.Print(t.Format("2006-01-02 15:04:05\t"))
}
