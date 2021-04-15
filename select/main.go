package main

import (
	"fmt"
	"net/http"
	"time"
)

const timeout = 10 * time.Second

func Racer(a, b string) (string, error) {
	return ConfigurableRacer(a, b, timeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (string, error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout): // If a function uses time, and needs testing, try to make it customisable so you're not waiting an excessively long time for the test to finish
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		ch <- struct{}{}
	}()
	return ch
}
