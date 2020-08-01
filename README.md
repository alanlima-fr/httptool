# HTTPTool

The HTTP tool is a simple extension of the net/http package written in Go.
It is not a Framework and offers additional. This package provides a slightly different handling of the Handler type by
exploiting the return to handle potential errors.

The ResponseWriter type has been increased to retrieve the state of the HTTP response such as 
code, length and whether the header has been written.

[![Golang](https://img.shields.io/badge/go-lang-%2347cafa.svg)](https://golang.org/) 
[![Godoc](https://godoc.org/github.com/nanoninja/httptool?status.svg)](https://pkg.go.dev/github.com/nanoninja/httptool?tab=doc) 
[![Build Status](https://travis-ci.org/nanoninja/httptool.svg)](https://travis-ci.org/nanoninja/httptool) 
[![Coverage Status](https://coveralls.io/repos/github/nanoninja/httptool/badge.svg?branch=master)](https://coveralls.io/github/nanoninja/httptool?branch=master) 
[![Go Report Card](https://goreportcard.com/badge/github.com/nanoninja/httptool)](https://goreportcard.com/report/github.com/nanoninja/httptool)
[![codebeat badge](https://codebeat.co/badges/0ce06064-931b-41ba-b29e-dcfbb6c577f3)](https://codebeat.co/projects/github-com-nanoninja-httptool-master)

## Installation

```shell script
go get github.com/nanoninja/httptool
```

## Getting Started

```go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nanoninja/httptool"
)

func logAccess(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		ip := httptool.ClientIP(r)
		rw := w.(httptool.ResponseWriter)

		logger.Printf("[nanoninja] %s %s %s %d\n", ip, r.Method, r.RequestURI, rw.Status())
	})
}

func logError(next httptool.HandlerFunc, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next.ServeHTTP(w, r); err != nil {
			logger.Println(err)
		}
	}
}

func greeting(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("Hello, Gophers"))
	return err
}

func main() {
    errLogger := log.New(os.Stderr, "", log.Lshortfile)
    accessLogger := log.New(os.Stdout, "", log.Lshortfile)
 
    mux := http.NewServeMux()
    mux.Handle("/", logError(greeting, errLogger))
 
    handler := logAccess(mux, accessLogger)
    handler = httptool.ResponseHandler(handler)
    handler = httptool.RecoveryHandler(handler, errLogger)
    
    log.Fatalln(http.ListenAndServe(":3000", handler))
}
```

## License

HTTPTool is licensed under the Creative Commons Attribution 3.0 License, and code is licensed under a [BSD license](https://github.com/nanoninja/httptool/blob/master/LICENSE).