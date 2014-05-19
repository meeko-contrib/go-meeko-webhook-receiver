// Copyright (c) 2013-2014 The go-meeko-webhook-receiver AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package server

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
		token = "yrqvoErxKn"
	)

	var handlerInvoked bool

	userHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerInvoked = true
		w.WriteHeader(http.StatusAccepted)
	})

	handler := AuthenticatedServer(token, userHandler)

	rw := httptest.NewRecorder()

	Convey("Receiving a valid HTTP POST request", t, func() {
		req, err := http.NewRequest("POST", "http://example.com?token="+token, nil)
		if err != nil {
			t.Fatal(err)
		}

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
		token          = "yrqvoErxKn"
		incorrectToken = "f656wn6x0t"
	)

	var handlerInvoked bool

	userHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerInvoked = true
		w.WriteHeader(http.StatusAccepted)
	})

	handler := AuthenticatedServer(token, userHandler)

	rw := httptest.NewRecorder()

	Convey("Receiving a HTTP POST request with invalid token", t, func() {
		req, err := http.NewRequest("POST", "http://example.com?token="+incorrectToken, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler.ServeHTTP(rw, req)

		Convey("The request should be rejected with 401 Unauthorized", func() {
			So(handlerInvoked, ShouldBeFalse)
			So(rw.Code, ShouldEqual, http.StatusUnauthorized)
		})
	})
}
