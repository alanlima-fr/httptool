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

func TestChainFunc(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)

	middleware := func(next Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			return next.ServeHTTP(w, r)
		})
	}

	h := func(w http.ResponseWriter, _ *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	mw := ChainFunc(h, middleware)
	err := mw.ServeHTTP(w, r)

	equal(t, nil, err)
	equal(t, http.StatusOK, w.Code)
}
