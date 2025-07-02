package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	k8sAPI      = "https://kubernetes.default.svc"
	openapiPath = "/openapi/v2"
	tokenPath   = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	caCertPath  = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

func openapiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		addCORS(w)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	token, err := os.ReadFile(tokenPath)
	if err != nil {
		http.Error(w, "cannot read SA token", http.StatusInternalServerError)
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: getTLSConfig(caCertPath),
		},
	}

	req, _ := http.NewRequest("GET", k8sAPI+openapiPath, nil)
	req.Header.Set("Authorization", "Bearer "+string(token))

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "upstream error: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	addCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, resp.Body)
}

func addCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
