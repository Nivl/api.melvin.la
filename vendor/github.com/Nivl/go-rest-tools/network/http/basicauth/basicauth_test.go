package basicauth_test

import (
	"testing"

	"github.com/Nivl/go-rest-tools/network/http/basicauth"
	"github.com/stretchr/testify/assert"
)

func TestValidBasicAuthHeader(t *testing.T) {
	cases := []struct {
		description string

		// Input to look into
		headers []string

		// Data to look for
		realm string

		// Result expected
		expectedName string
		expectedPsw  string
	}{
		{
			description:  "valid basic with no real",
			headers:      []string{"basic dXNlcjpwYXNzd29yZA=="},
			realm:        "",
			expectedName: "user",
			expectedPsw:  "password",
		},
		{
			description:  "valid basic with realm",
			headers:      []string{`basic dXNlcjpwYXNzd29yZA== realm="myRealm"`},
			realm:        "myRealm",
			expectedName: "user",
			expectedPsw:  "password",
		},
		{
			description:  "valid basic with realm at a different place",
			headers:      []string{`basic realm="myRealm" dXNlcjpwYXNzd29yZA==`},
			realm:        "myRealm",
			expectedName: "user",
			expectedPsw:  "password",
		},
		{
			description: "valid basic in a list",
			headers: []string{
				`basic dXNlcjpwYXNzd29yZA== realm="myRealm"`, // wrong realm
				"basic dXNlcjpwYXNzd29yZA==",                 // valid
			},
			realm:        "",
			expectedName: "user",
			expectedPsw:  "password",
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			name, psw, _ := basicauth.ParseAuthHeader(tc.headers, "basic", tc.realm)

			assert.Equal(t, tc.expectedName, name)
			assert.Equal(t, tc.expectedPsw, psw)
		})
	}
}

func TestInvalidBasicAuthHeader(t *testing.T) {
	cases := []struct {
		description string

		// Input to look into
		headers []string

		// Data to look for
		realm string
	}{
		{
			description: "Wrong credential",
			headers:     []string{"basic user:password"},
			realm:       "",
		},
		{
			description: "Wrong Bearer",
			headers:     []string{`baic dXNlcjpwYXNzd29yZA==`},
			realm:       "",
		},
		{
			description: "Missing Bearer",
			headers:     []string{`dXNlcjpwYXNzd29yZA==`},
			realm:       "",
		},
		{
			description: "Missing Bearer with realm",
			headers:     []string{`dXNlcjpwYXNzd29yZA== realm="myRealm"`},
			realm:       "myRealm",
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			name, psw, _ := basicauth.ParseAuthHeader(tc.headers, "basic", tc.realm)

			assert.Empty(t, name)
			assert.Empty(t, psw)
		})
	}
}

func TestValidPasswordAuthHeader(t *testing.T) {
	cases := []struct {
		description string

		// Input to look into
		headers []string

		// Data to look for
		realm string

		// Result expected
		expectedPsw string
	}{
		{
			description: "valid with no real",
			headers:     []string{"password cGFzc3dvcmQ="},
			realm:       "",
			expectedPsw: "password",
		},
		{
			description: "valid with realm",
			headers:     []string{`password cGFzc3dvcmQ= realm="myRealm"`},
			realm:       "myRealm",
			expectedPsw: "password",
		},
		{
			description: "valid basic with realm at a different place",
			headers:     []string{`password realm="myRealm" cGFzc3dvcmQ=`},
			realm:       "myRealm",
			expectedPsw: "password",
		},
		{
			description: "valid basic in a list",
			headers: []string{
				`password cGFzc3dvcmQ= realm="myRealm"`, // wrong
				"password cGFzc3dvcmQ=",                 // valid
			},
			realm:       "",
			expectedPsw: "password",
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			name, psw, _ := basicauth.ParseAuthHeader(tc.headers, "password", tc.realm)

			assert.Empty(t, name)
			assert.Equal(t, tc.expectedPsw, psw)
		})
	}
}

func TestInvalidPasswordAuthHeader(t *testing.T) {
	cases := []struct {
		description string

		// Input to look into
		headers []string

		// Data to look for
		realm string
	}{
		{
			description: "Wrong credential",
			headers:     []string{"password user:password"},
			realm:       "",
		},
		{
			description: "Wrong Bearer",
			headers:     []string{`psw cGFzc3dvcmQ=`},
			realm:       "",
		},
		{
			description: "Missing Bearer",
			headers:     []string{`cGFzc3dvcmQ=`},
			realm:       "",
		},
		{
			description: "Missing Bearer with realm",
			headers:     []string{`cGFzc3dvcmQ= realm="myRealm"`},
			realm:       "myRealm",
		},
	}

	for _, tc := range cases {
		t.Run(tc.description, func(t *testing.T) {
			name, psw, _ := basicauth.ParseAuthHeader(tc.headers, "password", tc.realm)

			assert.Empty(t, name)
			assert.Empty(t, psw)
		})
	}
}
