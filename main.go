package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var host = flag.String("host", ":8000", "Host adresse to use for static web server")
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("client=", r.RemoteAddr, " URL=", r.URL.Path)
		http.FileServer(http.Dir(".")).ServeHTTP(w, r)
	})
	http.ListenAndServe(*host, nil)
}
