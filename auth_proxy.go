package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ProxyHandler struct {
	Username, Password, ProxyBase string
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
	log.Printf("Requesting URL (%s) for client (%s) with headers (%v).", url, original.RemoteAddr, original.Header)
	req, err := http.NewRequest(original.Method, url, nil)
	if err != nil {
		log.Printf("Got error when making new request (method: %s, url: %s): %v", original.Method, url, err)
		return nil, err
	}
	req.SetBasicAuth(p.Username, p.Password)
	p.copyHeaders(original.Header, req.Header)
	cli := &http.Client{}
	resp, e = cli.Do(req)
	return
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestUrl := fmt.Sprintf("%s%s", p.ProxyBase, req.RequestURI)
	resp, e := p.MakeProxiedRequest(req, requestUrl)
	if e != nil {
		log.Printf("Got error when requesting url (%s): %v", requestUrl, e)
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
	var user, pass, base string
	var ok bool

	if user, ok = os.LookupEnv("USERNAME"); !ok {
		log.Fatal("ERROR: Must set USERNAME env.")
	}

	if pass, ok = os.LookupEnv("PASSWORD"); !ok {
		log.Fatal("ERROR: Must set PASSWORD env.")
	}

	if base, ok = os.LookupEnv("PROXY_BASE"); !ok {
		log.Fatal("ERROR: Must set PROXY_BASE env.")
	}

	proxy := ProxyHandler{Username: user, Password: pass, ProxyBase: base}
	server := &http.Server{
		Addr:    ":8989",
		Handler: &proxy,
	}

	log.Fatal(server.ListenAndServe())
}
