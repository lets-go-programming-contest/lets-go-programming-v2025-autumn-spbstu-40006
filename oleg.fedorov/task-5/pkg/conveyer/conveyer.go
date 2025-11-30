package conveyer

import (
	"context"
	"errors"
	"sync"
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

type Conveyer interface {
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

func New(size int) *conveyerImpl {
	return &conveyerImpl{
		size:     size,
		channels: make(map[string]chan string),
		handlers: make([]handler, 0),
	}
}

func (conveyer *conveyerImpl) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	inputName, outputName string,
) {
	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	inputChan := conveyer.getOrCreateChannel(inputName)
	outputChan := conveyer.getOrCreateChannel(outputName)

	conveyer.handlers = append(conveyer.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, inputChan, outputChan)
		},
	})
}

func (conveyer *conveyerImpl) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string, outputName string,
) {
	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	inputs := make([]chan string, len(inputNames))
	for index, name := range inputNames {
		inputs[index] = conveyer.getOrCreateChannel(name)
	}
	outputChan := conveyer.getOrCreateChannel(outputName)

	conveyer.handlers = append(conveyer.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, inputs, outputChan)
		},
	})
}

func (conveyer *conveyerImpl) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string, outputNames []string,
) {
	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	inputChan := conveyer.getOrCreateChannel(inputName)
	outputs := make([]chan string, len(outputNames))
	for index, name := range outputNames {
		outputs[index] = conveyer.getOrCreateChannel(name)
	}

	conveyer.handlers = append(conveyer.handlers, handler{
		run: func(ctx context.Context) error {
			return fn(ctx, inputChan, outputs)
		},
	})
}

func (conveyer *conveyerImpl) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(conveyer.handlers))

	for _, hand := range conveyer.handlers {
		wg.Add(1)
		go func(handler handler) {
			defer wg.Done()
			if err := handler.run(ctx); err != nil {
				errCh <- err
			}
		}(hand)
	}

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case err = <-errCh:
	}

	conveyer.closeAllChannels()

	wg.Wait()

	return err
}

func (conveyer *conveyerImpl) Send(inputName string, data string) error {
	conveyer.mu.RLock()
	defer conveyer.mu.RUnlock()

	ch, exists := conveyer.channels[inputName]
	if !exists {
		return errors.New("chan not found")
	}

	select {
	case ch <- data:
		return nil
	default:
		return errors.New("channel is full")
	}
}

func (conveyer *conveyerImpl) Recv(outputName string) (string, error) {
	conveyer.mu.RLock()
	ch, exists := conveyer.channels[outputName]
	conveyer.mu.RUnlock()

	if !exists {
		return "", errors.New("chan not found")
	}

	data, ok := <-ch
	if !ok {
		return "undefined", nil
	}

	return data, nil
}

func (conveyer *conveyerImpl) getOrCreateChannel(name string) chan string {
	if ch, exists := conveyer.channels[name]; exists {
		return ch
	}
	ch := make(chan string, conveyer.size)
	conveyer.channels[name] = ch
	return ch
}

func (conveyer *conveyerImpl) closeAllChannels() {
	conveyer.mu.Lock()
	defer conveyer.mu.Unlock()

	for name, ch := range conveyer.channels {
		close(ch)
		delete(conveyer.channels, name)
	}
}
