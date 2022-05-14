package config

import (
	"encoding/json"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type config struct {
	AppID        string    `json:"app_id"`
	AppKey       string    `json:"app_key"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type,omitempty"`
	Expiry       time.Time `json:"expiry,omitempty"`
	OAuth2Conf   *oauth2.Config
}

func loadConfig(filePath string) (cfg *config, err error) {

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "Config file does not exist: %s", filePath)
	}

	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "An error occurred while trying to read config file: %s", filePath)
	}

	if err := json.Unmarshal(configFile, &cfg); err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal JSON content: %s", filePath)
	}

	validateConfig(cfg)

	return cfg, nil
}

func validateConfig(c *config) error {

	var fieldsMissing []string
	if c.AppID == "" {
		fieldsMissing = append(fieldsMissing, "app_id")
	}

	if c.AppKey == "" {
		fieldsMissing = append(fieldsMissing, "app_key")
	}

	if len(fieldsMissing) > 0 {
		errMsg := "The following fields appear missing from config:"
		return errors.Wrap(errors.New(errMsg), strings.Join(fieldsMissing, "\n"))
	}

	return nil
}

func (c *config) writeCfgFile(filePath string, oauth2Resp *oauth2.Token) error {

	cfg := &config{
		AppID:        c.AppID,
		AppKey:       c.AppKey,
		AccessToken:  oauth2Resp.AccessToken,
		RefreshToken: oauth2Resp.RefreshToken,
		TokenType:    oauth2Resp.TokenType,
		Expiry:       oauth2Resp.Expiry,
	}

	jsonData, err := json.MarshalIndent(&cfg, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "Unable to parse JSON data: %#v", cfg)
	}

	if err := os.WriteFile(filePath, jsonData, 0600); err != nil {
		return errors.Wrapf(err, "Unable to write JSON data to config file: %v", filePath)
	}

	return nil
}

// Get the path of the configuration file
// On Unix/Linux: $XDG_DATA_HOME/go-inoreader.json
// On Windows: %APPDATA%\go-inoreader.json
func getCfgFilePath() string {

	homeDir, _ := os.UserHomeDir()
	var fileName string

	switch runtime.GOOS {
	case "windows":
		fileName = path.Join(homeDir, "/AppData/Roaming/go-inoreader.json")

	default:
		fileName = path.Join(homeDir, "/.local/share/go-inoreader.json")
	}

	return fileName
}
