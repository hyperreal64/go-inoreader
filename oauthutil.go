package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// Token ---
type Token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

// AuthTemplate ---
type AuthTemplate struct {
	Title   string
	Message string
}

var (
	oauthConf = oauth2.Config{
		ClientID:     os.Getenv("INOREADER_CLIENT_ID"),
		ClientSecret: os.Getenv("INOREADER_CLIENT_SECRET"),
		Scopes:       []string{"read", "write"},
		RedirectURL:  "http://localhost:8080/oauth/redirect",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.inoreader.com/oauth2/auth?",
			TokenURL: "https://www.inoreader.com/oauth2/token",
		},
	}
)

func generateOauthStateCookie(w http.ResponseWriter) string {

	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

// HandleInoreaderLogin ---
func HandleInoreaderLogin(w http.ResponseWriter, r *http.Request) {

	oauthState := generateOauthStateCookie(w)
	url := oauthConf.AuthCodeURL(oauthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleInoreaderCallback ---
func HandleInoreaderCallback(w http.ResponseWriter, r *http.Request) {

	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("Invalid OAuth state")
		http.Redirect(w, r, "/go-inoreader", http.StatusTemporaryRedirect)
		return
	}

	token, err := exchangeToken(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/go-inoreader", http.StatusTemporaryRedirect)
	}

	data, err := getUserDataFromInoreader(token)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/go-inoreader", http.StatusTemporaryRedirect)
		return
	}

	filePath := getCfgFilePath()
	cf, err := newCfgFile(filePath, data)
	cf.Contents.Oauth2Response = token
	cf.writeCfgFile()
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/go-inoreader", http.StatusTemporaryRedirect)
		return
	}

	http.Redirect(w, r, "/go-inoreader", http.StatusTemporaryRedirect)
}

func exchangeToken(code string) (*oauth2.Token, error) {

	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Code exchange was wrong: %s", err.Error())
	}

	return token, nil
}

// GetOAuthResponse ---
func (cf *Configuration) GetOAuthResponse(ctx context.Context) *http.Client {

	token := new(oauth2.Token)
	token = cf.Oauth2Response

	return oauthConf.Client(ctx, token)
}

func getUserDataFromInoreader(token *oauth2.Token) ([]byte, error) {

	url := "https://www.inoreader.com/reader/api/0/user-info"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize GET request")
	}

	bearerToken := fmt.Sprintf("Bearer %s", token.AccessToken)
	req.Header.Add("Authorization", bearerToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "HTTP GET request failed")
	}

	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read response body")
	}

	return contents, nil
}

// ServeTemplate ---
func ServeTemplate(w http.ResponseWriter, r *http.Request) {

	authTempl := AuthTemplate{"go-inoreader", "You may close this page now and return to go-inoreader in terminal."}

	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, authTempl); err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
