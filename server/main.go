package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

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

func validateResource(resource string) bool {
	if resource == "" {
		return false
	}

	if resource == "*" {
		return false
	}

	matched, _ := regexp.MatchString(`[*]`, resource)
	return !matched
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Data reading error", http.StatusBadRequest)
		return
	}

	var policy Policy
	err = json.Unmarshal(body, &policy)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	valid := false

	if len(policy.PolicyDocument.Statement) > 0 {
		valid = true
		for _, statement := range policy.PolicyDocument.Statement {
			if !validateResource(statement.Resource) {
				valid = false
				break
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	if valid {
		fmt.Fprintf(w, "true")
	} else {
		fmt.Fprintf(w, "false")
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/verify", verifyHandler).Methods("POST")

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
	http.ListenAndServe(":8080", router)
}
