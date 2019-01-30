/*
 * Session Affinity with Kubernetes
 * Copyright (c) 2019 Pivotal Software, Inc.
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

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

func _HandleRequests(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "/" {
		fmt.Fprintf(w, "backend[%s]\n", _GetHostname())
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
