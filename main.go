package main

import {
  "crypto/tls"
  "crypto/x509"
  "fmt"
  "io"
  "log"
  "net/http"
  "os"
}

const (
  k8sAPI = "https://kubernetes.default.svc"
  openapiPath = "/openapi/v2"
  tokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
  caCertPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
)

func openapiHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodOptions {
    addCors(w)
    w.WriteHeader(http.StatusNoContent)
    return
    }
