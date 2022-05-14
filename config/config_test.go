package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	var cfg = &config{}

	cfg, err := loadConfig(getCfgFilePath())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", &cfg)
}
