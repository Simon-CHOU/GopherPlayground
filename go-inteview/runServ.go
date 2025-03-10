package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Secure World!")
}

// func main() {
// 	http.HandleFunc("/", helloHandler)
// 	// err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		fmt.Printf("Server failed to start: %v\n", err)
// 	}
// }

// % curl http://localhost:8080/
// Hello, Secure World!%
