package git

// fanout taken and adapted from
// https://github.com/voxelbrain/pixelpixel/blob/master/fanout.go

import "sync"

// fanout is a small structure implementing a more or less generic,
// thread-safe fanout. Fanouts are created on input channels and
// propagate each received value to all consumers in order.
// Consumers can close their channels.
type fanout struct {
	*sync.RWMutex
	consumers map[chan string]struct{}
	closing   map[chan string]bool
}

// newFanout creates a new fanout from a channel.
func newFanout(c <-chan string) *fanout {
	f := &fanout{
		RWMutex:   &sync.RWMutex{},
		consumers: map[chan string]struct{}{},
		closing:   map[chan string]bool{},
	}
	go f.loop(c)
	return f
}

// Output creates a new consumer output.
func (f *fanout) Output() <-chan string {
	c := make(chan string)
	f.Lock()
	defer f.Unlock()
	f.consumers[c] = struct{}{}
	return c
}

// Close a consumer channel, stopping propagation for this particular
// consumer.
func (f *fanout) Close(rc <-chan string) {
	f.RLock()
	defer f.RUnlock()

	// Lookup original channel because we can't call close()
	// on a receive-only channel
	var c chan string
	for i := range f.consumers {
		if i == rc {
			c = i
		}
	}

	// If channel is not in consumers map are is already about to close
	// don't try to do it again.
	if _, ok := f.closing[c]; c == nil || ok {
		return
	}
	f.closing[c] = true

	// Wait for the current broadcast to finish (effectively unlocking
	// the mutex) and delete the consumer from the map.
	go func() {
		f.Lock()
		defer f.Unlock()
		delete(f.consumers, c)
		delete(f.closing, c)
		close(c)
	}()

	// Eat the values possibly left in channel in case the consumer
	// doesn't.
	go func() {
		for {
			_, ok := <-c
			// If the channel is closed it has been removed from
			// the consumers map by the previous goroutine. Stop eating.
			if !ok {
				return
			}
		}
	}()
}

func (f *fanout) loop(c <-chan string) {
	for v := range c {
		f.broadcast(v)
	}
	f.closeConsumers()
}

func (f *fanout) closeConsumers() {
	f.RLock()
	defer f.RUnlock()
	for c := range f.consumers {
		f.Close(c)
	}
}

func (f *fanout) broadcast(v string) {
	f.RLock()
	defer f.RUnlock()
	for c := range f.consumers {
		c <- v
	}
}
