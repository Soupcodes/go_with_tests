package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*
When testing calls to third party apis, ideally, we don't want to actually rely on their responses
as these can be slow and doesn't enable us to test edge cases. Use httptest instead
*/

func TestRacer(t *testing.T) {

	t.Run("Return the fastest responding url", func(t *testing.T) {
		// Create a dummy server with a mux handler to simulate a backend api response against an open local port
		slowServer := makeDummyServer(20 * time.Millisecond)
		fastServer := makeDummyServer(0 * time.Millisecond)
		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL // localhost url:port
		fastURL := fastServer.URL

		want := fastURL
		got, err := Racer(slowURL, fastURL)

		if err != nil {
			t.Fatalf("did not expect an error but got one %v", err)
		}

		if got != want {
			t.Fatalf("Got: %s, Want: %s", got, want)
		}
	})

	t.Run("Timeout if requests take longer than 10 secs", func(t *testing.T) {
		srv := makeDummyServer(20 * time.Millisecond)
		defer srv.Close()

		timeout := 10 * time.Millisecond
		_, err := ConfigurableRacer(srv.URL, srv.URL, timeout)

		if err == nil {
			t.Error("Expected an error, got none")
		}
	})

}

func makeDummyServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}
