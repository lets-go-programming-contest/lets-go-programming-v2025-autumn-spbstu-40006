package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrChanFull     = errors.New("channel is full")
)

const undefinedData = "undefined"

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	handlers []func(context.Context) error
	mu       sync.RWMutex
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
		mu:       sync.RWMutex{},
	}
}

func (conveyer *conveyerImpl) getOrCreateChannel(name string) chan string {
	if existingChannel, exists := conveyer.channels[name]; exists {
		return existingChannel
	}

	newChannel := make(chan string, conveyer.size)
	conveyer.channels[name] = newChannel

	return newChannel
}

func (conveyer *conveyerImpl) getOrCreateChannels(names ...string) {
	for _, name := range names {
		conveyer.getOrCreateChannel(name)
	}
}

func (conveyer *conveyerImpl) RegisterDecorator(
	decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
	inputName string,
	outputName string,
) {
	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	conveyer.getOrCreateChannels(inputName, outputName)

	inputChannel := conveyer.channels[inputName]
	outputChannel := conveyer.channels[outputName]

	conveyer.handlers = append(conveyer.handlers, func(ctx context.Context) error {
		return decoratorFunc(ctx, inputChannel, outputChannel)
	})
}

func (conveyer *conveyerImpl) RegisterMultiplexer(
	multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	conveyer.getOrCreateChannels(inputNames...)
	conveyer.getOrCreateChannel(outputName)

	inputChannels := make([]chan string, len(inputNames))
	for index, name := range inputNames {
		inputChannels[index] = conveyer.channels[name]
	}

	outputChannel := conveyer.channels[outputName]

	conveyer.handlers = append(conveyer.handlers, func(ctx context.Context) error {
		return multiplexerFunc(ctx, inputChannels, outputChannel)
	})
}

func (conveyer *conveyerImpl) RegisterSeparator(
	separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	conveyer.getOrCreateChannel(inputName)
	conveyer.getOrCreateChannels(outputNames...)

	inputChannel := conveyer.channels[inputName]
	outputChannels := make([]chan string, len(outputNames))

	for index, name := range outputNames {
		outputChannels[index] = conveyer.channels[name]
	}

	conveyer.handlers = append(conveyer.handlers, func(ctx context.Context) error {
		return separatorFunc(ctx, inputChannel, outputChannels)
	})
}

func (conveyer *conveyerImpl) Run(ctx context.Context) error {
	errorGroup, groupContext := errgroup.WithContext(ctx)

	for _, handlerFunc := range conveyer.handlers {
		currentHandler := handlerFunc

		errorGroup.Go(func() error {
			return currentHandler(groupContext)
		})
	}

	if err := errorGroup.Wait(); err != nil {
		return fmt.Errorf("conveyer run: %w", err)
	}

	return nil
}

func (conveyer *conveyerImpl) Send(inputName string, data string) error {
	conveyer.mu.RLock()
	defer conveyer.mu.RUnlock()

	channel, exists := conveyer.channels[inputName]
	if !exists {
		return ErrChanNotFound
	}

	select {
	case channel <- data:
		return nil
	default:
		return ErrChanFull
	}
}

func (conveyer *conveyerImpl) Recv(outputName string) (string, error) {
	conveyer.mu.RLock()
	channel, exists := conveyer.channels[outputName]
	conveyer.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	data, isOpen := <-channel
	if !isOpen {
		return undefinedData, nil
	}

	return data, nil
}
