package main

import (
	"fmt"
	"net/http"
)

type Store interface {
	Fetch() string
	Cancel()
}

func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context() // Capture the current context, which holds a done channel
		data := make(chan string, 1)

		go func() {
			data <- store.Fetch() // Run the query you need concurrently and add result to the data channel
		}()

		// Blocking select statement will either return the data to the user or cancel the store, depending on what request comes back first
		// Select RACES asynchronous processes against each other and exits when one completes
		select {
		case d := <-data:
			fmt.Fprint(w, d)
		case <-ctx.Done():
			store.Cancel()
		}
	}
}
