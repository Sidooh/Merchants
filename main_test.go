package main

import (
	"github.com/stretchr/testify/assert"
	"merchants.sidooh/api"
	"net/http/httptest"
	"testing"
)

func TestRouting(t *testing.T) {
	tests := []struct {
		description  string // description of the test case
		route        string // route path to test
		expectedCode int    // expected HTTP status code
	}{
		{
			description:  "get HTTP status 200",
			route:        "/200",
			expectedCode: 200,
		},
		{
			description:  "get 401, when unauthorized",
			route:        "/401",
			expectedCode: 401,
		},
	}

	app := api.Server()

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest("GET", test.route, nil)

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, 1)

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}

}
