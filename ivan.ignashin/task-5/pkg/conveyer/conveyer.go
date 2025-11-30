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
		mu:       sync.Mutex{},
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
	handlerFunc func(ctx context.Context, input chan string, output chan string) error,
	inputName string,
	outputName string,
) {
	c.getOrCreateChan(inputName)
	c.getOrCreateChan(outputName)

	worker := func(ctx context.Context) error {
		input := c.getOrCreateChan(inputName)
		output := c.getOrCreateChan(outputName)

		return handlerFunc(ctx, input, output)
	}

	c.workers = append(c.workers, worker)
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
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

		return handlerFunc(ctx, inputs, output)
	}

	c.workers = append(c.workers, worker)
}

func (c *Conveyer) RegisterSeparator(
	handlerFunc func(ctx context.Context, input chan string, outputs []chan string) error,
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

		return handlerFunc(ctx, input, outputs)
	}

	c.workers = append(c.workers, worker)
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
	waitGroup := sync.WaitGroup{}
	errorChannel := make(chan error, len(c.workers))

	for _, currentWorker := range c.workers {
		waitGroup.Add(1)

		workerCopy := currentWorker
		go func() {
			defer waitGroup.Done()

			if err := workerCopy(ctx); err != nil {
				select {
				case errorChannel <- err:
				default:
				}
			}
		}()
	}

	go func() {
		waitGroup.Wait()
		close(errorChannel)
	}()

	select {
	case err := <-errorChannel:
		return err
	case <-ctx.Done():
		waitGroup.Wait()
		return nil
	}
}
