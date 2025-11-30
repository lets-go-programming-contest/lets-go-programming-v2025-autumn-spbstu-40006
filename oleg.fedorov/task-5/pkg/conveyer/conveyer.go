package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrChanNotFound = errors.New("chan not found")
	ErrChanFull     = errors.New("channel is full")
)

type conveyerImpl struct {
	size     int
	channels map[string]chan string
	handlers []handler
	mu       sync.RWMutex
}

type handler struct {
	run func(ctx context.Context) error
}

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]handler, 0),
		mu:       sync.RWMutex{},
	}
}

func (conveyer *conveyerImpl) RegisterDecorator(
	decoratorFunc func(ctx context.Context, input chan string, output chan string) error,
	inputName string,
	outputName string,
) {
	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	inputChannel := conveyer.getOrCreateChannel(inputName)
	outputChannel := conveyer.getOrCreateChannel(outputName)

	conveyer.handlers = append(conveyer.handlers, handler{
		run: func(ctx context.Context) error {
			return decoratorFunc(ctx, inputChannel, outputChannel)
		},
	})
}

func (conveyer *conveyerImpl) RegisterMultiplexer(
	multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	inputChannels := make([]chan string, len(inputNames))
	for index, name := range inputNames {
		inputChannels[index] = conveyer.getOrCreateChannel(name)
	}

	outputChannel := conveyer.getOrCreateChannel(outputName)

	conveyer.handlers = append(conveyer.handlers, handler{
		run: func(ctx context.Context) error {
			return multiplexerFunc(ctx, inputChannels, outputChannel)
		},
	})
}

func (conveyer *conveyerImpl) RegisterSeparator(
	separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	inputChannel := conveyer.getOrCreateChannel(inputName)
	outputChannels := make([]chan string, len(outputNames))

	for index, name := range outputNames {
		outputChannels[index] = conveyer.getOrCreateChannel(name)
	}

	conveyer.handlers = append(conveyer.handlers, handler{
		run: func(ctx context.Context) error {
			return separatorFunc(ctx, inputChannel, outputChannels)
		},
	})
}

func (conveyer *conveyerImpl) Run(ctx context.Context) error {
	var waitGroup sync.WaitGroup

	errorChannel := make(chan error, len(conveyer.handlers))

	for _, handlerItem := range conveyer.handlers {
		waitGroup.Add(1)

		currentHandler := handlerItem
		go func() {
			defer waitGroup.Done()

			if err := currentHandler.run(ctx); err != nil {
				errorChannel <- err
			}
		}()
	}

	var runError error

	select {
	case <-ctx.Done():
		runError = ctx.Err()
	case runError = <-errorChannel:
	}

	conveyer.closeAllChannels()

	waitGroup.Wait()

	return runError
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
		return "undefined", nil
	}

	return data, nil
}

func (conveyer *conveyerImpl) getOrCreateChannel(name string) chan string {
	if existingChannel, exists := conveyer.channels[name]; exists {
		return existingChannel
	}

	newChannel := make(chan string, conveyer.size)
	conveyer.channels[name] = newChannel

	return newChannel
}

func (conveyer *conveyerImpl) closeAllChannels() {
	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	for name, channel := range conveyer.channels {
		close(channel)
		delete(conveyer.channels, name)
	}
}
