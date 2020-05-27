package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Api      string `json:"api"`
	AppsDomain string `json:"apps_domain"`
	User     string `json:"user"`
	Password string `json:"pass"`
	Org      string `json:"org"`
	Space    string `json:"space"`
	Timeout  TimeoutConfig `json:"timeout"`
}

type TimeoutConfig struct {
	Staging int `json:"staging"` // minutes
	Startup int `json:"startup"` // minutes
	FirstHealthyResponse int `json:"firstHealthyResponse"` // seconds
}

func (c *Config) LoadConfig(filename string) error {
	json_byte, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	err = json.Unmarshal(json_byte, &c)
	if err != nil {
		return err
	}

	// defaults
	if c.AppsDomain == "" {
	    c.AppsDomain = "cfapps.com"
	}

	return nil
}
