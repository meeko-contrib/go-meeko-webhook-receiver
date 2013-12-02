/*
   Copyright (C) 2013  Salsita s.r.o.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see {http://www.gnu.org/licenses/}.
*/

package receiver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Make sure that incoming requests are authenticated before it is passed to
// the user-defined request handler.
func TestAuthenticatedServer(t *testing.T) {
	const (
		realm    = "cider-pivotal-tracker"
		username = "pepa"
		password = "yrqvoErxKn"
	)

	var handlerInvoked bool

	userHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerInvoked = true
		w.WriteHeader(http.StatusAccepted)
	})

	handler := authenticatedServer(realm, username, password, userHandler)

	rw := httptest.NewRecorder()

	Convey("Receiving a valid HTTP POST request", t, func() {
		req, err := http.NewRequest("POST", "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetBasicAuth(username, password)

		handler.ServeHTTP(rw, req)

		Convey("The user-defined handler should be invoked", func() {
			So(handlerInvoked, ShouldBeTrue)
			So(rw.Code, ShouldEqual, http.StatusAccepted)
		})
	})
}

// Make sure that unauthenticated requests are refused and the user-defined
// request handler is never invoked.
func TestAuthenticatedServer_Unauthorized(t *testing.T) {
	const (
		realm         = "cider-pivotal-tracker"
		username      = "pepa"
		password      = "yrqvoErxKn"
		incorrectPass = "f656wn6x0t"
	)

	var handlerInvoked bool

	userHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerInvoked = true
		w.WriteHeader(http.StatusAccepted)
	})

	handler := authenticatedServer(realm, username, password, userHandler)

	rw := httptest.NewRecorder()

	Convey("Receiving a valid HTTP POST request", t, func() {
		req, err := http.NewRequest("POST", "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.SetBasicAuth(username, incorrectPass)

		handler.ServeHTTP(rw, req)

		Convey("The user-defined handler should be invoked", func() {
			So(handlerInvoked, ShouldBeFalse)
			So(rw.Code, ShouldEqual, http.StatusUnauthorized)
		})
	})
}
