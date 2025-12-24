package conveyer

import (
	"context"
	"errors"
	"sync"
)

var errChanNotFound = errors.New("chan not found")

const undefinedData = "undefined"

type impl struct {
	size int

	mu       sync.RWMutex
	channels map[string]chan string
	workers  []func(ctx context.Context) error

	running bool
}

func newImpl(size int) *impl {
	if size < 0 {
		size = 0
	}
	return &impl{
		size:     size,
		channels: make(map[string]chan string),
		workers:  make([]func(ctx context.Context) error, 0),
	}
}

func (c *impl) ensureChan(name string) chan string {
	if ch, ok := c.channels[name]; ok {
		return ch
	}
	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *impl) RegisterDecorator(fn decoratorFunc, input string, output string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	in := c.ensureChan(input)
	out := c.ensureChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *impl) RegisterMultiplexer(fn multiplexerFunc, inputs []string, output string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inChans := make([]chan string, 0, len(inputs))
	for _, name := range inputs {
		inChans = append(inChans, c.ensureChan(name))
	}
	out := c.ensureChan(output)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, inChans, out)
	})
}

func (c *impl) RegisterSeparator(fn separatorFunc, input string, outputs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	in := c.ensureChan(input)

	outChans := make([]chan string, 0, len(outputs))
	for _, name := range outputs {
		outChans = append(outChans, c.ensureChan(name))
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, in, outChans)
	})
}

func (c *impl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return errors.New("conveyer already running")
	}
	c.running = true
	workers := append([]func(context.Context) error(nil), c.workers...)
	c.mu.Unlock()

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(workers))

	errCh := make(chan error, 1)

	for _, w := range workers {
		worker := w
		go func() {
			defer wg.Done()
			if err := worker(runCtx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	var runErr error
	select {
	case <-ctx.Done():
	case runErr = <-errCh:
	}

	cancel()
	wg.Wait()

	c.mu.Lock()
	chans := make([]chan string, 0, len(c.channels))
	for _, ch := range c.channels {
		chans = append(chans, ch)
	}
	c.running = false
	c.mu.Unlock()

	for _, ch := range chans {
		close(ch)
	}

	return runErr
}

func (c *impl) Send(input string, data string) error {
	c.mu.RLock()
	ch, ok := c.channels[input]
	c.mu.RUnlock()

	if !ok {
		return errChanNotFound
	}

	ch <- data
	return nil
}

func (c *impl) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.channels[output]
	c.mu.RUnlock()

	if !ok {
		return "", errChanNotFound
	}

	v, ok := <-ch
	if !ok {
		return undefinedData, nil
	}

	return v, nil
}
