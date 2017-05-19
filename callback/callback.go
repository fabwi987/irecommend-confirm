package callback

import (
	_ "crypto/sha512"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

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

	/* Send to API for update

	recommendation.Referral.ReferralUserID = profile["user_id"].(string)
	recommendation.Referral.Name = profile["name"].(string)
	recommendation.Referral.Picture = profile["picture"].(string)
	recommendation.Referral.ProfileURL = profile["publicProfileUrl"].(string)
	recommendation.Referral.Headline = profile["headline"].(string)

	var refNew models.Referral
	refNew.Idreferrals = recommendation.Referral.Idreferrals
	refNew.Headline = recommendation.Referral.Headline
	refNew.Mail = recommendation.Referral.Mail
	refNew.Name = recommendation.Referral.Name
	refNew.Picture = recommendation.Referral.Picture
	refNew.ProfileURL = recommendation.Referral.ProfileURL
	refNew.ReferralUserID = recommendation.Referral.ReferralUserID
	refNew.Telephone = recommendation.Referral.Telephone

	//To be done by api
	_, err = models.CurrEnv.Db.UpdateReferral(refNew)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("views/confirmLinkedin.html")
	t.Execute(w, recommendation)*/

}
