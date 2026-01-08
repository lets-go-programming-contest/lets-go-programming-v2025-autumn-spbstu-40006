package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrUndefined    = errors.New("undefined")
)

type DecoratorFunc func(ctx context.Context, input chan string, output chan string) error
type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type Conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mu       sync.RWMutex
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(ctx context.Context) error, 0),
	}
}

func (c *Conveyer) getOrCreateChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.channels[name]; exists {
		return ch
	}

	ch := make(chan string, c.size)
	c.channels[name] = ch
	return ch
}

func (c *Conveyer) getChannel(name string) (chan string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ch, exists := c.channels[name]
	if !exists {
		return nil, ErrChanNotFound
	}
	return ch, nil
}

func (c *Conveyer) RegisterDecorator(
	fn DecoratorFunc,
	input string,
	output string,
) {
	inputChan := c.getOrCreateChannel(input)
	outputChan := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChan)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	fn MultiplexerFunc,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, len(inputs))
	for i, input := range inputs {
		inputChans[i] = c.getOrCreateChannel(input)
	}
	outputChan := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChans, outputChan)
	})
}

func (c *Conveyer) RegisterSeparator(
	fn SeparatorFunc,
	input string,
	outputs []string,
) {
	inputChan := c.getOrCreateChannel(input)
	outputChans := make([]chan string, len(outputs))
	for i, output := range outputs {
		outputChans[i] = c.getOrCreateChannel(output)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return fn(ctx, inputChan, outputChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.ctx, c.cancel = context.WithCancel(ctx)
	defer c.cancel()

	errChan := make(chan error, len(c.handlers))

	for _, handler := range c.handlers {
		c.wg.Add(1)
		go func(h func(ctx context.Context) error) {
			defer c.wg.Done()
			if err := h(c.ctx); err != nil {
				select {
				case errChan <- err:
				case <-c.ctx.Done():
				}
			}
		}(handler)
	}

	select {
	case err := <-errChan:
		c.cancel()
		c.wg.Wait()
		c.closeAllChannels()
		return err
	case <-c.ctx.Done():
		c.wg.Wait()
		c.closeAllChannels()
		return c.ctx.Err()
	}
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, ch := range c.channels {
		select {
		case <-ch:
		default:
		}
		close(ch)
		delete(c.channels, name)
	}
}

func (c *Conveyer) Send(input string, data string) error {
	ch, err := c.getChannel(input)
	if err != nil {
		return fmt.Errorf("%w: %s", err, input)
	}

	select {
	case ch <- data:
		return nil
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
		return errors.New("channel is full")
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	ch, err := c.getChannel(output)
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, output)
	}

	select {
	case data, ok := <-ch:
		if !ok {
			return "undefined", nil
		}
		return data, nil
	case <-c.ctx.Done():
		return "", c.ctx.Err()
	}
}
