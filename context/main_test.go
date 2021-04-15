package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SpyStore struct {
	response string
	t        *testing.T
}

// func (s *SpyStore) assertWasCancelled() {
// 	s.t.Helper()
// 	if !s.cancelled {
// 		s.t.Errorf("store was not told to cancel")
// 	}
// }

// func (s *SpyStore) assertWasNotCancelled() {
// 	s.t.Helper()
// 	if s.cancelled {
// 		s.t.Errorf("store was told to cancel")
// 	}
// }

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	data := make(chan string, 1)

	go func() {
		var result string
		for _, c := range s.response {
			select {
			case <-ctx.Done(): // ctx holds a Done channel
				s.t.Log("Spy store got cancelled")
				return
			default:
				time.Sleep(100 * time.Millisecond) // Artificially create a delay in getting the response so cancelling can be tested
				result += string(c)
			}
		}
		data <- result
	}()

	// Blocking select statement will either return the data to the user or cancel the store, depending on what request comes back first
	// Select RACES asynchronous processes against each other and exits when one completes
	select {
	case <-ctx.Done():
		return "", ctx.Err() // context also holds an error type that can be returned as an error
	case res := <-data:
		return res, nil
	}
}

type SpyResponseWriter struct {
	written bool
}

func (s *SpyResponseWriter) Header() http.Header {
	s.written = true
	return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
	s.written = true
	return 0, errors.New("not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
	s.written = true
}

func TestServer(t *testing.T) {

	t.Run("returns data from store", func(t *testing.T) {
		data := "hello, world"
		store := &SpyStore{response: data, t: t}
		srv := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		srv.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
		}
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
		res := &SpyResponseWriter{}                // This will capture the response of the request
		srv.ServeHTTP(res, req)

		if res.written {
			t.Error("a response should not have been written")
		}
	})
}
