package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	kclient "sigs.k8s.io/controller-runtime/pkg/client"
	kconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	opt, err := getConfig()
	if err != nil {
		log.Fatalf("Failed to get config: %s\n", err)
	}

	conf, err := kconfig.GetConfig()
	if err != nil {
		log.Fatalf("Failed to get kube config: %s\n", err)
	}
	c, err := kclient.New(conf, kclient.Options{})
	if err != nil {
		log.Fatalf("Failed to get kube client: %s\n", err)
	}

	var provider tokenProvider
	if opt.Oidc.TokenUrl != "" {
		provider = &oidcProvider{
			client: &http.Client{
				Timeout: 10 * time.Second,
			},
			tokenUrl:     opt.Oidc.TokenUrl,
			clientId:     opt.Oidc.ClientID,
			clientSecret: opt.Oidc.ClientSecret,
		}
	} else {
		log.Fatalln("No priovider configured")
	}

	r := refresher{
		name:      opt.SecretName,
		namespace: opt.SecretNamespace,
		key:       opt.SecretKey,
		Client:    c,
		provider:  provider,
	}

	ticker := time.NewTicker(time.Duration(opt.RefreshInterval) * time.Second)
	defer ticker.Stop()

	// TODO(glrf) liveness / readiness

	stopped := false
	for !stopped {
		err = r.refresh(ctx)
		if err != nil {
			// TODO(glrf) Retries
			log.Printf("Failed to refresh secret: %s\n", err)
		} else {
			log.Println("Refreshed token")
		}
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			stop()
			stopped = true
			log.Printf("Terminating.. \n")
		}
	}
}
