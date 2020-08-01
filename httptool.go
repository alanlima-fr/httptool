// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import "net/http"

// Logger is an interface representing the ability to log error.
type Logger interface {

	// Printf is used for print to the logger.
	// Arguments are handled in the manner of fmt.Printf.
	Printf(format string, v ...interface{})
}

// A Handler responds to an HTTP request.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as handlers.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// Middleware describes a middleware that can be applied to a httptool.handler.
type Middleware func(Handler) Handler

// Chain creates a new handler by wrapping middleware around a final httptool.Handler.
func Chain(next Handler, mw ...Middleware) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		if h := mw[i]; h != nil {
			next = h(next)
		}
	}

	return next
}

// Chain is wrap of httptool.Chain that can use an httptool.HandlerFunc.
func ChainFunc(next HandlerFunc, mw ...Middleware) Handler {
	return Chain(next, mw...)
}
