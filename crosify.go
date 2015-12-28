package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	address = flag.String("address", "0.0.0.0", "Listening address")
	port    = flag.String("port", "8080", "Listening port")
	target  = flag.String("url", "", "Target URL")
)

type proxy struct {
	rp   *httputil.ReverseProxy
	host string
}

func newProxy(targetURL string) *proxy {
	target, _ := url.Parse(targetURL)
	return &proxy{httputil.NewSingleHostReverseProxy(target), target.Host}
}

func (p *proxy) handle(rw http.ResponseWriter, req *http.Request) {
	req.Host = p.host
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	p.rp.ServeHTTP(rw, req)
}

func main() {
	flag.Parse()
	listen := *address + ":" + *port
	p := newProxy(*target)
	http.HandleFunc("/", p.handle)
	log.Printf("Listening on %s", listen)
	log.Printf("Proxying for %s", *target)
	log.Fatal(http.ListenAndServe(listen, nil))
}
