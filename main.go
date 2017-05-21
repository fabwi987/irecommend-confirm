package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fabwi987/irecommend-confirm/callback"
	"github.com/fabwi987/irecommend-confirm/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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

	r.HandleFunc("/callback", callback.CallbackHandler)
	r.HandleFunc("/confirm/{id}", LinkHandler)
	r.HandleFunc("/confirm/nonsocial/{id}", NosocialHandler)
	r.HandleFunc("/confirm/submit/{id}", SubmitHandler)

	http.ListenAndServe(":8081", handlers.LoggingHandler(os.Stdout, r))
}

func LinkHandler(w http.ResponseWriter, r *http.Request) {
	var client http.Client
	vars := mux.Vars(r)
	recommendationid := vars["id"]
	req, err := http.NewRequest("GET", os.Getenv("APP_HOST")+"recommendations/full/single/"+recommendationid, nil)
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
		IdReferral        uuid.UUID
		ReferralName      string
		ReferralPhone     string
		ReferralMail      string
	}{
		PositionTitle:     rec.Position.Title,
		PositionSubtitle:  rec.Position.Subtitle,
		PositionText:      rec.Position.Text,
		Idrecommendations: rec.Idrecommendations,
		IdReferral:        rec.Referral.Idreferrals,
		ReferralName:      rec.Referral.Name,
		ReferralPhone:     rec.Referral.Telephone,
		ReferralMail:      rec.Referral.Mail,
	}

	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("views/landing.html")
	t.Execute(w, layoutData)

}

func NosocialHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recommendationid := vars["id"]
	var client http.Client
	req, err := http.NewRequest("GET", os.Getenv("APP_HOST")+"recommendations/full/single/"+recommendationid, nil)
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	defer resp.Body.Close()
	var recommendation *models.Recommendation

	if err := json.NewDecoder(resp.Body).Decode(&recommendation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	refData := struct {
		Idrecommendations uuid.UUID
		IdReferral        uuid.UUID
		ReferralName      string
		ReferralPhone     string
		ReferralMail      string
		ReferralPicture   string
	}{
		Idrecommendations: recommendation.Idrecommendations,
		IdReferral:        recommendation.Referral.Idreferrals,
		ReferralName:      recommendation.Referral.Name,
		ReferralPhone:     recommendation.Referral.Telephone,
		ReferralMail:      recommendation.Referral.Mail,
		ReferralPicture:   recommendation.Referral.Picture,
	}

	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("views/confirmLinkedin.html")
	t.Execute(w, refData)

}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var confref models.Referral
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		decoder := schema.NewDecoder()
		err = decoder.Decode(&confref, r.PostForm)
	}

	var client http.Client
	vars := mux.Vars(r)
	recommendationid := vars["id"]
	req, err := http.NewRequest("GET", os.Getenv("APP_HOST")+"recommendations/full/single/"+recommendationid, nil)
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	defer resp.Body.Close()
	var rec *models.Recommendation

	if err := json.NewDecoder(resp.Body).Decode(&rec); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var updated bool
	if confref.Name != rec.Referral.Name {
		rec.Referral.Name = confref.Name
		updated = true
	}
	if confref.Telephone != rec.Referral.Telephone {
		rec.Referral.Telephone = confref.Telephone
		updated = true
	}
	if confref.Mail == rec.Referral.Mail {
		rec.Referral.Mail = confref.Mail
		updated = true
	}

	if updated {
		form := url.Values{}
		form.Add("Idreferrals", rec.Referral.Idreferrals.String())
		form.Add("Headline", rec.Referral.Headline)
		form.Add("Mail", rec.Referral.Mail)
		form.Add("Name", rec.Referral.Name)
		form.Add("Picture", rec.Referral.Picture)
		form.Add("ProfileURL", rec.Referral.ProfileURL)
		form.Add("ReferralUserID", rec.Referral.ReferralUserID)
		form.Add("Telephone", rec.Referral.Telephone)
		req, err = http.NewRequest("POST", os.Getenv("APP_HOST")+"referrals/single/"+rec.Referral.Idreferrals.String(), strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err = client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}

	req, err = http.NewRequest("POST", os.Getenv("APP_HOST")+"recommendations/single/"+recommendationid, nil)
	resp, err = client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("views/receipt.html")
	t.Execute(w, nil)

}
