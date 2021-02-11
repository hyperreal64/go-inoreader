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

// FIXME: refactor and make this more efficient
//  This most likely does not need to be a strucc
//  Store / pick up client ID and client secret from config
//  Unexport these, since they are all in one package
type cfgFile struct {
	FilePath string
	Contents *Configuration
}

// Configuration ---
type Configuration struct {
	Oauth2Response *oauth2.Token
}

// GetcfgFilePath ---
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

func newCfgFile(filePath string, data []byte) (*cfgFile, error) {

	conf := &Configuration{}
	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal JSON data")
	}

	return &cfgFile{
		FilePath: filePath,
		Contents: conf,
	}, nil
}

func (cf *cfgFile) writeCfgFile() error {

	jsonData, err := json.MarshalIndent(&cf.Contents, "", "    ")
	if err != nil {
		return errors.Wrap(err, "Could not marshal JSON data to file")
	}

	// Note: File permission mode argument '0666' assumes environment umask is set to 0022,
	// which is the default on most Linux distributions. After umask is applied, the result
	// would be 0644 (rw-r--r--).
	// As per https://golang.org/pkg/os/#FileMode, the mode bits have the same definition
	// on all systems, thus '0666' can be used for Windows and Unix alike.
	if err = ioutil.WriteFile(cf.FilePath, jsonData, 0666); err != nil {
		return errors.Wrapf(err, "Could not write to config file: %s\n", cf.FilePath)
	}

	return nil
}
