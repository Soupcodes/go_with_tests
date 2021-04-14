package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

const write = "write"
const sleep = "sleep"

type SpySleeper struct {
	Calls int
}

// This will record how many times the Sleep method would have actually have been called in the real invokation of the tested function
func (s *SpySleeper) Sleep() {
	s.Calls++
}

type CountdownOperationsSpy struct {
	Calls []string
}

func (s *CountdownOperationsSpy) Sleep() {
	s.Calls = append(s.Calls, sleep)
}

func (s *CountdownOperationsSpy) Write([]byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
	s.durationSlept = duration
}

func TestCountdown(t *testing.T) {
	t.Run("prints 3 to Go!", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		Countdown(buffer, &CountdownOperationsSpy{})

		got := buffer.String()
		want := `3
2
1
Go!`

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("Ensure sleeps happen in between writes", func(t *testing.T) {
		spySleepPrinter := &CountdownOperationsSpy{}
		Countdown(spySleepPrinter, spySleepPrinter) // The struct has the writer and reader methods to be able to satisfy the underlying type Countdown asks for

		want := []string{write, sleep, write, sleep, write, sleep, write}

		if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
			t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
		}
	})
}

func TestConfigurableSleeper(t *testing.T) {
	sleepTime := 5 * time.Second

	// There must be a better way to do mocking in Go? This particular composition requires a lot of jumping around to follow the logic...
	spyTime := &SpyTime{}
	sleeper := ConfigurableSleeper{sleepTime, spyTime.Sleep}
	sleeper.Sleep()

	// spyTime.durationSlept is equal because it itself is passed into ConfigurableSleeper and gets the duration from that struct when sleeper.Sleep() is called
	if spyTime.durationSlept != sleepTime {
		t.Errorf("should have slept for %v but slept for %v", sleepTime, spyTime.durationSlept)
	}
}
