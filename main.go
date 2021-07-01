package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	kconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// default to text formatter to get human readable output in case of
	// setup failure
	textf := new(log.TextFormatter)
	textf.TimestampFormat = "2006-01-02 15:04:05"
	textf.FullTimestamp = true
	log.SetFormatter(textf)

	conf, err := getConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to get config")
	}

	// Setup logger
	l, err := log.ParseLevel(conf.Log.Level)
	if err != nil {
		log.WithError(err).Error("Failed to set log level. Falling back to debug.")
		l = log.DebugLevel
	}
	log.SetLevel(l)

	switch conf.Log.Format {
	case JSONFormat:
		log.SetFormatter(&log.JSONFormatter{})
	default:
		// Keep the text logger
	}

	kc, err := kconfig.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to get kube config")
	}
	c, err := kclient.New(kc, kclient.Options{})
	if err != nil {
		log.WithError(err).Fatal("Failed to get kube client")
	}

	var provider tokenProvider
	if conf.Oidc.TokenUrl != "" {
		provider = &oidcProvider{
			client: &http.Client{
				Timeout: 10 * time.Second,
			},
			tokenUrl:     conf.Oidc.TokenUrl,
			clientId:     conf.Oidc.ClientID,
			clientSecret: conf.Oidc.ClientSecret,
		}
	} else if conf.DummyProvider {
		provider = &dummyProvider{}
	} else {
		log.Fatal("No priovider configured")
	}

	r := refresher{
		name:      conf.Secret.Name,
		namespace: conf.Secret.Namespace,
		key:       conf.Secret.Key,
		Client:    c,
		provider:  provider,
	}

	ticker := time.NewTicker(time.Duration(conf.RefreshInterval) * time.Second)
	defer ticker.Stop()

	// TODO(glrf) liveness / readiness

	stopped := false
	backoff := 500 * time.Millisecond
	maxbackoff := 30 * time.Second
	for !stopped {
		err = r.refresh(ctx)
		if err != nil {
			log.WithField("backoff", backoff).WithError(err).Error("Failed to refresh secret")
			select {
			case <-ctx.Done():
				// Will continue to next select which will handle the termination
			case <-time.After(backoff):
				backoff = 2 * backoff
				if backoff > maxbackoff {
					backoff = maxbackoff
				}
				continue
			}
		} else {
			log.Info("Refreshed token")
			backoff = 500 * time.Millisecond
		}
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			stop()
			stopped = true
			log.Warn("Terminating..")
		}
	}
}
