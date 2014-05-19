// Copyright (c) 2013-2014 The go-meeko-webhook-receiver AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package server

import "net/http"

// AuthenticatedServer turns a regular http.Handler into a handler that returns
// Unauthorized in case the token query parameter is not set correctly.
func AuthenticatedServer(token string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow the POST method only.
		if r.Method != "POST" {
			http.Error(w, "POST Method Expected", http.StatusMethodNotAllowed)
			return
		}

		// Make sure that the token query parameter is set correctly.
		if r.FormValue("token") != token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// If everything is ok, serve the user-defined handler.
		handler.ServeHTTP(w, r)
	})
}
