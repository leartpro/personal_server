package main

import (
	"encoding/json"
	"fmt"
	"github.com/rs/cors"
	"io/ioutil"
	"net/http"
)

type Config struct {
	ServerAddr string `json:"server_addr"`
	ServerPort string `json:"server_port"`
}

func main() {
	// read the config file
	configFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	// parse the config file
	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	// Create a new mux router
	mux := http.NewServeMux()

	// Register the contact handler
	mux.HandleFunc("/contact", handleContact)

	// Use cors middleware to allow cross-origin requests
	corsHandler := cors.Default().Handler(mux)

	// Start the server
	addr := fmt.Sprintf("%s:%s", config.ServerAddr, config.ServerPort)
	fmt.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, corsHandler); err != nil {
		panic(err)
	}
}

type Contact struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request found")
	if r.Method == "POST" {
		// Die Größe der Daten begrenzen
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		// JSON-Decoder erstellen
		decoder := json.NewDecoder(r.Body)

		// JSON in Go-Objekt konvertieren
		var contact Contact
		err := decoder.Decode(&contact)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Die Daten auf der Konsole ausgeben
		fmt.Println("Name:", contact.Name)
		fmt.Println("Email:", contact.Email)
		fmt.Println("Message:", contact.Message)
	} else {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
	}
}
