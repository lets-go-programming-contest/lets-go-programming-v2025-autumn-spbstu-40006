package conveyer

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Conveyer struct {
	size     int
	channels map[string]chan string
	handlers []func(context.Context) error
	mu       sync.RWMutex
}

func New(size int) *Conveyer {
	return &Conveyer{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
		mu:       sync.RWMutex{},
	}
}

func (conveyer *Conveyer) RegisterDecorator(
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

func (conveyer *Conveyer) RegisterMultiplexer(
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

func (conveyer *Conveyer) RegisterSeparator(
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

func (conveyer *Conveyer) Run(ctx context.Context) error {
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

func (conveyer *Conveyer) Send(inputName string, data string) error {
	conveyer.mu.RLock()
	defer conveyer.mu.RUnlock()

	channel, exists := conveyer.channels[inputName]
	if !exists {
		return fmt.Errorf("chan not found")
	}

	channel <- data

	return nil
}

func (conveyer *Conveyer) Recv(outputName string) (string, error) {
	conveyer.mu.RLock()
	channel, exists := conveyer.channels[outputName]
	conveyer.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("chan not found")
	}

	data, isOpen := <-channel
	if !isOpen {
		return "undefinedData", nil
	}

	return data, nil
}

func (conveyer *Conveyer) getOrCreateChannel(name string) chan string {
	if existingChannel, exists := conveyer.channels[name]; exists {
		return existingChannel
	}

	newChannel := make(chan string, conveyer.size)
	conveyer.channels[name] = newChannel

	return newChannel
}

func (conveyer *Conveyer) getOrCreateChannels(names ...string) {
	for _, name := range names {
		conveyer.getOrCreateChannel(name)
	}
}
