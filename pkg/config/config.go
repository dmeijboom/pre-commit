package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Checks []Check `json:"checks"`
}

func Load(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	if err := json.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, nil
}
