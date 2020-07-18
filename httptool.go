// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import "net/http"

// ErrNoHandlerProvided is used when no handler
// is given as an argument to httptool.Then and httptool.ThenFunc.
// A handler will be used by default returning this error.
const ErrNoHandlerProvided = Error("no handler was provided")

// Logger is an interface representing the ability to log error.
type Logger interface {

	// Printf is used for print to the logger.
	// Arguments are handled in the manner of fmt.Printf.
	Printf(format string, v ...interface{})
}

// Error is a trivial implementation of error.
type Error string

func (e Error) Error() string { return string(e) }

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

// Chain is a stack of handlers.
type Chain struct {
	handlers []func(Handler) Handler
}

// NewChain creates a new chain.
func NewChain(m ...func(Handler) Handler) Chain {
	return Chain{}.Use(m...)
}

// Use adds handler onto the chain stack.
func (c Chain) Use(handler ...func(Handler) Handler) Chain {
	handlers := make([]func(Handler) Handler, 0, len(c.handlers)+len(handler))
	handlers = append(handlers, c.handlers...)
	handlers = append(handlers, handler...)
	c.handlers = handlers

	return c
}

// Then chains the handlers and returns the final httptool.Handler.
// ThenFunc is a wrap of Then, but takes a HandlerFunc instead of a Handler.
//
//     h := httptool.NewChain(f1, f2, f3).Then(handler)
//
func (c Chain) Then(next Handler) Handler {
	if next == nil {
		next = HandlerFunc(func(http.ResponseWriter, *http.Request) error {
			return ErrNoHandlerProvided
		})
	}

	if n := len(c.handlers); n >= 1 {
		for i := n - 1; i >= 0; i-- {
			next = c.handlers[i](next)
		}
	}

	return next
}

// ThenFunc is a wrap of httptool.Then.
func (c Chain) ThenFunc(next HandlerFunc) Handler {
	return c.Then(next)
}
