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

var (
	_ Logger       = (*loggerTest)(nil)
	_ http.Handler = (*Recovery)(nil)
)

type loggerTest struct {
	logged bool
	value  string
}

func (l *loggerTest) Printf(format string, v ...interface{}) {
	l.logged = true
	l.value = fmt.Sprintf(format, v...)
}

func handlerTest(http.ResponseWriter, *http.Request) {
	panic("errors are values")
}

func TestRecovery(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)

	logger := &loggerTest{}
	recovery := Recovery{Next: http.HandlerFunc(handlerTest), Logger: logger}

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
