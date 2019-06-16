package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kyokomi/emoji"
)

func root(emoji string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		t := time.NewTicker(100 * time.Millisecond)
		defer t.Stop()
		flusher, _ := w.(http.Flusher)
		for {
			select {
			case <-req.Context().Done():
				return
			case <-t.C:
				_, err := w.Write([]byte(emoji))
				if err != nil {
					return
				}
				flusher.Flush()
			}
		}
	}
}

func main() {
	// Accept a command line flag "-httpaddr :8083"
	// This flag tells the server the http address to listen on
	httpaddr := flag.String("httpaddr", "localhost:8083",
		"the address/port to listen on for http \n"+
			"use :<port> to listen on all addresses\n")

	emojicode := flag.String("emoji", ":ghost:",
		"the emoji shortcode to stream \n"+
			"http://www.unicode.org/emoji/charts/emoji-list.html\n")

	flag.Parse()

	em := strings.TrimSuffix(emoji.Sprint(*emojicode), " ")
	http.HandleFunc("/", root(em))

	log.Println("Streaming emoji: ", em)
	log.Println("Listening on http:", *httpaddr)
	log.Fatal(http.ListenAndServe(*httpaddr, nil))
}
