// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeJSON(t *testing.T) {
	v := struct {
		Name string `json:"name"`
	}{}

	body := `{"name": "Gopher"}`
	r := httptest.NewRequest(http.MethodPost, "http://www.example.com", strings.NewReader(body))

	equal(t, nil, DecodeJSON(r, &v))
	equal(t, "Gopher", v.Name)
}

func TestClientIP(t *testing.T) {
	newRequest := func(remoteAddr, xRealIP string, xForwardedFor ...string) *http.Request {
		h := http.Header{}
		h.Set("X-Real-IP", xRealIP)

		for _, ip := range xForwardedFor {
			h.Set("X-Forwarded-For", ip)
		}

		return &http.Request{
			RemoteAddr: remoteAddr,
			Header:     h,
		}
	}

	addr1 := "134.26.30.10"
	addr2 := "120.13.64.5"
	addr3 := "120.13.64.5:8080"
	local := "127.0.0.0"

	tests := map[string]struct {
		request  *http.Request
		expected net.IP
	}{
		"NoHeader": {
			request:  newRequest(addr1, ""),
			expected: net.ParseIP(addr1),
		},
		"NoHeaderWithPort": {
			request:  newRequest(addr3, ""),
			expected: net.ParseIP(addr2),
		},
		"NoHeaderWithNoValidRemoteAddr": {
			request:  newRequest("127.0.0.0::", ""),
			expected: nil,
		},
		"XForwardedFor": {
			request:  newRequest("", "", addr1),
			expected: net.ParseIP(addr1),
		},
		"MultipleXForwardedFor": {
			request:  newRequest("", "", local, addr1, addr2),
			expected: net.ParseIP(addr2),
		},
		"MultipleXForwardedForFromString": {
			request: newRequest(
				"",
				"",
				fmt.Sprintf("%s, %s, %s, %s", "", addr1, local, addr2)),
			expected: net.ParseIP(addr1),
		},
		"MultipleEmptyXForwardedFor": {
			request:  newRequest(addr1, "", ", , , , ,"),
			expected: net.ParseIP(addr1),
		},
		"XRealIP": {
			request:  newRequest("", addr1),
			expected: net.ParseIP(addr1),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ip := ClientIP(tt.request)

			equal(t, true, ip.Equal(tt.expected))
			equal(t, tt.expected.String(), ip.String())
		})
	}
}

func TestIsSecure(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)

	t.Run("ReturnFalse", func(t *testing.T) {
		if actual := IsSecure(r); actual != false {
			equal(t, false, actual)
		}
	})

	t.Run("ReturnTrue", func(t *testing.T) {
		r.TLS = &tls.ConnectionState{}

		if actual := IsSecure(r); actual != true {
			equal(t, true, actual)
		}
	})
}

func TestIsXMLHTTPRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://www.example.com", nil)

	t.Run("ReturnFalse", func(t *testing.T) {
		if actual := IsXMLHTTPRequest(r); actual != false {
			equal(t, false, actual)
		}
	})

	t.Run("ReturnTrue", func(t *testing.T) {
		r.Header.Set("X-Requested-With", "XMLHttpRequest")

		if actual := IsXMLHTTPRequest(r); actual != true {
			equal(t, true, actual)
		}
	})
}
