// Copyright 2020 The Nanoninja Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httptool

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

// ClientIP returns the client's IP address.
func ClientIP(r *http.Request) net.IP {
	for _, addr := range []string{"X-Real-IP", "X-Forwarded-For"} {
		if ip := ParseIP(r.Header.Get(addr)); ip != nil {
			return ip
		}
	}

	for _, addr := range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		if ip := ParseIP(addr); ip != nil {
			return ip
		}
	}

	return ParseIP(r.RemoteAddr)
}

// DecodeJSON reads the body of an HTTP request looking for a JSON document.
func DecodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// IsSecure returns whether is https secure request.
func IsSecure(r *http.Request) bool {
	return r.TLS != nil
}

// IsXMLHTTPRequest returns whether is the request a Javascript XMLHttpRequest.
func IsXMLHTTPRequest(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// ParseIP checks if the IP is valid.
func ParseIP(ip string) net.IP {
	ip = strings.TrimSpace(ip)

	if ip == "" {
		return nil
	}

	if strings.ContainsRune(ip, ':') {
		ip, _, _ = net.SplitHostPort(ip)
	}

	return net.ParseIP(ip)
}
