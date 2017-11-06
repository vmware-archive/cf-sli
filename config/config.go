package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Api      string `json:"api"`
	User     string `json:"user"`
	Password string `json:"pass"`
	Domain   string `json:"domain"`
	Org      string `json:"org"`
	Space    string `json:"space"`
	Timeout  int `json:"timeout"`
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
	return nil
}
