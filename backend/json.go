package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	tmpls := []string{"./web/layout.html", "./web/" + tmpl + ".html"}
	t, err := template.ParseFiles(tmpls...)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing templates", err)
		return
	}

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error executing templates", err)
		return
	}
}
