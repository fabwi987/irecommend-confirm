package callback

import (
	_ "crypto/sha512"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fabwi987/irecommend-confirm/models"
	uuid "github.com/satori/go.uuid"

	"log"

	"golang.org/x/oauth2"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	domain := os.Getenv("AUTH0_DOMAIN")

	conf := &oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Scopes:       []string{"state", "openid", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://" + domain + "/authorize",
			TokenURL: "https://" + domain + "/oauth/token",
		},
	}

	code := r.URL.Query().Get("code")

	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Getting now the userInfo
	client := conf.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://" + domain + "/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var profile map[string]interface{}
	if err = json.Unmarshal(raw, &profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var APIclient http.Client
	req, err := http.NewRequest("GET", os.Getenv("APP_HOST")+"recommendations/full/single/"+r.URL.Query().Get("state"), nil)
	resp, err = APIclient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	defer resp.Body.Close()
	var recommendation *models.Recommendation

	if err := json.NewDecoder(resp.Body).Decode(&recommendation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//To be done by api
	form := url.Values{}
	form.Add("Idreferrals", recommendation.Referral.Idreferrals.String())
	form.Add("Headline", profile["headline"].(string))
	form.Add("Mail", recommendation.Referral.Mail)
	form.Add("Name", profile["name"].(string))
	form.Add("Picture", profile["picture"].(string))
	form.Add("ProfileURL", profile["publicProfileUrl"].(string))
	form.Add("ReferralUserID", profile["user_id"].(string))
	form.Add("Telephone", recommendation.Referral.Telephone)

	req, err = http.NewRequest("POST", os.Getenv("APP_HOST")+"referrals/single/"+recommendation.Referral.Idreferrals.String(), strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err = APIclient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Println("Thank you")

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
		ReferralName:      profile["name"].(string),
		ReferralPhone:     recommendation.Referral.Telephone,
		ReferralMail:      recommendation.Referral.Mail,
		ReferralPicture:   profile["picture"].(string),
	}

	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("views/confirmLinkedin.html")
	t.Execute(w, refData)

}
