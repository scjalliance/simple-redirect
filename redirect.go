package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	var prefix string
	flag.StringVar(&prefix, "prefix", "", "URL Prefix")
	flag.Parse()
	if prefix == "" {
		prefix = os.Getenv("REDIR_PREFIX")
	}
	if prefix == "" {
		prefix = "/"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s%s?%s\n", prefix, r.URL.Path, r.URL.RawQuery)
		http.Redirect(w, r, fmt.Sprintf("%s%s?%s", prefix, r.URL.Path, r.URL.RawQuery), 301)
	})
	log.Fatal(http.ListenAndServe(":80", nil))
}
