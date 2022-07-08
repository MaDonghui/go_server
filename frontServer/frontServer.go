package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Printf("readying to send all the bad-designed websites :D\nRight over port 8080")
	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.ListenAndServe(":8080", nil)
}
