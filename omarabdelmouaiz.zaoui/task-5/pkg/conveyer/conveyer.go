package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrChanNotFound = errors.New("chan not found")

type handler func(ctx context.Context) error

const undefined = "undefined"

type Pipeline struct {
	size       int
	channels   map[string]chan string
	handlers   []handler
	mutexChans sync.RWMutex
}

func New(size int) *Pipeline {
	return &Pipeline{
		size:       size,
		channels:   make(map[string]chan string),
		handlers:   []handler{},
		mutexChans: sync.RWMutex{},
	}
}

func (pipe *Pipeline) register(ch string) chan string {
	if _, exists := pipe.channels[ch]; !exists {
		pipe.channels[ch] = make(chan string, pipe.size)
	}

	return pipe.channels[ch]
}

func (pipe *Pipeline) RegisterDecorator(
	function func(
		ctx context.Context,
		input chan string,
		output chan string,
	) error,
	input string,
	output string,
) {
	pipe.mutexChans.Lock()

	in := pipe.register(input)
	out := pipe.register(output)
	pipe.handlers = append(pipe.handlers, handler(func(ctx context.Context) error {
		return function(ctx, in, out)
	}))

	pipe.mutexChans.Unlock()
}

func (pipe *Pipeline) RegisterMultiplexer(
	function func(
		ctx context.Context,
		inputs []chan string,
		output chan string,
	) error,
	inputs []string,
	output string,
) {
	pipe.mutexChans.Lock()

	ins := make([]chan string, len(inputs))
	for i, ch := range inputs {
		ins[i] = pipe.register(ch)
	}

	out := pipe.register(output)
	pipe.handlers = append(pipe.handlers, handler(func(ctx context.Context) error {
		return function(ctx, ins, out)
	}))

	pipe.mutexChans.Unlock()
}

func (pipe *Pipeline) RegisterSeparator(
	function func(
		ctx context.Context,
		input chan string,
		outputs []chan string,
	) error,
	input string,
	outputs []string,
) {
	pipe.mutexChans.Lock()

	inChan := pipe.register(input)
	outs := make([]chan string, len(outputs))

	for i, ch := range outputs {
		outs[i] = pipe.register(ch)
	}

	pipe.handlers = append(pipe.handlers, handler(func(ctx context.Context) error {
		return function(ctx, inChan, outs)
	}))

	pipe.mutexChans.Unlock()
}

func (pipe *Pipeline) Run(ctx context.Context) error {
	pipe.mutexChans.RLock()
	errgr, ctx := errgroup.WithContext(ctx)

	for _, handler := range pipe.handlers {
		errgr.Go(func() error {
			return handler(ctx)
		})
	}

	err := errgr.Wait()

	for _, ch := range pipe.channels {
		close(ch)
	}

	pipe.mutexChans.RUnlock()

	if err != nil {
		return fmt.Errorf("pipeline failed: %w", err)
	}

	return nil
}

func (pipe *Pipeline) Send(input string, data string) error {
	pipe.mutexChans.RLock()

	channel, exists := pipe.channels[input]

	pipe.mutexChans.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (pipe *Pipeline) Recv(output string) (string, error) {
	pipe.mutexChans.RLock()

	channel, exists := pipe.channels[output]

	pipe.mutexChans.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return undefined, nil
	}

	return data, nil
}
