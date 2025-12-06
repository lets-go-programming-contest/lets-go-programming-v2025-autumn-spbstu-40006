package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

const undefined = "undefined"

type ConveyerInterface interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)
	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)
	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)

	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type Pipeline struct {
	size     int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
	mu       sync.RWMutex
}

func New(size int) *Pipeline {
	return &Pipeline{
		size:     size,
		channels: make(map[string]chan string),
		handlers: []func(ctx context.Context) error{},
		mu:       sync.RWMutex{},
	}
}

func (pipe *Pipeline) getOrInitChannel(name string) chan string {
	pipe.mu.Lock()
	defer pipe.mu.Unlock()

	if channel, exists := pipe.channels[name]; exists {
		return channel
	}

	channel := make(chan string, pipe.size)
	pipe.channels[name] = channel

	return channel
}

func (pipe *Pipeline) RegisterDecorator(
	function func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	channelIn := pipe.getOrInitChannel(input)
	channelOut := pipe.getOrInitChannel(output)

	pipe.handlers = append(pipe.handlers, func(ctx context.Context) error {
		return function(ctx, channelIn, channelOut)
	})
}

func (pipe *Pipeline) RegisterMultiplexer(
	function func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	inStrings := make([]chan string, len(inputs))
	for i, name := range inputs {
		inStrings[i] = pipe.getOrInitChannel(name)
	}

	out := pipe.getOrInitChannel(output)

	pipe.handlers = append(pipe.handlers, func(ctx context.Context) error {
		return function(ctx, inStrings, out)
	})
}

func (pipe *Pipeline) RegisterSeparator(
	function func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	channelIn := pipe.getOrInitChannel(input)
	outs := make([]chan string, len(outputs))

	for i, name := range outputs {
		outs[i] = pipe.getOrInitChannel(name)
	}

	pipe.handlers = append(pipe.handlers, func(ctx context.Context) error {
		return function(ctx, channelIn, outs)
	})
}

func (pipe *Pipeline) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	for _, h := range pipe.handlers {
		handler := h

		group.Go(func() error {
			return handler(groupCtx)
		})
	}

	err := group.Wait()

	pipe.mu.Lock()
	for _, channel := range pipe.channels {
		close(channel)
	}
	pipe.mu.Unlock()

	if err != nil {
		return fmt.Errorf("conveyer run failed: %w", err)
	}

	return nil
}

func (pipe *Pipeline) Send(input string, data string) error {
	pipe.mu.RLock()
	channel, exists := pipe.channels[input]
	pipe.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}
	channel <- data

	return nil
}

func (pipe *Pipeline) Recv(output string) (string, error) {
	pipe.mu.RLock()
	channel, exists := pipe.channels[output]
	pipe.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	data, status := <-channel
	if !status {
		return undefined, nil
	}
	
	return data, nil
}
