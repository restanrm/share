package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var host = flag.String("host", ":8000", "Host adresse to use for static web server")
	var verbose = flag.Bool("v", false, "set verbosity")
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("client=", r.RemoteAddr, " Method=", r.Method, " URL=", r.URL.Path)
		if *verbose {
			r.ParseForm()
			log.Print("List of headers: ")
			for header, value := range r.Header {
				log.Print("\t", header, ": ", value)
			}
			log.Print("List parameters found in the request: ")
			for form, value := range r.Form {
				log.Print("\t", form, ": ", value)
			}
			log.Print()
		}
		http.FileServer(http.Dir(".")).ServeHTTP(w, r)
	})
	log.Fatal(http.ListenAndServe(*host, nil))
}
