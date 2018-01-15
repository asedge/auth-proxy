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

func (p *ProxyHandler) MakeProxiedRequest(method, url string, headers http.Header) (resp *http.Response, e error) {
	req, _ := http.NewRequest(method, url, nil)
	req.SetBasicAuth(p.Username, p.Password)
	for k, v := range headers {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	//log.Printf("Request headers: %v", req.Header)
	cli := &http.Client{}
	resp, e = cli.Do(req)
	return
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestUrl := fmt.Sprintf("%s%s", p.ProxyBase, r.RequestURI)
	log.Printf("Requesting URL (%s) for client (%s) with headers (%v).", requestUrl, r.RemoteAddr, r.Header)
	//log.Printf("Requesting URL (%s) for client (%s).", requestUrl, r.RemoteAddr)
	resp, e := p.MakeProxiedRequest(r.Method, requestUrl, r.Header)
	//log.Printf("Received headers: %v", resp.Header)
	if e != nil {
		log.Printf("Got error when requesting url (%s): %v", requestUrl, e)
	}
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Printf("Unable to read response body: %v", readErr)
	}
	resp.Body.Close()
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.Write(body)
}
