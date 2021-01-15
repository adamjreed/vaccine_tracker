package app

import (
	"covid/internal/config"
	"covid/internal/store"
	"github.com/sfreiberg/gotwilio"
)

type Container struct {
	Config         *config.Config
	Twilio         *gotwilio.Twilio
	ParameterStore store.Store
}

func (c *Container) Cleanup() {
	//clean up and long running things here
}

func CreateContainer() (*Container, error) {
	paramStore, err := store.NewAWS()
	if err != nil {
		return nil, err
	}

	conf, err := config.LoadConfig(paramStore)
	if err != nil {
		return nil, err
	}

	return &Container{
		Config:         conf,
		Twilio:         gotwilio.NewTwilioClient(conf.Twilio.AccountSID, conf.Twilio.AuthToken),
		ParameterStore: paramStore,
	}, nil
}
