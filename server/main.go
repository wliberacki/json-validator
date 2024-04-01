package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

// Typ reprezentujący strukturę dokumentu JSON
type Policy struct {
	PolicyName     string `json:"PolicyName"`
	PolicyDocument struct {
		Version   string `json:"Version"`
		Statement []struct {
			Sid      string   `json:"Sid"`
			Effect   string   `json:"Effect"`
			Action   []string `json:"Action"`
			Resource string   `json:"Resource"`
		} `json:"Statement"`
	} `json:"PolicyDocument"`
}

// Funkcja walidująca pole Resource
func validateResource(resource string) bool {
	// Sprawdzenie pustego pola
	if resource == "" {
		return false
	}

	// Sprawdzenie pojedynczego gwiazdka
	if resource == "*" {
		return false
	}

	//sprawdzanie calego stringa
	matched, _ := regexp.MatchString(`[*]`, resource)
	return !matched
}

// Funkcja obsługi żądania HTTP
func verifyHandler(w http.ResponseWriter, r *http.Request) {
	// Odczytanie danych JSON z ciała żądania
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Błąd odczytu danych", http.StatusBadRequest)
		return
	}

	// Dekodowanie JSON do struktury Policy
	var policy Policy
	err = json.Unmarshal(body, &policy)
	if err != nil {
		http.Error(w, "Nieprawidłowy format JSON", http.StatusBadRequest)
		return
	}

	// Walidacja pola Resource dla każdego oświadczenia
	// Walidacja pola Resource dla każdego oświadczenia
	valid := false

	if len(policy.PolicyDocument.Statement) > 0 {
		valid = true
		for _, statement := range policy.PolicyDocument.Statement {
			if !validateResource(statement.Resource) {
				valid = false
				break // Zatrzymanie iteracji, jeśli pole Resource jest nieprawidłowe
			}
		}
	}

	// Zwrócenie odpowiedzi
	w.WriteHeader(http.StatusOK) // Always return 200 status code
	if valid {
		fmt.Fprintf(w, "true") // Pole Resource jest poprawne
	} else {
		fmt.Fprintf(w, "false")
	}
}

func main() {
	// Utworzenie routera HTTP
	router := mux.NewRouter()

	// Zdefiniowanie endpointu do weryfikacji
	router.HandleFunc("/api/verify", verifyHandler).Methods("POST")

	// Ustawienie obsługi CORS
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Uruchamianie serwera na porcie 8080
	http.ListenAndServe(":8080", router)
}
