// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nanoninja/httptool"
)

func ExampleNewChain() {
	middleware := func(next httptool.Handler) httptool.Handler {
		return httptool.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			fmt.Println("Middleware")
			return next.ServeHTTP(w, r)
		})
	}

	index := func(w http.ResponseWriter, _ *http.Request) error {
		_, err := w.Write([]byte("Hello, Gophers"))
		return err
	}

	c := httptool.NewChain(middleware).ThenFunc(index)

	server := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := c.ServeHTTP(w, r); err != nil {
			fmt.Println("Error:", err.Error())
		}
	})

	log.Fatalln(http.ListenAndServe(":3000", server))
}
