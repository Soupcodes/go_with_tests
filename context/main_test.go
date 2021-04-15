package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SpyStore struct {
	response  string
	cancelled bool
	t         *testing.T
}

func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()
	if !s.cancelled {
		s.t.Errorf("store was not told to cancel")
	}
}

func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()
	if s.cancelled {
		s.t.Errorf("store was told to cancel")
	}
}

func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.response
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}

func TestServer(t *testing.T) {

	t.Run("returns data from store", func(t *testing.T) {
		data := "hello, world"
		store := &SpyStore{response: data, t: t}
		svr := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		svr.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
		}
		store.assertWasNotCancelled()
	})

	t.Run("Tell store to cancel work if request is cancelled", func(t *testing.T) {
		data := "Hello World!"
		store := &SpyStore{response: data, t: t}
		srv := Server(store)

		req := httptest.NewRequest(http.MethodGet, "/", nil)   // Send this as a request query into ServeHTTP
		cancelCtx, cancel := context.WithCancel(req.Context()) // Returns a copy of req.Context() (the parent Context) with a new Done channel. This done channel is closed when 'cancel' is invoked or the parents' done channel is closed
		// cancel: when invoked, tells an operation to abandon its work without waiting for it to stop
		// cancelCtx will capture the current state of the req.Context()
		time.AfterFunc(5*time.Millisecond, cancel) // invoke cancel after a delay
		req = req.WithContext(cancelCtx)           // reassign req with the updated context captured by cancelCtx
		res := httptest.NewRecorder()              // This will capture the response of the request
		srv.ServeHTTP(res, req)

		store.assertWasCancelled()
	})
}
