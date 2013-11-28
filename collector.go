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

package collector

import (
	"crypto/sha1"
	"encoding/base64"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	auth "github.com/abbot/go-http-auth"

	"cider.go"
	_ "cider.go/dialers/zmq"
)

// Internal Cider session.
var session cider.Session

// Serve POST requests using the handler passed into ListenAndServe.
// This function blocks until a signal is received. So signals are being
// handled by this function, no need to do it manually.
func ListenAndServe(handler http.HandlerFunc) {
	// Load all the required environment variables, panic if any is not set.
	var (
		addr     = mustBeSet(os.Getenv("HTTP_ADDR"))
		realm    = mustBeSet(os.Getenv("HTTP_AUTH_REALM"))
		username = mustBeSet(os.Getenv("HTTP_AUTH_USERNAME"))
		password = mustBeSet(os.Getenv("HTTP_AUTH_PASSWORD"))
	)

	// Compute SHA1 of password since that is what go-http-auth expects.
	hasher := sha1.New()
	hasher.Write([]byte(password))
	password = "{SHA}" + base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Initialise a Cider session.
	dialer := cider.MustNewDialer("zmq", nil)
	session = dialer.MustDial(cider.MustSessionConfigFromEnv())
	defer session.Close()

	// Prepare the internal HTTP request handler.
	secretFunc := func(user, realm string) string {
		if user == username {
			return password
		}
		return ""
	}

	authenticator := auth.NewBasicAuthenticator(realm, secretFunc)
	authHandler := authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		// Allow the POST method only.
		if r.Method != "POST" {
			http.Error(w, "POST Method Expected", http.StatusMethodNotAllowed)
			return
		}
		handler(w, &r.Request)
	})

	// Listen.
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// Start catching signals.
	var (
		interrupted bool

		closeCh    = make(chan struct{})
		closeAckCh = make(chan struct{})
	)

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ch:
			interrupted = true
			listener.Close()
		case <-closeCh:
			close(closeAckCh)
		}
	}()

	defer func() {
		close(closeCh)
		<-closeAckCh
	}()

	// Keep serving until interrupted.
	err = http.Serve(listener, authHandler)
	if err != nil && !interrupted {
		panic(err)
	}
}

// eventBody must be marshallable by encoding/json and the codec package.
func Publish(eventType string, eventBody interface{}) error {
	return session.Publish(eventType, eventBody)
}

// Helpers ---------------------------------------------------------------------

func mustBeSet(v string) string {
	if v == "" {
		panic("Required variable is not set")
	}
	return v
}
