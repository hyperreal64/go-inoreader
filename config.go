package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type config struct {
	*oauth2.Token
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

func writeCfgFile(filepath string, oauth2Resp *oauth2.Token) error {

	jsonData, err := json.MarshalIndent(&oauth2Resp, "", "    ")
	if err != nil {
		return errors.Wrap(err, "Could not marshal JSON data to file")
	}

	// Note: File permission mode argument '0666' assumes environment umask is set to 0022,
	// which is the default on most Linux distributions. After umask is applied, the result
	// would be 0644 (rw-r--r--).
	// As per https://golang.org/pkg/os/#FileMode, the mode bits have the same definition
	// on all systems, thus '0666' can be used for Windows and Unix alike.
	if err = ioutil.WriteFile(filepath, jsonData, 0666); err != nil {
		return errors.Wrapf(err, "Could not write to config file: %s\n", filepath)
	}

	return nil
}
