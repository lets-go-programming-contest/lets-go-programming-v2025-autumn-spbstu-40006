package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	mu      sync.Mutex
	chans   map[string]chan string
	size    int
	workers []func(context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		chans:   make(map[string]chan string),
		size:    size,
		workers: []func(context.Context) error{},
	}
}

func (c *Conveyer) getChan(name string) (chan string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, exists := c.chans[name]
	if !exists {
		return nil, ErrChanNotFound
	}
	return ch, nil
}

func (c *Conveyer) createChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.chans[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.chans[name] = ch
	return ch
}

func (c *Conveyer) RegisterDecorator(
	fn func(context.Context, chan string, chan string) error,
	inName string,
	outName string,
) {
	in := c.createChan(inName)
	out := c.createChan(outName)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(context.Context, []chan string, chan string) error,
	inNames []string,
	outName string,
) {
	ins := make([]chan string, len(inNames))
	for i, name := range inNames {
		ins[i] = c.createChan(name)
	}
	out := c.createChan(outName)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, ins, out)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn func(context.Context, chan string, []chan string) error,
	inName string,
	outNames []string,
) {
	in := c.createChan(inName)
	outs := make([]chan string, len(outNames))
	for i, name := range outNames {
		outs[i] = c.createChan(name)
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return fn(ctx, in, outs)
	})
}

func (c *Conveyer) Send(chanName string, value string) error {
	ch, err := c.getChan(chanName)
	if err != nil {
		return err
	}
	ch <- value
	return nil
}

func (c *Conveyer) Recv(chanName string) (string, error) {
	ch, err := c.getChan(chanName)
	if err != nil {
		return "", err
	}

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return val, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(c.workers))

	for _, worker := range c.workers {
		wg.Add(1)
		go func(w func(context.Context) error) {
			defer wg.Done()
			if err := w(ctx); err != nil {
				select {
				case errChan <- err:
				default:
				}
			}
		}(worker)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case err, ok := <-errChan:
		if ok {
			return err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
