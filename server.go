package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"jmdict/api"
)

func main() {
	port := flag.String("port", "8080", "Port to run server on")
	flag.Parse()

	// Set up routes
	http.HandleFunc("/api/word/", func(w http.ResponseWriter, r *http.Request) {
		api.Handler(w, r)
	})

	fmt.Printf("Server running on http://localhost:%s\n", *port)
	fmt.Printf("Try: http://localhost:%s/api/word/言葉\n", *port)

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
