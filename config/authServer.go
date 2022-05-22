package config

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
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
	redirectURL          string = "http://localhost:8081/oauth/redirect"
	authURL              string = "https://www.inoreader.com/oauth2/auth?"
	tokenURL             string = "https://www.inoreader.com/oauth2/token"
	authResponseTemplate        = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Done</title>
</head>
<body>
<h1>{{ .Title }}</h1>
<hr>
{{ .Message }}
</body>
</html>
`
)

// Generates an Oauth state cookie and returns it as a string.
func generateOauthStateCookie(w http.ResponseWriter) (string, error) {

	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state, nil
}

// Loads go-inoreader.json configuration file into a oauth2.Config struct
// alongside Oauth config data.
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

// Sends initial request to the API to start Oauth flow using
// the api_key and app_id from go-inoreader.json. Handles user
// login in browser.
func (c *config) handleInoreaderLogin(w http.ResponseWriter, r *http.Request) {

	oauthState, err := generateOauthStateCookie(w)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	url := c.OAuth2Conf.AuthCodeURL(oauthState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Receives callback from Inoreader API, verifies Oauth state value, converts
// auth code into a token, and writes token info to go-inoreader.json. After
// successful exchange, user is redirected to callback URL, which is a page
// that displays the authResponseTemplate.
func (c *config) handleInoreaderCallback(w http.ResponseWriter, r *http.Request) {

	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("Invalid OAuth state")
		return
	}

	token, err := c.OAuth2Conf.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		log.Println(err)
	}

	if err := c.writeCfgFile(getCfgFilePath(), token); err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/go-inoreader", http.StatusTemporaryRedirect)
}

// Takes a context.Context, loads go-inoreader.json config with token info,
// initializes and returns a resty.Client with context and token data.
func Oauth2RestyClient(ctx context.Context) *resty.Client {

	c, err := loadConfig(getCfgFilePath())
	if err != nil {
		log.Println(err)
	}

	token := new(oauth2.Token)
	token.AccessToken = c.AccessToken
	token.RefreshToken = c.RefreshToken
	token.TokenType = c.TokenType
	token.Expiry = c.Expiry

	return resty.NewWithClient(c.OAuth2Conf.Client(ctx, token))
}

// Serves authResponseTemplate
func serveTemplate(w http.ResponseWriter, r *http.Request) {

	templateMsg := &authTemplate{Title: "Done", Message: "You may close this page and return to go-inoreader in your terminal."}
	var t = template.Must(template.New("authResponse").Parse(authResponseTemplate))
	if err := t.Execute(w, templateMsg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Initiates OAuth flow
func Init() {

	ctx, cancel := context.WithCancel(context.Background())
	config, err := loadConfig(getCfgFilePath())
	if err != nil {
		log.Fatalln(err)
	}
	config.getOauthConf()

	http.HandleFunc("/", config.handleInoreaderLogin)
	http.HandleFunc("/oauth/redirect", config.handleInoreaderCallback)
	http.HandleFunc("/go-inoreader", func(w http.ResponseWriter, r *http.Request) {
		serveTemplate(w, r)
		defer cancel()
	})

	// TODO: This can be improved
	srv := &http.Server{Addr: ":8081"}
	log.Println("Server listening on http://localhost:8081")
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Server error: %s", err.Error())
		}
	}()
	<-ctx.Done()

	log.Println("Done")
}
