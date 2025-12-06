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
		workers:  []func(ctx context.Context) error{},
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

	val, ok := <-ch
	if !ok {
		return "undefined", nil
	}
	return val, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(c.workers))

	for _, w := range c.workers {
		wg.Add(1)
		worker := w
		go func() {
			defer wg.Done()
			if err := worker(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		c.mu.Lock()
		for _, ch := range c.channels {
			close(ch)
		}
		c.mu.Unlock()
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
