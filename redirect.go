package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gentlemanautomaton/signaler"
)

func main() {
	// Capture shutdown signals
	shutdown := signaler.New().Capture(os.Interrupt, syscall.SIGTERM)

	// Parse arguments and environment
	var prefix string
	flag.StringVar(&prefix, "prefix", "", "URL Prefix")
	flag.Parse()
	if prefix == "" {
		prefix = os.Getenv("REDIR_PREFIX")
	}
	if prefix == "" {
		prefix = "/"
	}

	// Register the handler with the default muxer
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		location := fmt.Sprintf("%s%s?%s", prefix, r.URL.Path, r.URL.RawQuery)
		fmt.Println(location)
		http.Redirect(w, r, location, 301)
	})

	// Create a server (which will use the default muxer by default)
	s := &http.Server{Addr: ":80"}

	// Tell the server to stop gracefully when a shutdown signal is received
	stopped := shutdown.Then(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		s.Shutdown(ctx)
	})

	// Always cleanup and wait until the shutdown has completed
	defer stopped.Wait()
	defer shutdown.Trigger()

	// Run the server and print the final result
	log.Println(s.ListenAndServe())
}
