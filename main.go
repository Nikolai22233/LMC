package main

import (
	"LMC/internal/application"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", application.HandleCalculate)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
