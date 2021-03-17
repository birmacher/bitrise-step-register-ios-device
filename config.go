package main

import (
	"github.com/bitrise-io/go-steputils/stepconf"
)

type Config struct {
	APIKeyPath     stepconf.Secret `env:"api_key_path"`
	APIIssuer      string          `env:"api_issuer"`
	BuildAPIToken  string          `env:"build_api_token"`
	BuildURL       string          `env:"build_url"`
	DeviceName     string          `env:"device_name"`
	DeviceUDID     string          `env:"device_udid"`
	DevicePlatform string          `env:"device_platform"`
}
