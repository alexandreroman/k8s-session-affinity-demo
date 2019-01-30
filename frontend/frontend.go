/*
 * Session Affinity with Kubernetes
 * Copyright (c) 2019 Pivotal Software, Inc.
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var _tr = &http.Transport{
	// Disable connection cache: make sure a new connection is used
	// each time we try to connect to the backend.
	DisableKeepAlives: true,
}
var _client = &http.Client{Transport: _tr}

func _GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Cannot get hostname", err)
	}
	return hostname
}

func _GetServerPort() int {
	port := os.Getenv("PORT")
	if port != "" {
		intPort, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal("Cannot parse environment variable: PORT")
		}
		return intPort
	}
	return 8080
}

func _GetBackend() string {
	backendHost := os.Getenv("BACKEND_HOST")
	if backendHost == "" {
		backendHost = "localhost"
	}
	backendPort := os.Getenv("BACKEND_PORT")
	if backendPort == "" {
		backendPort = "8080"
	}
	return fmt.Sprintf("%s:%s", backendHost, backendPort)
}

func _HandleRequests(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "/" {
		resp, err := _client.Get(fmt.Sprintf("http://%s", _GetBackend()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Backend not available: %s\n", err)
			return
		}
		defer resp.Body.Close()
		msg, err := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "frontend[%s] -> %s\n", _GetHostname(), msg)
	}
}

func main() {
	log.Printf("Process PID: %d", os.Getpid())

	port := _GetServerPort()
	log.Printf("Listening on port: %d", port)

	http.HandleFunc("/", _HandleRequests)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("Cannot start web server", err)
	}
}
