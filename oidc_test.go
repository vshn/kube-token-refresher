package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetToken(t *testing.T) {

	var response string
	code := 200
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		fmt.Fprintln(w, response)
	}))
	defer ts.Close()
	p := oidcProvider{
		client:       &http.Client{},
		tokenUrl:     ts.URL,
		clientId:     "testid",
		clientSecret: "testsecret",
	}

	ctx := context.Background()

	tcs := map[string]struct {
		input string
		code  int
		fail  bool
	}{
		"Valid200Resonse_ExpectNotFail": {
			input: `{"access_token": "aaa"}`,
			code:  200,
			fail:  false,
		},
		"Valid400Resonse_ExpectFail": {
			input: `{"error": "invalid_client", "error_description": "unauthorized"}`,
			code:  401,
			fail:  true,
		},
		"Empty200Resonse_ExpectFail": {
			input: `{}`,
			code:  200,
			fail:  true,
		},
		"Invalid200Resonse_ExpectFail": {
			input: `invalid`,
			code:  200,
			fail:  true,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			response = tc.input
			code = tc.code
			_, err := p.GetToken(ctx)
			if !tc.fail && err != nil {
				t.Fatalf("failed to get token: %s", err.Error())
			}
			if tc.fail && err == nil {
				t.Fatalf("should have failed but didn't")
			}
		})
	}

}
