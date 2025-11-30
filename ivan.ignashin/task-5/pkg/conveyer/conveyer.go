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

	if channel, exists := c.channels[name]; exists {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel
	return channel
}

func (c *Conveyer) getChan(name string) (chan string, error) {
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
	c.getOrCreateChan(inputName)
	c.getOrCreateChan(outputName)

	handler := func(ctx context.Context) error {
		input := c.getOrCreateChan(inputName)
		output := c.getOrCreateChan(outputName)
		return fn(ctx, input, output)
	}

	c.workers = append(c.workers, handler)
}

func (c *Conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	for _, name := range inputNames {
		c.getOrCreateChan(name)
	}
	c.getOrCreateChan(outputName)

	handler := func(ctx context.Context) error {
		inputs := make([]chan string, len(inputNames))
		for i, name := range inputNames {
			inputs[i] = c.getOrCreateChan(name)
		}

		output := c.getOrCreateChan(outputName)
		return fn(ctx, inputs, output)
	}

	c.workers = append(c.workers, handler)
}

func (c *Conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	c.getOrCreateChan(inputName)
	for _, name := range outputNames {
		c.getOrCreateChan(name)
	}

	handler := func(ctx context.Context) error {
		input := c.getOrCreateChan(inputName)
		outputs := make([]chan string, len(outputNames))
		for i, name := range outputNames {
			outputs[i] = c.getOrCreateChan(name)
		}

		return fn(ctx, input, outputs)
	}

	c.workers = append(c.workers, handler)
}

func (c *Conveyer) Send(inputName string, data string) error {
	channel, err := c.getChan(inputName)
	if err != nil {
		return ErrChanNotFound
	}

	channel <- data
	return nil
}

func (c *Conveyer) Recv(outputName string) (string, error) {
	channel, err := c.getChan(outputName)
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
	errCh := make(chan error, len(c.workers))

	for _, worker := range c.workers {
		wg.Add(1)

		currentWorker := worker
		go func() {
			defer wg.Done()

			if err := currentWorker(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		wg.Wait()
		return nil
	}
}
