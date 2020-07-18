// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import "net/http"

// Recovery is an http.Handler that recovers from all panics.
type Recovery struct {
	Next   http.Handler
	Logger Logger
}

func (r Recovery) log(format string, v ...interface{}) {
	if r.Logger != nil {
		r.Logger.Printf(format, v...)
	}
}

func (r Recovery) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			r.log("%v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()

	r.Next.ServeHTTP(w, req)
}
