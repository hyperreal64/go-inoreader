package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"golang.org/x/oauth2"
)

type token struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

type authTemplate struct {
	Title   string
	Message string
}

var (
	oauthConf = oauth2.Config{
		ClientID:     os.Getenv("INOREADER_CLIENT_ID"),
		ClientSecret: os.Getenv("INOREADER_CLIENT_SECRET"),
		Scopes:       []string{"read", "write"},
		RedirectURL:  "http://localhost:8081/oauth/redirect",
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
func handleInoreaderLogin(w http.ResponseWriter, r *http.Request) {

	oauthState := generateOauthStateCookie(w)
	url := oauthConf.AuthCodeURL(oauthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleInoreaderCallback ---
func handleInoreaderCallback(w http.ResponseWriter, r *http.Request) {

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

	writeCfgFile(getCfgFilePath(), token)
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

func (cf *config) getOAuthResponse(ctx context.Context) *http.Client {

	token := new(oauth2.Token)
	token = cf.Token

	return oauthConf.Client(ctx, token)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {

	authTempl := authTemplate{"go-inoreader", "You may close this page now and return to go-inoreader in terminal."}

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
