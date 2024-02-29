package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofireflyio/aiac/v3/libaiac"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// ------------------------------------------------------------------------
func handleOptions(w http.ResponseWriter, r *http.Request) {
	// Respond to pre-flight request with CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Authorization")
	w.WriteHeader(http.StatusOK)
}

// ------------------------------------------------------------------------
func generateIAC(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Unmarshal the request body into a struct
	var requestData struct {
		Prompt string `json:"i_want"`
	}
	if err := json.Unmarshal(body, &requestData); err != nil {
		http.Error(w, "Failed to parse request JSON", http.StatusBadRequest)
		return
	}

	question := "generate " + requestData.Prompt
	log.Println(question)

	apiKey := os.Getenv("OPENAI_API_KEY")
	options := &libaiac.NewClientOptions{
		ApiKey: apiKey,
	}

	client := libaiac.NewClient(options)

	ctx := context.TODO()

	// Call the library to generate IAC code
	iacCode, err := client.GenerateCode(
		ctx,
		libaiac.ModelGPT35Turbo0301,
		string(question),
	)
	if err != nil {
		http.Error(w, "Failed to generate IAC code", http.StatusInternalServerError)
		return
	}

	// Respond with the generated code
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(iacCode.Code))
	// w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "text/plain")
	// responseData := struct {
	// 	IACCode string `json:"iac_code"`
	// }{
	// 	IACCode: iacCode.Code,
	// }
	// json.NewEncoder(w).Encode(responseData)

}

// ------------------------------------------------------------------------
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8086"
	}

	cors := handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "CORRELATIONID"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
	)

	router := mux.NewRouter()
	// Handle CORS pre-flight requests
	router.HandleFunc("/api/generate", handleOptions).Methods("OPTIONS")
	router.HandleFunc("/api/generate", generateIAC).Methods("POST")
	router.Use(cors)

	log.Println("Starting server at port ", port)

	log.Fatal(http.ListenAndServe(port, (router)))

}

// ------------------------------------------------------------------------
