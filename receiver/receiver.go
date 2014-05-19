// Copyright (c) 2013-2014 The go-meeko-webhook-receiver AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package receiver

import (
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/meeko-contrib/go-meeko-webhook-receiver/receiver/server"

	"github.com/meeko/go-meeko/agent"
)

// API functions ---------------------------------------------------------------

// Serve POST requests using the handler passed into ListenAndServe.
// This function blocks until a signal is received. So signals are being
// handled by this function, no need to do it manually.
func ListenAndServe(handler http.Handler) {
	if err := runListenAndServe(handler); err != nil {
		os.Exit(1)
	}
}

func runListenAndServe(handler http.Handler) error {
	// Make sure agent is terminated properly.
	defer agent.Terminate()

	// Load all the required environment variables, panic if any is not set.
	// This is placed here and not outside to make testing easier (possible).
	// The applications do not have to really connect to Cider to run tests.
	var (
		addr  = os.Getenv("LISTEN_ADDRESS")
		token = os.Getenv("ACCESS_TOKEN")
	)
	switch {
	case addr == "":
		return agent.Logging.Critical("LISTEN_ADDRESS variable is not set")
	case token == "":
		return agent.Logging.Critical("ACCESS_TOKEN variable is not set")
	}

	// Start catching interrupts.
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// Listen.
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return agent.Logging.Critical(err)
	}

	// Start processing interrupts.
	interruptedCh := make(chan bool, 1)
	go func() {
		<-signalCh
		interruptedCh <- true
		listener.Close()
	}()

	// Keep serving until interrupted.
	err = http.Serve(listener, server.AuthenticatedServer(token, handler))
	if err != nil {
		select {
		case <-interruptedCh:
		default:
			return agent.Logging.Critical(err)
		}
	}
	return nil
}
