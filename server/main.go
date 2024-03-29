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

// Funkcja do walidacji pola Resource
func validateResource(resource string) bool {
	// Sprawdzenie pojedynczego gwiazdka
	if resource == "*" {
		return false
	}

	// Dodatkowe walidacje (opcjonalne)
	// Przykład: sprawdzanie innych znaków specjalnych za pomocą regexu
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9:/\-_.]+$`, resource)
	return matched
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
	valid := true

	for _, statement := range policy.PolicyDocument.Statement {
		if !validateResource(statement.Resource) {
			valid = false
			break // Stop iterating if any Resource is invalid
		}
	}

	// Zwrócenie odpowiedzi "OK"
	if valid {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "true") // Resource is valid (no single asterisk)
	} else {
		http.Error(w, "Pole Resource zawiera niedozwolone znaki", http.StatusBadRequest)
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
	http.ListenAndServe(":3000", router)
}
