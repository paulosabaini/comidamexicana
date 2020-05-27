package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondHTML(w http.ResponseWriter, msg string) {
	tmpl := template.Must(template.ParseFiles("template/message.html", "template/header.html", "template/footer.html"))
	tmpl.ExecuteTemplate(w, "message", msg)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
