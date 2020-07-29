// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import (
	"bufio"
	"net"
	"net/http"
)

// A ResponseWriter interface is used by an HTTP handler to construct an HTTP response.
// It is an extension of ResponseWriter interface of the net/http package.
type ResponseWriter interface {
	// A ResponseWriter interface is used by an HTTP handler to
	// construct an HTTP response.
	// see https://golang.org/pkg/net/http/#ResponseWriter
	http.ResponseWriter

	// The Flusher interface is implemented by ResponseWriters that allow
	// an HTTP handler to flush buffered data to the client.
	// see https://golang.org/pkg/net/http/#Flusher
	http.Flusher

	// The Hijacker interface is implemented by ResponseWriters that allow
	// an HTTP handler to take over the connection.
	// see https://golang.org/pkg/net/http/#Hijacker
	http.Hijacker

	// Pusher is the interface implemented by ResponseWriters that support
	// HTTP/2 server push. For more background,
	// see https://golang.org/pkg/net/http/#Pusher
	http.Pusher

	// Len returns the number of bytes of the unread portion of the buffer.
	Len() int

	// Status returns the status code of the response or 0 if the response has
	// not been written.
	Status() int

	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
}

type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	length      int
	wroteHeader bool
}

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter.
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     0,
		wroteHeader:    false,
	}
}

// Flush sends any buffered data to the client.
func (w *responseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		if !w.Written() {
			w.WriteHeader(http.StatusOK)
		}

		flusher.Flush()
	}
}

// Hijack lets the caller take over the connection.
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// Push initiates an HTTP/2 server push.
func (w *responseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := w.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}

	return http.ErrNotSupported
}

// Len returns the number of bytes of the unread portion of the buffer.
func (w *responseWriter) Len() int {
	return w.length
}

// Status returns the status code of the response or 0 if the response has
// not been written.
func (w *responseWriter) Status() int {
	return w.statusCode
}

// Write writes the headers described in h to w.
func (w *responseWriter) Write(p []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(p)
	w.length += n
	return
}

// Written returns whether or not the ResponseWriter has been written.
func (w *responseWriter) Written() bool {
	return w.wroteHeader
}

// WriteHeader sends an HTTP response header with the provided
// status code.
func (w *responseWriter) WriteHeader(statusCode int) {
	if w.Written() {
		return
	}

	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
	w.wroteHeader = true
}
