package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tkanos/gonfig"
	"golang.org/x/oauth2"
)

type config struct {
	appID  string
	appSec string
	*oauth2.Token
}

// type appDevInfo struct {
// 	appID  string
// 	appSec string
// }

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
