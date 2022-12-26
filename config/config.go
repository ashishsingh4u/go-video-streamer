package config

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	SERVER_IP string
	PORT      string
}

func GetConfig(config *Configuration, params ...string) (err error) {

	env := "dev"

	if len(params) > 0 {
		env = params[0]
	}

	fileName := fmt.Sprintf("./config/%s_config.yaml", env)

	err = gonfig.GetConf(fileName, config)

	return
}
