package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	flag "github.com/spf13/pflag"
)

type configuration struct {
	RefreshInterval int                 `koanf:"interval"`
	Secret          secretConfiguration `koanf:"secret"`
	Oidc            oidcConfiguration   `koanf:"oidc"`
	Log             logConfiguration    `koanf:"log"`

	DummyProvider bool `koanf:"dummyprovider"`
}

type secretConfiguration struct {
	Name      string `koanf:"name"`
	Namespace string `koanf:"namespace"`
	Key       string `koanf:"key"`
}

type logConfiguration struct {
	Level  string `koanf:"level"`
	Format string `koanf:"format"`
}

const (
	TextFormat = "text"
	JSONFormat = "json"

	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
)

type oidcConfiguration struct {
	TokenUrl     string `koanf:"tokenurl"`
	ClientID     string `koanf:"clientid"`
	ClientSecret string `koanf:"clientsecret"`
}

func getConfig() (configuration, error) {
	c := configuration{
		Secret: secretConfiguration{
			Key: "token",
		},
		RefreshInterval: 500,
		DummyProvider:   false,
		Log: logConfiguration{
			Level:  LevelInfo,
			Format: TextFormat,
		},
	}
	k := koanf.New(".")

	f := flag.NewFlagSet("config", flag.ContinueOnError)
	configFile := f.StringP("config", "f", "", "path to configuration file")
	f.String("secret.name", "", "name of the secret to update")
	f.String("secret.namespace", "", "namespace of the secret to update")
	f.String("secret.key", c.Secret.Key, "name of the secret to update")

	f.Int("interval", c.RefreshInterval, "interval in seconds to update the token")
	f.String("log.level", c.Log.Level, "the log level, on of [debug, info, warn]")
	f.String("log.format", c.Log.Format, "how to format the log, on of [text, json]")

	f.String("oidc.tokenurl", "", "the OIDC token endpoint")
	f.String("oidc.clientid", "", "the OIDC client id")
	f.String("oidc.clientsecret", "", "the OIDC client secret")
	f.Parse(os.Args[1:])

	if *configFile != "" {
		if err := k.Load(file.Provider(*configFile), yaml.Parser()); err != nil {
			return c, fmt.Errorf("unable to read file %s: %w", *configFile, err)
		}
	}

	if err := k.Load(env.Provider("KTR_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "KTR_")), "_", ".")
	}), nil); err != nil {
		return c, fmt.Errorf("unable to parse environemnt variables: %w", err)
	}

	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		return c, fmt.Errorf("unable to parse commandline flags: %w", err)
	}

	if err := k.Unmarshal("", &c); err != nil {
		return c, fmt.Errorf("could not load config: %w", err)
	}

	if err := validateConfig(c); err != nil {
		return c, err
	}
	return c, nil
}

func validateConfig(c configuration) error {
	if c.Secret.Name == "" {
		return errors.New("secret name may not be empty")
	}
	if c.Secret.Namespace == "" {
		return errors.New("secret namespace may not be empty")
	}
	if c.Secret.Key == "" {
		return errors.New("secret key may not be empty")
	}
	return nil
}
