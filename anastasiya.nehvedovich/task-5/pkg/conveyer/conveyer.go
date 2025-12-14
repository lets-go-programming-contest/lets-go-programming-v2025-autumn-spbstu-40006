package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	mu        sync.RWMutex
	streams   map[string]chan string
	processes []func(ctx context.Context) error
	bufSize   int
}

func New(size int) *Conveyer {
	return &Conveyer{
		mu:        sync.RWMutex{},
		streams:   make(map[string]chan string),
		processes: make([]func(ctx context.Context) error, 0),
		bufSize:   size,
	}
}

func (c *Conveyer) ensureChan(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exists := c.streams[name]; exists {
		return channel
	}

	channel := make(chan string, c.bufSize)
	c.streams[name] = channel

	return channel
}

func (c *Conveyer) RegisterDecorator(
	callback func(context.Context, chan string, chan string) error,
	inputName string,
	outputName string,
) {
	c.ensureChan(inputName)
	c.ensureChan(outputName)

	c.processes = append(c.processes, func(ctx context.Context) error {
		input := c.ensureChan(inputName)
		output := c.ensureChan(outputName)

		return callback(ctx, input, output)
	})
}

func (c *Conveyer) RegisterMultiplexer(
	callback func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	for _, name := range inputNames {
		c.ensureChan(name)
	}

	c.ensureChan(outputName)

	c.processes = append(c.processes, func(ctx context.Context) error {
		inputs := make([]chan string, len(inputNames))

		for index, name := range inputNames {
			inputs[index] = c.ensureChan(name)
		}

		return callback(ctx, inputs, c.ensureChan(outputName))
	})
}

func (c *Conveyer) RegisterSeparator(
	callback func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	c.ensureChan(inputName)

	for _, name := range outputNames {
		c.ensureChan(name)
	}

	c.processes = append(c.processes, func(ctx context.Context) error {
		outputs := make([]chan string, len(outputNames))

		for index, name := range outputNames {
			outputs[index] = c.ensureChan(name)
		}

		return callback(ctx, c.ensureChan(inputName), outputs)
	})
}

func (c *Conveyer) Send(pipeName string, data string) error {
	c.mu.RLock()
	channel, exists := c.streams[pipeName]
	c.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(pipeName string) (string, error) {
	c.mu.RLock()
	channel, exists := c.streams[pipeName]
	c.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	value, isOpen := <-channel

	if !isOpen {
		return "undefined", nil
	}

	return value, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	errorGroup, groupCtx := errgroup.WithContext(ctx)

	for _, processor := range c.processes {
		processorCopy := processor

		errorGroup.Go(func() error {
			return processorCopy(groupCtx)
		})
	}

	if err := errorGroup.Wait(); err != nil {
		return fmt.Errorf("conveyer run error: %w", err)
	}

	return nil
}
