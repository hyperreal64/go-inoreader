package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tkanos/gonfig"
	"golang.org/x/oauth2"
)

type authTemplate struct {
	Title   string
	Message string
}

var (
	oauthConf = oauth2.Config{
		ClientID:     os.Getenv("INO_APP_ID"),
		ClientSecret: os.Getenv("INO_APP_SEC"),
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

	writeCfgFile(token)
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

type config struct {
	appID  string
	appSec string
	*oauth2.Token
}

// Init --- Initiate Oauth flow
func Init() {

	ctx, cancel := context.WithCancel(context.Background())
	http.HandleFunc("/", handleInoreaderLogin)
	http.HandleFunc("/oauth/redirect", handleInoreaderCallback)
	http.HandleFunc("/go-inoreader", func(w http.ResponseWriter, r *http.Request) {
		serveTemplate(w, r)
		cancel()
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

// Get the path of the configuration file
// On Unix/Linux: $XDG_DATA_HOME/go-inoreader.json
// On Windows: %APPDATA%\go-inoreader.json
func getCfgFilePath() string {

	homeDir, _ := os.UserHomeDir()
	var fileName string

	switch runtime.GOOS {
	case "windows":
		fileName = path.Join(homeDir, "/AppData/Local/go-inoreader.json")

	default:
		fileName = path.Join(homeDir, "/.local/share/go-inoreader.json")
	}

	return fileName
}

func writeCfgFile(oauth2Resp *oauth2.Token) error {

	jsonData, err := json.MarshalIndent(&oauth2Resp, "", "    ")
	if err != nil {
		return errors.Wrap(err, "Could not marshal JSON data to file")
	}

	filepath := getCfgFilePath()
	if err = ioutil.WriteFile(filepath, jsonData, 0644); err != nil {
		return errors.Wrapf(err, "Could not write to config file: %s\n", filepath)
	}

	return nil
}

func config2Client() (*resty.Client, error) {
	cf := &config{}
	if err := gonfig.GetConf(getCfgFilePath(), cf); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := cf.getOAuthResponse(ctx)

	return resty.NewWithClient(client), nil
}
