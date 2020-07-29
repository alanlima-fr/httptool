// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import (
	"encoding/json"
	"net/http"
)

// ResponseHandler is a http.handler that extends the http.ResponseWriter written in Go.
func ResponseHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(NewResponseWriter(w), r)
	})
}

// EncodeJSON writes the JSON document of v to the body response.
func EncodeJSON(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}
