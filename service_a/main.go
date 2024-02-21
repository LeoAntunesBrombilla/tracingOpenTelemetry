package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeoAntunesBrombilla/tracingOpenTelemetry/serviceA/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

var tracer = otel.Tracer("service_a")

func validateAndForwardCEPHandler(w http.ResponseWriter, r *http.Request) {
	_, validateCEPSpan := tracer.Start(context.Background(), "validateCEP")

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
	validateCEPSpan.End()

	ctx, span := otel.Tracer("service_a").Start(context.Background(), "callServiceB")
	defer span.End()

	serviceBURL := fmt.Sprintf("http://service_b:8081/processCEP?cep=%s", url.QueryEscape(cep))

	req, err := http.NewRequestWithContext(ctx, "GET", serviceBURL, nil)
	if err != nil {
		fmt.Println("Error passing the context", err)
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error requesting ", err)
		return
	}

	defer resp.Body.Close()
	if err != nil {
		http.Error(w, "Error forwarding request to Service B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func main() {
	tracing.InitTracer()
	http.HandleFunc("/validateCEP", validateAndForwardCEPHandler)

	port := "8080"
	log.Printf("Service A listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
