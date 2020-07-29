// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import (
	"encoding/json"
	"net/http"
)

// EncodeJSON writes the JSON document of v to the body response.
func EncodeJSON(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}
