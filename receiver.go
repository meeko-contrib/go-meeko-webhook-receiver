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
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/cider/go-cider/cider/services/logging"
	"github.com/cider/go-cider/cider/services/pubsub"
	zlogging "github.com/cider/go-cider/cider/transports/zmq3/logging"
	zpubsub "github.com/cider/go-cider/cider/transports/zmq3/pubsub"

	zmq "github.com/pebbe/zmq3"
)

// API functions ---------------------------------------------------------------

var (
	Logger *logging.Service
	PubSub *pubsub.Service
)

// Serve POST requests using the handler passed into ListenAndServe.
// This function blocks until a signal is received. So signals are being
// handled by this function, no need to do it manually.
func ListenAndServe(handler http.Handler) {
	// Load all the required environment variables, panic if any is not set.
	// This is placed here and not outside to make testing easier (possible).
	// The applications do not have to really connect to Cider to run tests.
	var (
		alias = mustBeSet(os.Getenv("CIDER_ALIAS"))
		addr  = mustBeSet(os.Getenv("LISTEN_ADDRESS"))
		token = mustBeSet(os.Getenv("ACCESS_TOKEN"))
	)

	// Start catching interrupts.
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// Initialise Logging service from environmental variables.
	var err error
	Logger, err = logging.NewService(func() (logging.Transport, error) {
		factory := zlogging.NewTransportFactory()
		factory.MustReadConfigFromEnv("CIDER_ZMQ3_LOGGING_").MustBeFullyConfigured()
		return factory.NewTransport(alias)
	})
	if err != nil {
		panic(err)
	}
	Logger.Info("Logging service initialised\n")

	// Make sure ZeroMQ is terminated properly.
	defer func() {
		Logger.Info("Waiting for ZeroMQ context to terminate...\n")
		Logger.Flush()
		zmq.Term()
	}()

	// Initialise PubSub service from environmental variables.
	PubSub, err = pubsub.NewService(func() (pubsub.Transport, error) {
		factory := zpubsub.NewTransportFactory()
		factory.MustReadConfigFromEnv("CIDER_ZMQ3_PUBSUB_").MustBeFullyConfigured()
		return factory.NewTransport(alias)
	})
	if err != nil {
		panic(Logger.Critical(err))
	}
	defer PubSub.Close()
	Logger.Info("PubSub service initialised\n")

	// Listen.
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(Logger.Critical(err))
	}

	// Start processing interrupts.
	var interrupted bool
	go func() {
		<-signalCh
		interrupted = true
		listener.Close()
	}()

	// Keep serving until interrupted.
	err = http.Serve(listener, authenticatedServer(token, handler))
	if err != nil && !interrupted {
		panic(Logger.Critical(err))
	}
}

// Helpers ---------------------------------------------------------------------

func mustBeSet(v string) string {
	if v == "" {
		panic("Required variable is not set")
	}
	return v
}

func authenticatedServer(token string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make sure that the token query parameter is set correctly.
		if r.FormValue("token") != token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Allow the POST method only.
		if r.Method != "POST" {
			http.Error(w, "POST Method Expected", http.StatusMethodNotAllowed)
			return
		}

		// If everything is ok, serve the user-defined handler.
		handler.ServeHTTP(w, r)
	})
}
