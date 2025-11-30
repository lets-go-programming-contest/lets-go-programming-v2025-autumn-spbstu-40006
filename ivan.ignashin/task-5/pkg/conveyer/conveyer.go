package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound = errors.New("chan not found")
)

type Conveyer struct {
	mu       sync.Mutex
	channels map[string]chan string
	size     int
	handlers []func(ctx context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		size:     size,
		handlers: make([]func(ctx context.Context) error, 0),
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exists := c.channels[name]; exists {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel
	return channel
}

func (c *Conveyer) getChannel(name string) (chan string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel, exists := c.channels[name]
	if !exists {
		return nil, ErrChanNotFound
	}
	return channel, nil
}

func (c *Conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	inputName string,
	outputName string,
) {
	input := c.getOrCreateChannel(inputName)
	output := c.getOrCreateChannel(outputName)

	handler := func(ctx context.Context) error {
		return fn(ctx, input, output)
	}

	c.handlers = append(c.handlers, handler)
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	inputs := make([]chan string, len(inputNames))
	for i, name := range inputNames {
		inputs[i] = c.getOrCreateChannel(name)
	}
	output := c.getOrCreateChannel(outputName)

	handler := func(ctx context.Context) error {
		return fn(ctx, inputs, output)
	}

	c.handlers = append(c.handlers, handler)
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	input := c.getOrCreateChannel(inputName)
	outputs := make([]chan string, len(outputNames))
	for i, name := range outputNames {
		outputs[i] = c.getOrCreateChannel(name)
	}

	handler := func(ctx context.Context) error {
		return fn(ctx, input, outputs)
	}

	c.handlers = append(c.handlers, handler)
}

func (c *Conveyer) Send(inputName string, data string) error {
	channel, err := c.getChannel(inputName)
	if err != nil {
		return ErrChanNotFound
	}

	channel <- data
	return nil
}

func (c *Conveyer) Recv(outputName string) (string, error) {
	channel, err := c.getChannel(outputName)
	if err != nil {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return "undefined", nil
	}
	return data, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}
	errCh := make(chan error, len(c.handlers))

	for _, handler := range c.handlers {
		wg.Add(1)
		go func(h func(ctx context.Context) error) {
			defer wg.Done()
			if err := h(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(handler)
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
