package main

import (
	"errors"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type option struct {
	SecretName      string
	SecretNamespace string
	SecretKey       string

	RefreshInterval int
	Oidc            oidcOption
}

type oidcOption struct {
	TokenUrl     string
	ClientID     string
	ClientSecret string
}

func getConfig() (option, error) {
	o := option{
		RefreshInterval: 595,
	}

	configFile := flag.StringP("config", "f", "", "path to configuration file")
	flag.Parse()

	// Get config file
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/kube-token-refresher")
	viper.AddConfigPath("$HOME/.config/kube-token-refresher")
	viper.AddConfigPath("$HOME/.kube-token-refresher")

	viper.SetEnvPrefix("ktr")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.BindEnv("secretName")
	viper.BindEnv("secretNamespace")
	viper.BindEnv("secretKey")
	viper.BindEnv("oidc.tokenUrl")
	viper.BindEnv("oidc.clientID")
	viper.BindEnv("oidc.clientSecret")

	if *configFile != "" {
		viper.SetConfigFile(*configFile)
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil && !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		return o, err
	}
	if err != nil && errors.As(err, &viper.ConfigFileNotFoundError{}) &&
		*configFile != "" {
		return o, err
	}

	err = viper.Unmarshal(&o)
	if err != nil {
		return o, err
	}

	if err = validateConfig(o); err != nil {
		return o, err
	}
	return o, nil
}

func validateConfig(o option) error {
	if o.SecretName == "" {
		return errors.New("secret name may not be empty")
	}
	if o.SecretNamespace == "" {
		return errors.New("secret namespace may not be empty")
	}
	if o.SecretKey == "" {
		return errors.New("secret key may not be empty")
	}
	return nil
}
