package main

// TODO(glrf) Maybe switch to a full blown OIDC client, which would enable us to
//  get the actual expiration time of the token.

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type oidcProvider struct {
	client *http.Client

	tokenUrl     string
	clientId     string
	clientSecret string
}

func (p *oidcProvider) GetToken(context.Context) ([]byte, error) {
	raw, err := p.client.PostForm(p.tokenUrl, url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {p.clientId},
		"client_secret": {p.clientSecret},
	})
	if err != nil {
		return nil, err
	}
	if raw.Body == nil {
		return nil, fmt.Errorf("missing body in response")
	}
	defer raw.Body.Close()
	body, err := io.ReadAll(raw.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	if raw.StatusCode != http.StatusOK {
		errResp := struct {
			Error       string `json:"error"`
			Description string `json:"error_description"`
		}{}
		err := json.Unmarshal(body, &errResp)
		if err != nil {
			return nil, fmt.Errorf("error %s: uable to parse error response %w", raw.Status, err)
		}
		return nil, fmt.Errorf("error %s: %s", raw.Status, errResp.Description)
	}

	resp := struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, fmt.Errorf("uable to parse response %w", err)
	}

	if resp.AccessToken == "" {
		return nil, fmt.Errorf("invalid response, no token provided")
	}

	return []byte(resp.AccessToken), nil
}
