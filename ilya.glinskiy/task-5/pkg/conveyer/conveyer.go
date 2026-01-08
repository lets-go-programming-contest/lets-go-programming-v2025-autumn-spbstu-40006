package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrChanFull     = errors.New("channel is full")
)

type DecoratorFunc func(ctx context.Context, input chan string, output chan string) error

type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error

type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type Conveyer struct {
	size      int
	channels  map[string]chan string
	handlers  []func(ctx context.Context) error
	mu        sync.RWMutex
	wg        sync.WaitGroup
	runCtx    context.Context
	runCancel context.CancelFunc
	isRunning bool
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:      size,
		channels:  make(map[string]chan string),
		handlers:  make([]func(ctx context.Context) error, 0),
		isRunning: false,
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
		return nil, fmt.Errorf("%w: %s", ErrChanNotFound, name)
	}

	return ch, nil
}

func (c *Conveyer) RegisterDecorator(
	decoratorFn DecoratorFunc,
	input string,
	output string,
) {
	inputChan := c.getOrCreateChannel(input)
	outputChan := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return decoratorFn(ctx, inputChan, outputChan)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	multiplexerFn MultiplexerFunc,
	inputs []string,
	output string,
) {
	inputChans := make([]chan string, len(inputs))

	for i, input := range inputs {
		inputChans[i] = c.getOrCreateChannel(input)
	}

	outputChan := c.getOrCreateChannel(output)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return multiplexerFn(ctx, inputChans, outputChan)
	})
}

func (c *Conveyer) RegisterSeparator(
	separatorFn SeparatorFunc,
	input string,
	outputs []string,
) {
	inputChan := c.getOrCreateChannel(input)
	outputChans := make([]chan string, len(outputs))

	for i, output := range outputs {
		outputChans[i] = c.getOrCreateChannel(output)
	}

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return separatorFn(ctx, inputChan, outputChans)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	c.mu.Lock()

	if c.isRunning {
		c.mu.Unlock()
		return errors.New("conveyer is already running")
	}

	c.runCtx, c.runCancel = context.WithCancel(ctx)
	c.isRunning = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.isRunning = false
		c.mu.Unlock()
		c.runCancel()
	}()

	errChan := make(chan error, len(c.handlers))

	for _, handler := range c.handlers {
		c.wg.Add(1)

		handlerFn := handler
		go func() {
			defer c.wg.Done()

			if err := handlerFn(c.runCtx); err != nil {
				select {
				case errChan <- fmt.Errorf("handler error: %w", err):
				case <-c.runCtx.Done():
				}
			}
		}()
	}

	select {
	case err := <-errChan:
		c.runCancel()
		c.wg.Wait()
		c.closeAllChannels()

		return fmt.Errorf("conveyer stopped with error: %w", err)

	case <-c.runCtx.Done():
		c.wg.Wait()
		c.closeAllChannels()

		if errors.Is(c.runCtx.Err(), context.Canceled) {
			return nil
		}

		return fmt.Errorf("conveyer context error: %w", c.runCtx.Err())
	}
}

func (c *Conveyer) closeAllChannels() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, channel := range c.channels {
		for {
			select {
			case <-channel:
			default:
				close(channel)
				delete(c.channels, name)
				goto nextChannel
			}
		}

	nextChannel:
	}
}

func (c *Conveyer) Send(input string, data string) error {
	channel, err := c.getChannel(input)
	if err != nil {
		return fmt.Errorf("send error: %w", err)
	}

	select {
	case channel <- data:
		return nil
	case <-c.runCtx.Done():
		return fmt.Errorf("send cancelled: %w", c.runCtx.Err())
	default:
		return fmt.Errorf("send error: %w", ErrChanFull)
	}
}

func (c *Conveyer) Recv(output string) (string, error) {
	channel, err := c.getChannel(output)
	if err != nil {
		return "", fmt.Errorf("receive error: %w", err)
	}

	select {
	case data, ok := <-channel:
		if !ok {
			return "undefined", nil
		}

		return data, nil
	case <-c.runCtx.Done():
		return "", fmt.Errorf("receive cancelled: %w", c.runCtx.Err())
	}
}
