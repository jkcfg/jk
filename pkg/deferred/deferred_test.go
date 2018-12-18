package deferred

import (
	"context"
	"testing"
	"time"
)

type menuFunc func([]byte, error, bool)

func (fn menuFunc) Error(s Serial, err error) {
	fn(nil, err, false)
}

func (fn menuFunc) Data(s Serial, data []byte) {
	fn(data, nil, false)
}

func (fn menuFunc) End(s Serial) {
	fn(nil, nil, true)
}

// TestCancel checks that cancelling a deferred resolves it as an
// error.
func TestCancel(t *testing.T) {
	ds := New()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	ser := ds.RegisterWithContext(ctx, func(ctx context.Context) ([]byte, error) {
		select {
		case <-time.After(1 * time.Second):
			break
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		return []byte{}, nil
	}, menuFunc(func(data []byte, err error, end bool) {
		if data != nil {
			t.Errorf("data is not nil, but all I did was cancel")
		}
		if end {
			t.Errorf("end is called unexpectedly")
		}
		if err == nil {
			t.Fatal("Error not given, but expected one")
		}
	}))
	ds.Cancel(ser)
	ds.Wait()
}

// TestTimeout checks that a deferred resolve can time out, and this
// results in an error.
func TestTimeout(t *testing.T) {
	ds := New()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	ds.RegisterWithContext(ctx, func(ctx context.Context) ([]byte, error) {
		select {
		case <-time.After(1 * time.Second):
			break
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		return []byte{}, nil
	}, menuFunc(func(data []byte, err error, end bool) {
		if data != nil {
			t.Errorf("data is not nil, but all I did was cancel")
		}
		if end {
			t.Errorf("end is called unexpectedly")
		}
		if err == nil {
			t.Fatal("Error not given, but expected one")
		}
	}))
	// don't explicitly cancel; wait for it to time out
	ds.Wait()
}
