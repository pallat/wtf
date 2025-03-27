package httpclient

import (
	"net/http"
	"testing"
)

func TestOption(t *testing.T) {
	t.Run("should be able to set Authorization header", func(t *testing.T) {
		token := "eyJheader.payload.signature"
		want := "Bearer " + token

		req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
		opt := AuthorizationOption(token)

		opt(req)

		got := req.Header.Get("Authorization")
		if got != want {
			t.Errorf("Authorization header got: %s, want: %s", got, want)
		}
	})
}
