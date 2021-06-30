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

	tcs := []struct {
		input string
		code  int
		fail  bool
	}{
		{
			input: `{"access_token": "aaa"}`,
			code:  200,
			fail:  false,
		},
		{
			input: `{"error": "invalid_client", "error_description": "unauthorized"}`,
			code:  401,
			fail:  true,
		},
		{
			input: `{}`,
			code:  200,
			fail:  true,
		},
		{
			input: `invalid`,
			code:  200,
			fail:  true,
		},
	}

	for _, tc := range tcs {
		response = tc.input
		code = tc.code
		_, err := p.GetToken(ctx)
		if !tc.fail && err != nil {
			t.Fatalf("[input: %s] failed to get token: %s", tc.input, err.Error())
		}
		if tc.fail && err == nil {
			t.Fatalf("[input: %s] should have failed but didn't", tc.input)
		}
	}

}
