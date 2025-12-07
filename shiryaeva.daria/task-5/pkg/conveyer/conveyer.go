package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

const UndefinedData = "undefined"

type Conveyer struct {
	size     int
	channels map[string]chan string
	workers  []func(ctx context.Context) error
	mu       sync.Mutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		workers:  make([]func(ctx context.Context) error, 0),
		mu:       sync.Mutex{},
	}
}

func (c *Conveyer) getOrCreateChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch

	return ch
}

func (c *Conveyer) getChan(name string) (chan string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch, exists := c.channels[name]
	if !exists {
		return nil, ErrChanNotFound
	}

	return ch, nil
}

func (c *Conveyer) RegisterDecorator(
	handler func(ctx context.Context, input, output chan string) error,
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

func (c *Conveyer) Send(inputName, data string) error {
	ch, err := c.getChan(inputName)
	if err != nil {
		return ErrChanNotFound
	}

	ch <- data

	return nil
}

func (c *Conveyer) Recv(outputName string) (string, error) {
	ch, err := c.getChan(outputName)
	if err != nil {
		return "", ErrChanNotFound
	}

	data, ok := <-ch
	if !ok {
		return UndefinedData, nil
	}

	return data, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	waitGroup := sync.WaitGroup{}
	errCh := make(chan error, len(c.workers))

	for _, worker := range c.workers {
		waitGroup.Add(1)

		go func(w func(ctx context.Context) error) {
			defer waitGroup.Done()

			if err := w(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(worker)
	}

	go func() {
		waitGroup.Wait()
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		waitGroup.Wait()

		return nil
	}
}
