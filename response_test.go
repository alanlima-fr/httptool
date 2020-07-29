// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncodeJSON(t *testing.T) {
	w := httptest.NewRecorder()
	v := struct {
		Name string `json:"name"`
	}{
		Name: "Gopher",
	}

	equal(t, nil, EncodeJSON(w, &v))
	equal(t, http.StatusOK, w.Code)
	equal(t, w.Body.String(), `{"name":"Gopher"}`+"\n")
}
