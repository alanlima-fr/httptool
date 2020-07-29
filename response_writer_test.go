// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	_ ResponseWriter      = (*responseWriter)(nil)
	_ http.ResponseWriter = (*responseWriter)(nil)
	_ http.Flusher        = (*responseWriter)(nil)
	_ http.Hijacker       = (*responseWriter)(nil)
	_ http.Pusher         = (*responseWriter)(nil)
)

func TestResponseWriterStatus(t *testing.T) {
	r := httptest.NewRecorder()
	w := NewResponseWriter(r)

	t.Run("ReturnZero", func(t *testing.T) {
		equal(t, 0, w.Status())
		equal(t, http.StatusOK, r.Code)
	})

	t.Run("ReturnOK", func(t *testing.T) {
		w.WriteHeader(http.StatusOK)

		equal(t, http.StatusOK, r.Code)
		equal(t, http.StatusOK, w.Status())
		equal(t, r.Code, w.Status())
	})
}

func TestResponseWriterWritten(t *testing.T) {
	r := httptest.NewRecorder()
	w := NewResponseWriter(r)

	t.Run("ReturnFalse", func(t *testing.T) {
		equal(t, false, w.Written())
	})

	t.Run("ReturnTrue", func(t *testing.T) {
		w.WriteHeader(http.StatusTeapot)
		equal(t, true, w.Written())

		w.WriteHeader(http.StatusContinue)

		equal(t, true, w.Written())
		equal(t, http.StatusTeapot, w.Status())
	})
}

func TestResponseWriterLen(t *testing.T) {
	r := httptest.NewRecorder()
	w := NewResponseWriter(r)

	t.Run("ReturnZero", func(t *testing.T) {
		equal(t, 0, r.Body.Len())
		equal(t, 0, w.Len())
		equal(t, r.Body.Len(), w.Len())
	})

	t.Run("ReturnGreaterThanZero", func(t *testing.T) {
		n, err := w.Write([]byte("Concurrency is not parallelism."))

		equal(t, nil, err)
		equal(t, n, w.Len())
		equal(t, n, r.Body.Len())
		equal(t, r.Body.Len(), w.Len())
	})
}

func TestResponseWriterBody(t *testing.T) {
	r := httptest.NewRecorder()
	w := NewResponseWriter(r)
	s := "Don't communicate by sharing memory, share memory by communicating."

	_, err := w.Write([]byte(s))

	equal(t, nil, err)
	equal(t, s, r.Body.String())
}

func TestResponseWriterWithFlusher(t *testing.T) {
	r := httptest.NewRecorder()
	w := NewResponseWriter(r)
	s := "A little copying is better than a little dependency."

	t.Run("StatusReturnZero", func(t *testing.T) {
		equal(t, 0, w.Status())
	})

	t.Run("FlushReturnFalse", func(t *testing.T) {
		equal(t, false, r.Flushed)
	})

	t.Run("FlushReturnTrue", func(t *testing.T) {
		_, err := r.WriteString(s)

		equal(t, nil, err)
		equal(t, false, r.Flushed)

		w.Flush()
		equal(t, true, r.Flushed)
	})

	t.Run("StatusReturnDefaultOK", func(t *testing.T) {
		equal(t, http.StatusOK, w.Status())
	})
}

type hijackRecorder struct {
	http.ResponseWriter
	hijacked bool
}

func (r *hijackRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	r.hijacked = true
	return nil, nil, nil
}

func newResponseRecorder() *hijackRecorder {
	return &hijackRecorder{}
}

func TestResponseWriterHijacker(t *testing.T) {
	r := newResponseRecorder()
	w := NewResponseWriter(r)

	t.Run("HijackedReturnFalse", func(t *testing.T) {
		equal(t, false, r.hijacked)
	})

	t.Run("HijackedReturnTrue", func(t *testing.T) {
		conn, rw, err := w.Hijack()

		equal(t, true, r.hijacked)
		equal(t, nil, conn)
		equal(t, (*bufio.ReadWriter)(nil), rw)
		equal(t, nil, err)
	})
}

type pusherRecorder struct {
	http.ResponseWriter
	err    error
	target string
	opts   *http.PushOptions
	pushed bool
}

func newPusherRecorder() *pusherRecorder {
	return &pusherRecorder{}
}

func (r *pusherRecorder) Push(target string, opts *http.PushOptions) error {
	r.target = target
	r.opts = opts
	r.pushed = true

	return r.err
}

func TestResponseWriterPush(t *testing.T) {
	r := newPusherRecorder()
	w := NewResponseWriter(r)

	target := "static/css/main.css"
	opts := &http.PushOptions{Method: http.MethodHead}

	t.Run("PushedReturnFalse", func(t *testing.T) {
		equal(t, false, r.pushed)
	})

	t.Run("PushedReturnTrue", func(t *testing.T) {
		err := w.Push(target, opts)

		equal(t, true, r.pushed)
		equal(t, nil, err)
		equal(t, target, r.target)
		equal(t, opts, r.opts)
		equal(t, opts.Method, r.opts.Method)
	})
}

type pusherNotSupportedRecorder struct {
	http.ResponseWriter
}

func newPusherNotSupportedRecorder() pusherNotSupportedRecorder {
	return pusherNotSupportedRecorder{}
}

func (*pusherNotSupportedRecorder) Push(string, *http.PushOptions) error {
	return http.ErrNotSupported
}

func TestResponseWriterPushErrorNotSupported(t *testing.T) {
	r := newPusherNotSupportedRecorder()
	w := NewResponseWriter(r)
	err := w.Push("", nil)

	equal(t, http.ErrNotSupported, err)
}

func TestResponseWriterHandler(t *testing.T) {
	var handler http.Handler

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := w.(ResponseWriter)
		rw.WriteHeader(http.StatusFound)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)

	handler = ResponseHandler(handler)
	handler.ServeHTTP(w, r)

	equal(t, http.StatusFound, w.Code)
}
