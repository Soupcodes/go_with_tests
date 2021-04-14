package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

/*
Using mocking, we can inject 'dependencies' to spy on calls in the functions we are testing.
Eg. in this program, we wouldn't want to actually wait a second between every number we're counting down
from during testing
*/

type Sleeper interface {
	Sleep()
}

type DefaultSleeper struct{}

func (d *DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

type ConfigurableSleeper struct {
	duration time.Duration
	sleep    func(time.Duration)
}

func (c *ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
}

const (
	start     = 5
	finalWord = "Go!"
)

func Countdown(w io.Writer, sleeper Sleeper) { // Defines that the sleeper struct needs to have the behaviour to be able to Sleep() under the hood...
	for i := start; i > 0; i-- {
		fmt.Fprintln(w, i)
		sleeper.Sleep()
	}
	fmt.Fprint(w, finalWord)
}

func main() {
	duration := 1 * time.Second
	sleeper := &ConfigurableSleeper{duration: duration, sleep: time.Sleep}
	Countdown(os.Stdout, sleeper)
}
