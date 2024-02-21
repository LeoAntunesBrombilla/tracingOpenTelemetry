package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func isValidCEP(cep string) bool {
	matched, _ := regexp.MatchString(`^\d{5}-?\d{3}$`, cep)
	return matched
}

func validateAndForwardCEPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		CEP string `json:"cep"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Could not decode request body", http.StatusBadRequest)
		return
	}

	if !isValidCEP(input.CEP) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	cep := input.CEP
	serviceBURL := fmt.Sprintf("http://service_b:8081/processCEP?cep=%s", url.QueryEscape(cep))
	fmt.Println(serviceBURL)

	resp, err := http.Get(serviceBURL)
	fmt.Println("LALALALLALA")
	if err != nil {
		fmt.Println("LALALALLALA")
		fmt.Println(err)
		http.Error(w, "Error forwarding request to Service B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/validateCEP", validateAndForwardCEPHandler)

	port := "8080"
	log.Printf("Service A listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
