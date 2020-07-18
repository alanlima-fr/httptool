// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func errorf(t *testing.T, expected, actual interface{}) {
	t.Errorf("\nexpected: %v\n  actual: %v", expected, actual)
}

func equal(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		errorf(t, expected, actual)
	}
}

func TestError(t *testing.T) {
	e := Error("test error")
	equal(t, "test error", e.Error())
}

func TestHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)

	handler := HandlerFunc(func(w http.ResponseWriter, _ *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	})

	err := handler.ServeHTTP(w, r)

	equal(t, nil, err)
	equal(t, http.StatusOK, w.Code)
}

func TestChain(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)

	handler := func(next Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			return next.ServeHTTP(w, r)
		})
	}

	t.Run("WithNoError", func(t *testing.T) {
		h := func(w http.ResponseWriter, _ *http.Request) error {
			w.WriteHeader(http.StatusInternalServerError)
			return nil
		}

		c := NewChain(handler)
		err := c.ThenFunc(h).ServeHTTP(w, r)

		equal(t, nil, err)
		equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("WithErrNoHandlerProvided", func(t *testing.T) {
		c := NewChain(handler)
		err := c.Then(nil).ServeHTTP(w, r)

		equal(t, ErrNoHandlerProvided, err)
	})
}
