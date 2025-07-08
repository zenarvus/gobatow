package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func ProxyServer()error{
	handleRequest:=func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			handleHTTPS(w, r)
		} else {
			handleHTTP(w, r)
		}
	}
    server := &http.Server{
        Addr:         ":"+proxyPort,
        Handler:      http.HandlerFunc(handleRequest),
        ReadTimeout:  10 * time.Second,  // Set the read timeout
        WriteTimeout: 10 * time.Second,  // Set the write timeout
        IdleTimeout:  30 * time.Second,  // Set the idle timeout
    }
    log.Println("Starting proxy server on port "+proxyPort)
    log.Fatal(server.ListenAndServe())
    return nil
}

func handleHTTP(w http.ResponseWriter, r *http.Request){
	if isPageBlocked(r.Host) {
		http.Error(w, "Blocked domain", http.StatusForbidden)
		return
	}

	log.Printf("HTTP Request: %s %s", r.Method, r.URL.String())

	// Create a transport with timeouts
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,  // Set the dial timeout
			KeepAlive: 30 * time.Second, // Set the keep-alive period
		}).DialContext,
		TLSHandshakeTimeout: 5 * time.Second,  // Set the TLS handshake timeout
		ResponseHeaderTimeout: 10 * time.Second,  // Set the response header timeout
		ExpectContinueTimeout: 1 * time.Second,  // Set the expect continue timeout
	}

	outReq := new(http.Request)
	*outReq = *r
	outReq.RequestURI = ""

	resp, err := transport.RoundTrip(outReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleHTTPS(w http.ResponseWriter, r *http.Request) {
	if isPageBlocked(r.Host) {
		http.Error(w, "Blocked domain", http.StatusForbidden)
		return
	}

	log.Printf("HTTPS Request: %s %s", r.Method, r.URL.String())

	destConn, err := net.DialTimeout("tcp", r.Host, 5*time.Second)  // Set the dial timeout
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer destConn.Close()

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	// Use go-routines to copy data with timeouts
	done := make(chan struct{})
	go copyDataWithTimeout(destConn, clientConn, done)
	go copyDataWithTimeout(clientConn, destConn, done)
	
	<-done
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func copyDataWithTimeout(dst io.Writer, src io.Reader, done chan struct{}) {
	io.Copy(dst, src)
	done <- struct{}{}
}
