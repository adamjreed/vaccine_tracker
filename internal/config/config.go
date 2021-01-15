package config

import (
	"covid/internal/store"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Stage  string       `json:"stage"`
	Check  CheckConfig  `json:"check"`
	SMS    SMSConfig    `json:"sms"`
	Twilio TwilioConfig `json:"twilio"`
}

type CheckConfig struct {
	URL    string `json:"url"`
	Phrase string `json:"phrase"`
}

type SMSConfig struct {
	To   []string `json:"to"`
	From string `json:"from"`
}

type TwilioConfig struct {
	AccountSID      string `json:"account_sid"`
	AuthTokenSecret string `json:"auth_token_secret"`
	AuthToken       string `json:"auth_token"`
}

func LoadConfig(aws *store.AWS) (*Config, error) {
	var config Config
	stage := os.Getenv("STAGE")

	configJson, err := aws.GetParam(fmt.Sprintf("/%s/vaccine_tracker/config", stage), false)
	if err != nil {
		return nil, err
	}

	if configJson == "" {
		return nil, errors.New("could not load config")
	}

	err = json.Unmarshal([]byte(configJson), &config)
	if err != nil {
		return nil, err
	}

	if config.Twilio.AuthToken == "" {
		config.Twilio.AuthToken, err = aws.GetParam(config.Twilio.AuthTokenSecret, true)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}
