package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/oauth2"
)

type authTemplate struct {
	Title   string
	Message string
}

var scopes = []string{"read", "write"}

const (
	redirectURL string = "http://localhost:8081/oauth/redirect"
	authURL     string = "https://www.inoreader.com/oauth2/auth?"
	tokenURL    string = "https://www.inoreader.com/oauth2/token"
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

func (c *config) getOauthConf() {

	oauthConf := &oauth2.Config{
		ClientID:     c.AppID,
		ClientSecret: c.AppKey,
		Scopes:       scopes,
		RedirectURL:  redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	c.OAuth2Conf = oauthConf
}

// HandleInoreaderLogin ---
func (c *config) handleInoreaderLogin(w http.ResponseWriter, r *http.Request) {

	oauthState := generateOauthStateCookie(w)
	url := c.OAuth2Conf.AuthCodeURL(oauthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// HandleInoreaderCallback ---
func (c *config) handleInoreaderCallback(w http.ResponseWriter, r *http.Request) {

	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("Invalid OAuth state")
		http.Redirect(w, r, "/go-inoreader", http.StatusTemporaryRedirect)
		return
	}

	token, err := c.exchangeToken(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/go-inoreader", http.StatusTemporaryRedirect)
	}

	c.writeCfgFile(getCfgFilePath(), token)

	http.Redirect(w, r, "/go-inoreader", http.StatusTemporaryRedirect)
}

func (c *config) exchangeToken(code string) (*oauth2.Token, error) {

	token, err := c.OAuth2Conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Code exchange was wrong: %s", err.Error())
	}

	return token, nil
}

func oauth2RestyClient(ctx context.Context) *resty.Client {

	c := loadConfig(getCfgFilePath())

	token := new(oauth2.Token)
	token.AccessToken = c.AccessToken
	token.RefreshToken = c.RefreshToken
	token.TokenType = c.TokenType
	token.Expiry = c.Expiry

	httpClient := c.OAuth2Conf.Client(ctx, token)
	return resty.NewWithClient(httpClient)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {

	authTempl := authTemplate{"Done", "You may close this page now and return to go-inoreader in terminal."}

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

// Init - Initiate OAuth flow
func Init() {

	ctx, cancel := context.WithCancel(context.Background())
	config := loadConfig(getCfgFilePath())
	config.getOauthConf()

	http.HandleFunc("/", config.handleInoreaderLogin)
	http.HandleFunc("/oauth/redirect", config.handleInoreaderCallback)
	http.HandleFunc("/go-inoreader", func(w http.ResponseWriter, r *http.Request) {
		serveTemplate(w, r)
		defer cancel()
	})

	srv := &http.Server{Addr: ":8081"}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Server error: %s", err.Error())
		}
	}()
	<-ctx.Done()
	if err := srv.Shutdown(ctx); err != nil && err != context.Canceled {
		log.Println(err)
	}
	log.Println("Done")
}
