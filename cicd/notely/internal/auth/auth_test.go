package auth

import (
	"net/http"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	testsWithAuth := map[string]struct {
		input string
		want  string
	}{
		"simple": {input: "abc123", want: "abc123"},
		"no_key": {input: "", want: ""},
	}

	for name, tc := range testsWithAuth {
		t.Run(name, func(t *testing.T) {
			header := http.Header{}
			header.Add("Authorization", "ApiKey "+tc.input)
			got, err := GetAPIKey(header)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}

	t.Run("no_auth", func(t *testing.T) {
		header := http.Header{}
		got, err := GetAPIKey(header)
		if err != ErrNoAuthHeaderIncluded {
			t.Fatal("Wrong error raised from missing auth header")
		}
		if got != "" {
			t.Fatal("Wrong value recieved from missing auth header")
		}
	})

	t.Run("malformed_auth", func(t *testing.T) {
		header := http.Header{}
		header.Add("Authorization", "ApiKey123")
		got, err := GetAPIKey(header)
		if err != ErrMalformedAuthHeader {
			t.Fatal("Wrong error raised from malformed auth header")
		}
		if got != "" {
			t.Fatal("Wrong value recieved from malformed auth header")
		}
	})
}
