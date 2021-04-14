package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

/*
fmt.Printf prints to stdout which is tough for a test framework to capture,
so what needs to happen is to be able to INJECT the dependency of printing.
A good way to do this is by passing around types that have the behaviour of the io.Writer interface,
a general purpose type that allows you to control where to redirect your stdouts

Types of io.Writers:
- bytes.Buffer{}
- http.ResponseWriter
- os.Stdout
*/

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
	Greet(w, "world")
}

func main() {
	log.Fatal(http.ListenAndServe(":5000", http.HandlerFunc(MyGreeterHandler)))
}
