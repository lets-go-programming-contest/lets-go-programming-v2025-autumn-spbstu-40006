package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

const undefinedData = "undefined"

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	bufSize  int
	mu       sync.Mutex
	channels map[string]chan string
	handlers []func(ctx context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		bufSize:  size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
	}
}

func (conv *Conveyer) provideChannel(id string) chan string {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	if channel, ok := conv.channels[id]; ok {
		return channel
	}
	channel := make(chan string, conv.bufSize)
	conv.channels[id] = channel
	return channel
}

func (conv *Conveyer) get(id string) (chan string, bool) {
	conv.mu.Lock()
	defer conv.mu.Unlock()
	ch, ok := conv.channels[id]
	return ch, ok
}

func (conv *Conveyer) RegisterDecorator(
	handlerFn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	in := conv.provideChannel(input)
	out := conv.provideChannel(output)

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, in, out)
	})
}

func (conv *Conveyer) RegisterMultiplexer(
	handlerFn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	ins := make([]chan string, 0, len(inputs))
	for _, id := range inputs {
		ins = append(ins, conv.provideChannel(id))
	}
	out := conv.provideChannel(output)

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, ins, out)
	})
}

func (conv *Conveyer) RegisterSeparator(
	handlerFn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	in := conv.provideChannel(input)
	outs := make([]chan string, 0, len(outputs))
	for _, id := range outputs {
		outs = append(outs, conv.provideChannel(id))
	}
	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, in, outs)
	})
}

func (conv *Conveyer) Run(ctx context.Context) error {
	group, groupCtx := errgroup.WithContext(ctx)

	for _, handlerFunc := range conv.handlers {
		handFunc := handlerFunc
		group.Go(func() error {
			return handFunc(groupCtx)
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("pipeline stopped: %w", err)
	}

	return nil
}

func (conv *Conveyer) Send(input string, data string) error {
	channel, ok := conv.get(input)
	if !ok {
		return ErrChanNotFound
	}

	channel <- data
	return nil
}

func (conv *Conveyer) Recv(output string) (string, error) {
	channel, ok := conv.get(output)
	if !ok {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return undefinedData, nil
	}

	return data, nil
}
