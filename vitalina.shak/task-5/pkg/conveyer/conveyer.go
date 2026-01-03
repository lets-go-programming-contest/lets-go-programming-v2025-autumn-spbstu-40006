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
		mu:       sync.Mutex{},
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
	}
}

func (conv *Conveyer) provideChannel(channelName string) chan string {
	conv.mu.Lock()
	defer conv.mu.Unlock()

	channel, exists := conv.channels[channelName]
	if exists {
		return channel
	}

	newChannel  := make(chan string, conv.bufSize)
	conv.channels[channelName] = newChannel

	return newChannel
}

func (conv *Conveyer) get(channelName string) (chan string, bool) {
	conv.mu.Lock()
	defer conv.mu.Unlock()
	channel, exists := conv.channels[channelName]

	return channel, exists
}

func (conv *Conveyer) RegisterDecorator(
	handlerFn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	inputChannel := conv.provideChannel(input)
	outputChannel := conv.provideChannel(output)

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, inputChannel, outputChannel)
	})
}

func (conv *Conveyer) RegisterMultiplexer(
	handlerFn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	inputChannels := make([]chan string, 0, len(inputs))
	for _, inputID := range inputs {
		inputChannels = append(inputChannels, conv.provideChannel(inputID))
	}

	outputChannel := conv.provideChannel(output)

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, inputChannels, outputChannel)
	})
}

func (conv *Conveyer) RegisterSeparator(
	handlerFn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	inputChannel := conv.provideChannel(input)
	outputChannels := make([]chan string, 0, len(outputs))
	for _, outputID := range outputs {
		outputChannels = append(outputChannels, conv.provideChannel(outputID))
	}

	conv.handlers = append(conv.handlers, func(ctx context.Context) error {
		return handlerFn(ctx, inputChannel, outputChannels)
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
	channel, exists := conv.get(input)
	if !exists {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (conv *Conveyer) Recv(output string) (string, error) {
	channel, exists  := conv.get(output)
	if !exists {
		return "", ErrChanNotFound
	}

	data, isOpen  := <-channel

	if !isOpen  {
		return undefinedData, nil
	}

	return data, nil
}
