package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/jaschaephraim/lrserver"
)

// Livereload
// Include the following code to have livereload functionnality in you webpage:
// <script>
//   document.write('<script src="http://' + (location.host || 'localhost').split(':')[0] +
//   ':35729/livereload.js?snipver=1"></' + 'script>')
// </script>

func startLiveReload() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	// get current path
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	log.Println("Listening on path: ", ex)
	err = watcher.Add(ex)
	if err != nil {
		log.Fatalln(err)
	}

	lr := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	go lr.ListenAndServe()

	for {
		select {
		case event := <-watcher.Events:
			lr.Reload(event.Name)
		case err := <-watcher.Errors:
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func main() {
	var host = flag.String("host", ":8000", "Host adresse to use for static web server")
	var verbose = flag.Bool("v", false, "set verbosity")
	var liveReload = flag.Bool("lr", false, "Enable livereload functionnality")
	flag.Parse()

	if *liveReload {
		go startLiveReload()
	}

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

	if *verbose {
		log.Print("Liste of interfaces: ", getInterfaces())
	}
	log.Print("Listening on ", *host)
	log.Fatal(http.ListenAndServe(*host, nil))
}

func getInterfaces() (intList []string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Print("Failed to retrieve list of interfaces")
	} else {
		for _, i := range interfaces {
			addrs, err := i.Addrs()
			if err != nil {
				log.Println("Failed to retrieve list of addresses for interface ", i.Name)
			}
			for _, j := range addrs {
				ip, _, err := net.ParseCIDR(j.String())
				if err != nil {
					log.Println("Failed to extract IP from interface")
				}
				intList = append(intList, ip.String())
			}
		}
	}
	return intList
}
