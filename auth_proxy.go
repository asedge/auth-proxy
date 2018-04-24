package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

type ProxyHandler struct {
	Username, Password string
}

func (p *ProxyHandler) copyHeaders(source, destination http.Header) {
	for k, values := range source {
		for _, v := range values {
			destination.Add(k, v)
		}
	}
}

// Proxies an incoming request on behalf of the requester.
func (p *ProxyHandler) MakeProxiedRequest(original *http.Request, url string) (resp *http.Response, e error) {
	requestorAddr, _, _ := net.SplitHostPort(original.RemoteAddr)
	log.Printf("Requesting URL (%s) for client (%s).", url, requestorAddr)
	req, err := http.NewRequest(original.Method, url, nil)
	if err != nil {
		log.Printf("Got error when making new request (method: %s, url: %s): %v", original.Method, url, err)
		return nil, err
	}
	req.SetBasicAuth(p.Username, p.Password)
	req.Header.Add("X-Forward-For", requestorAddr)
	p.copyHeaders(original.Header, req.Header)
	cli := &http.Client{}
	resp, e = cli.Do(req)
	return
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp, e := p.MakeProxiedRequest(req, req.RequestURI)
	if e != nil {
		log.Printf("Got error when requesting url (%s): %v", req.RequestURI, e)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("Unable to read response body: %v", readErr)
	}
	defer resp.Body.Close()
	p.copyHeaders(resp.Header, w.Header())
	// Set the status code from the proxied request.
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func main() {
	var user, pass string
	var ok bool

	if user, ok = os.LookupEnv("USERNAME"); !ok {
		log.Fatal("ERROR: Must set USERNAME env.")
	}

	if pass, ok = os.LookupEnv("PASSWORD"); !ok {
		log.Fatal("ERROR: Must set PASSWORD env.")
	}

	proxy := ProxyHandler{Username: user, Password: pass}
	server := &http.Server{
		Addr:    ":8989",
		Handler: &proxy,
	}

	log.Fatal(server.ListenAndServe())
}
