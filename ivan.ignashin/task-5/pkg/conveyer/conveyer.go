package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type conveyer interface {
	RegisterDecorator(fn func(context.Context, chan string, chan string) error, input string, output string)
	RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error, inputs []string, output string)
	RegisterSeparator(fn func(context.Context, chan string, []chan string) error, input string, outputs []string)
	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type Conveyer struct {
	size    int
	chans   map[string]chan string
	runners []func(context.Context) error
	mu      sync.Mutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:  size,
		chans: make(map[string]chan string),
	}
}

func (c *Conveyer) getOrCreate(id string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch, ok := c.chans[id]
	if !ok {
		ch = make(chan string, c.size)
		c.chans[id] = ch
	}
	return ch
}

func (c *Conveyer) get(id string) (chan string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch, ok := c.chans[id]
	if !ok {
		return nil, ErrChanNotFound
	}
	return ch, nil
}

func (c *Conveyer) RegisterDecorator(fn func(context.Context, chan string, chan string) error,
	input string, output string) {

	in := c.getOrCreate(input)
	out := c.getOrCreate(output)

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyer) RegisterMultiplexer(fn func(context.Context, []chan string, chan string) error,
	inputs []string, output string) {

	var inCh []chan string
	for _, id := range inputs {
		inCh = append(inCh, c.getOrCreate(id))
	}
	out := c.getOrCreate(output)

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, inCh, out)
	})
}

func (c *Conveyer) RegisterSeparator(fn func(context.Context, chan string, []chan string) error,
	input string, outputs []string) {

	in := c.getOrCreate(input)
	var out []chan string
	for _, id := range outputs {
		out = append(out, c.getOrCreate(id))
	}

	c.runners = append(c.runners, func(ctx context.Context) error {
		return fn(ctx, in, out)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(c.runners))

	for _, r := range c.runners {
		wg.Add(1)
		go func(run func(context.Context) error) {
			defer wg.Done()
			if err := run(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(r)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err, ok := <-errCh:
		if ok {
			return err
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Conveyer) Send(id string, data string) error {
	ch, err := c.get(id)
	if err != nil {
		return err
	}
	ch <- data
	return nil
}

func (c *Conveyer) Recv(id string) (string, error) {
	ch, err := c.get(id)
	if err != nil {
		return "", err
	}
	v, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return v, nil
}

var _ conveyer = (*Conveyer)(nil)
