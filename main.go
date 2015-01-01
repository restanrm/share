package main 

import (
	"net/http"
	"flag"
)

func main() {
	var host = flag.String("host", ":8000", "Host adresse to use for static web server")
	flag.Parse()
	http.ListenAndServe(*host, http.FileServer(http.Dir(".")))
}
