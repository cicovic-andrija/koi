package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

var _serverControl control

type control struct {
	bootBlock      sync.Once
	shutdownBlock  sync.Once
	shutdownSignal sync.WaitGroup
	running        bool
	closed         bool
	failure        chan error

	https    *http.Server
	endpoint string
}

func (c *control) boot() {
	c.init()
	c.start()
	c.wait()
}

func (c *control) init() {
	c.https = &http.Server{
		Addr:     c.endpoint,
		Handler:  multiHandler(),
		ErrorLog: log.New(io.Discard, "", 0),
	}
}

func (c *control) start() {
	c.bootBlock.Do(func() {
		c.shutdownSignal.Add(1)
		c.failure = make(chan error)
		go func() {
			err := c.https.ListenAndServe()
			c.shutdownSignal.Done()
			if !errors.Is(err, http.ErrServerClosed) {
				c.failure <- fmt.Errorf("error: server failed unexpectedly: %v", err)
			}
		}()
		c.running = true
		trace(_control, "server started listening on %s", c.https.Addr)
	})
}

func (c *control) wait() {
	c.assertRunning()
	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	trace(_control, "main: waiting indefinitely for interrupt signal or server failure...")
	select {
	case <-interrupts:
		trace(_control, "main: interrupt signal received")
		err := c.shutdown()
		if err != nil {
			panic(err)
		}
		trace(_control, "main: server closed")
	case err := <-c.failure:
		trace(_control, "main: failure signal received")
		panic(err)
	}
}

func (c *control) shutdown() (err error) {
	c.assertRunning()
	c.shutdownBlock.Do(func() {
		if err = c.https.Shutdown(context.TODO()); err != nil {
			err = fmt.Errorf("error: shutdown: %v", err)
		}
		c.shutdownSignal.Wait()
		c.closed = true
	})
	return
}

func (c *control) assertRunning() {
	if !c.running {
		panic("error: server not running")
	}
}
