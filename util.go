package main

import "net/http"

func maybeNotify(url string) {
	if url != "" {
		http.Get(url)
	}
}
