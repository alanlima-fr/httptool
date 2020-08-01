// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

type loggerTest struct {
	logged bool
	value  string
}

func (l *loggerTest) Printf(format string, v ...interface{}) {
	l.logged = true
	l.value = fmt.Sprintf(format, v...)
}

func TestRecoveryHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)

	handlerTest := func(w http.ResponseWriter, r *http.Request) {
		panic("errors are values")
	}

	logger := &loggerTest{}
	recovery := RecoveryHandler(http.HandlerFunc(handlerTest), logger)

	t.Run("NotLoggedPanic", func(t *testing.T) {
		equal(t, false, logger.logged)
		equal(t, "", logger.value)
	})

	recovery.ServeHTTP(w, r)

	t.Run("LoggedPanic", func(t *testing.T) {
		equal(t, true, logger.logged)
		equal(t, "errors are values", logger.value)
		equal(t, http.StatusInternalServerError, w.Code)
		equal(t, true, strings.Contains(w.Body.String(), http.StatusText(http.StatusInternalServerError)))
	})
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
