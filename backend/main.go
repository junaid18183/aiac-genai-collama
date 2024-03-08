package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofireflyio/aiac/v4/libaiac"
	"github.com/gofireflyio/aiac/v4/libaiac/ollama"
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

// Send initial request to OLLAMA_API to keep alive.
func sendKeepAliveRequest() error {
	OLLAMA_API_BASE_URL := os.Getenv("OLLAMA_API_BASE_URL")
	url := OLLAMA_API_BASE_URL + "/generate"
	payload := map[string]interface{}{
		"model":      "codellama",
		"keep_alive": -1,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Optionally read the response if needed
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println("Initial Keep Alive Response: ", string(body))

	return nil
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

	// apiKey := os.Getenv("OPENAI_API_KEY")
	OLLAMA_API_BASE_URL := os.Getenv("OLLAMA_API_BASE_URL")
	options := &libaiac.NewClientOptions{
		// ApiKey: apiKey,
		OllamaURL: OLLAMA_API_BASE_URL,
		Backend:   libaiac.BackendOllama,
	}

	client := libaiac.NewClient(options)
	ctx := context.TODO()

	// Call the library to generate IAC code
	iacCode, err := client.GenerateCode(
		ctx,
		ollama.ModelCodeLlama,
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
	// Send an initial request to OLLAMA_API to keep alive before starting the server.
	err := sendKeepAliveRequest()
	if err != nil {
		log.Fatalf("Failed to send keep alive request: %v", err)
	}

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
