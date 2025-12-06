package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	channels map[string]chan string
	size     int
	workers  []func(ctx context.Context) error
	mu       sync.Mutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		size:     size,
		workers:  make([]func(ctx context.Context) error, 0),
	}
}

func (c *Conveyer) getOrCreateChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, ok := c.channels[name]; ok {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) getChan(name string) (chan string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, ok := c.channels[name]
	if !ok {
		return nil, ErrChanNotFound
	}
	return ch, nil
}

func (c *Conveyer) RegisterDecorator(
	handler func(ctx context.Context, input chan string, output chan string) error,
	inputName, outputName string,
) {
	c.getOrCreateChan(inputName)
	c.getOrCreateChan(outputName)

	worker := func(ctx context.Context) error {
		input := c.getOrCreateChan(inputName)
		output := c.getOrCreateChan(outputName)
		return handler(ctx, input, output)
	}

	c.workers = append(c.workers, worker)
}

func (c *Conveyer) RegisterSeparator(
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	c.getOrCreateChan(inputName)
	for _, name := range outputNames {
		c.getOrCreateChan(name)
	}

	worker := func(ctx context.Context) error {
		input := c.getOrCreateChan(inputName)
		outputs := make([]chan string, len(outputNames))
		for i, name := range outputNames {
			outputs[i] = c.getOrCreateChan(name)
		}
		return handler(ctx, input, outputs)
	}

	c.workers = append(c.workers, worker)
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	for _, name := range inputNames {
		c.getOrCreateChan(name)
	}
	c.getOrCreateChan(outputName)

	worker := func(ctx context.Context) error {
		inputs := make([]chan string, len(inputNames))
		for i, name := range inputNames {
			inputs[i] = c.getOrCreateChan(name)
		}
		output := c.getOrCreateChan(outputName)
		return handler(ctx, inputs, output)
	}

	c.workers = append(c.workers, worker)
}

func (c *Conveyer) Send(ctx context.Context, name, data string) error {
	ch, err := c.getChan(name)
	if err != nil {
		return ErrChanNotFound
	}
	select {
	case ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Conveyer) Recv(ctx context.Context, name string) (string, error) {
	ch, err := c.getChan(name)
	if err != nil {
		return "", ErrChanNotFound
	}
	select {
	case val, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return val, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(c.workers))

	for _, w := range c.workers {
		wg.Add(1)
		worker := w
		go func() {
			defer wg.Done()
			done := make(chan struct{})
			go func() {
				if err := worker(ctx); err != nil {
					select {
					case errCh <- err:
					default:
					}
				}
				close(done)
			}()
			select {
			case <-ctx.Done():
			case <-done:
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}

	return nil
}
