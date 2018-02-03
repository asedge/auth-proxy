package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestProxyHandlerReturnsStatusCodeFromServer(t *testing.T) {
	expectedStatusCode := 404

	// Set up fake remote server to proxy to.
	remoteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(expectedStatusCode)
		fmt.Fprintln(w, "Not found")
	}))
	defer remoteServer.Close()

	// Set up proxy server.
	handler := &ProxyHandler{Username: "joe", Password: "secret", ProxyBase: remoteServer.URL}
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	// Request setup & execution
	url := fmt.Sprintf("%s/foo", proxyServer.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		t.Fatalf("Could not make request: %v", req)
	}

	// Assertions
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("Expected Status Code to be %v but got %v", expectedStatusCode, resp.StatusCode)
	}
}

func TestProxyHandlerReturnsHeadersFromServer(t *testing.T) {
	headerMap := map[string][]string{
		"foo": {"bar"},
	}

	// Set up fake remote server to proxy to.
	remoteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for k, values := range headerMap {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
		fmt.Fprintln(w, "Hello, client")
	}))
	defer remoteServer.Close()

	// Set up proxy server.
	handler := &ProxyHandler{Username: "joe", Password: "secret", ProxyBase: remoteServer.URL}
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	// Request setup & execution
	url := fmt.Sprintf("%s/foo", proxyServer.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	resp, e := client.Do(req)
	if e != nil {
		t.Fatalf("Could not make request: %v", req)
	}

	// Assertions
	if !reflect.DeepEqual(resp.Header["Foo"], headerMap["foo"]) {
		t.Errorf("Expected returned header 'foo' to contain %s but got %s.", headerMap["foo"], resp.Header["Foo"])
	}
}

func TestProxyHandlerForwardsHeadersFromClient(t *testing.T) {
	var actualHeader http.Header
	headerMap := map[string][]string{
		"foo": {"bar"},
	}

	// Set up fake remote server to proxy to.
	remoteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualHeader = r.Header
		fmt.Fprintln(w, "Hello, client")
	}))
	defer remoteServer.Close()

	// Set up proxy server.
	handler := &ProxyHandler{Username: "joe", Password: "secret", ProxyBase: remoteServer.URL}
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	// Request setup & execution
	url := fmt.Sprintf("%s/foo", proxyServer.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	for k, values := range headerMap {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}
	client := &http.Client{}
	_, e := client.Do(req)
	if e != nil {
		t.Fatalf("Could not make request: %v", req)
	}

	// Assertions
	if !reflect.DeepEqual(actualHeader["Foo"], headerMap["foo"]) {
		t.Errorf("Expected header 'foo' to contain %s but got %s.", headerMap["foo"], actualHeader["Foo"])
	}
}

func TestProxyHandlerUsesRequestMethod(t *testing.T) {
	expectedMethod := "POST"
	var actualMethod string

	// Set up fake remote server to proxy to.
	remoteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		actualMethod = r.Method
		fmt.Fprintln(w, "Hello, client")
	}))
	defer remoteServer.Close()

	// Set up proxy server.
	handler := &ProxyHandler{Username: "joe", Password: "secret", ProxyBase: remoteServer.URL}
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	// Request setup & execution
	url := fmt.Sprintf("%s/foo", proxyServer.URL)
	req, err := http.NewRequest(expectedMethod, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{}
	_, e := client.Do(req)
	if e != nil {
		t.Fatalf("Could not make request: %v", req)
	}

	// Assertions
	if actualMethod != expectedMethod {
		t.Errorf("Request used method %s but server proxied method %s", expectedMethod, actualMethod)
	}
}
