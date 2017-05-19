package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/fabwi987/irecommend-confirm/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/confirm/{id}", LinkHandler)
	r.HandleFunc("/confirm/interested/{id}", InterestedHandler)

	http.ListenAndServe(":8081", handlers.LoggingHandler(os.Stdout, r))
}

func LinkHandler(w http.ResponseWriter, r *http.Request) {
	var client http.Client
	vars := mux.Vars(r)
	recommendationid := vars["id"]
	req, err := http.NewRequest("GET", "http://localhost:3000/recommendations/full/single/"+recommendationid, nil)
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	defer resp.Body.Close()
	var rec *models.Recommendation

	if err := json.NewDecoder(resp.Body).Decode(&rec); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	layoutData := struct {
		PositionTitle     string
		PositionSubtitle  string
		PositionText      string
		Idrecommendations uuid.UUID
	}{
		PositionTitle:     rec.Position.Title,
		PositionSubtitle:  rec.Position.Subtitle,
		PositionText:      rec.Position.Text,
		Idrecommendations: rec.Idrecommendations,
	}

	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("views/landing.html")
	t.Execute(w, layoutData)

}

func InterestedHandler(w http.ResponseWriter, r *http.Request) {
}
