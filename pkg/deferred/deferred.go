// Package deferred has machinery for providing JavaScript with
// deferred values. These arise because we want to have library
// functions for JavaScript to request things that will only later be
// supplied -- either because it will take some time to get, and it's
// better if JavaScript doesn't have to block (e.g., an HTTP request),
// or because we're waiting for something elsewhere to happen (e.g.,
// watching a resource for changes).
//
// In JavaScript, deferred values will often end up being represented
// by promises; but the protocol allows other representations, and in
// particular, some requests may result in an sequence of values
// rather than just a single value. The protocol also allows for
// cancelling a deferred value.
package deferred

import (
	"context"
	"sync"
)

var (
	globalDeferreds = New()
)

// New returns a pointer to an empty, initialised Deferreds.
func New() *Deferreds {
	return &Deferreds{
		cancels: make(map[Serial]context.CancelFunc),
	}
}

// Register schedules an action to be performed later, with the result
// sent to `resolver`, using the global deferred scheduler.
func Register(p performFunc, r resolver) Serial {
	return globalDeferreds.Register(p, r)
}

// Cancel cancels the fulfilment of a deferred value; see
// `*Deferreds.Cancel` below.
func Cancel(s Serial) {
	globalDeferreds.Cancel(s)
}

// Wait blocks until all outstanding deferred values in the global
// scheduler are fulfilled.
func Wait() {
	globalDeferreds.Wait()
}

// Serial is a serial number used to identify deferreds between Go and
// JavaScript.
type Serial uint64

// Deferreds does the bookkeeping for deferred values.
type Deferreds struct {
	serialMu sync.Mutex
	serial   Serial

	outstanding sync.WaitGroup
	cancelsMu   sync.Mutex
	cancels     map[Serial]context.CancelFunc
}

// responder is the interface for a deferred request to use to send
// its response.
type resolver interface {
	Error(Serial, error)
	Data(Serial, []byte)
	End(Serial)
}

type performFunc func(context.Context) ([]byte, error)

// Register adds a request to those being tracked, and returns the
// serial number to give back to the runtime.
func (d *Deferreds) Register(perform performFunc, r resolver) Serial {
	d.serialMu.Lock()
	s := d.serial
	d.serial++
	d.serialMu.Unlock()
	d.outstanding.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	d.cancelsMu.Lock()
	d.cancels[s] = cancel
	d.cancelsMu.Unlock()
	go func(s Serial) {
		defer func() {
			d.Cancel(s)
			d.outstanding.Done()
		}()
		b, err := perform(ctx)
		if err != nil {
			r.Error(s, err)
			return
		}
		r.Data(s, b)
	}(s)
	return s
}

// Cancel indicates that the caller no longer needs the deferred value
// identified by the serial number. It can be called more than once
// with the same serial number, without panicking. The deferred value
// may still be sent, since there is a race condition between the
// cancellation and the operation completing.
func (d *Deferreds) Cancel(s Serial) {
	d.cancelsMu.Lock()
	c, ok := d.cancels[s]
	delete(d.cancels, s)
	d.cancelsMu.Unlock()
	if ok {
		c()
	}
}

// Wait blocks until all outstanding deferred requests are fulfilled.
func (d *Deferreds) Wait() {
	d.outstanding.Wait()
}
