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
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nanoninja/httptool"
)

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		rw := w.(httptool.ResponseWriter)
		ip := httptool.ClientIP(r)

		log.Printf("[nanoninja] %s %s %s %d\n", ip, r.Method, r.RequestURI, rw.Status())
	})
}

func handleError(logger *log.Logger, h httptool.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.ServeHTTP(w, r); err != nil {
			logger.Println(err)
		}
	}
}

func index(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("Hello, Gophers"))
	return err
}

func main() {
	logger := log.New(os.Stderr, "", log.Lshortfile)

	mux := http.NewServeMux()
	mux.Handle("/", handleError(logger, index))

	server := httptool.Recovery{
		Next:   httptool.ResponseHandler(logRequest(mux)),
		Logger: logger,
	}

	log.Fatalln(http.ListenAndServe(":3000", server))
}
```

## License

HTTPTool is licensed under the Creative Commons Attribution 3.0 License, and code is licensed under a [BSD license](https://github.com/nanoninja/httptool/blob/master/LICENSE).