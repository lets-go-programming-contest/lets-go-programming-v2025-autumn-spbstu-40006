package conveyer

import (
	"context"
	"errors"
	"sync"
)

const (
	errMsgChanNotFound = "chan not found"
	valUndefined       = "undefined"
)

var errChanNotFound = errors.New(errMsgChanNotFound)

type Conveyer struct {
	mu       sync.RWMutex
	channels map[string]chan string
	tasks    []func(context.Context) error
	size     int
}

func New(size int) *Conveyer {
	return &Conveyer{
		mu:       sync.RWMutex{},
		channels: make(map[string]chan string),
		tasks:    make([]func(context.Context) error, 0),
		size:     size,
	}
}

func (c *Conveyer) getOrInitChannel(name string) chan string {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, ok := c.channels[name]; ok {
		return channel
	}

	channel := make(chan string, c.size)
	c.channels[name] = channel

	return channel
}

func (c *Conveyer) RegisterDecorator(
	handlerFunc func(context.Context, chan string, chan string) error,
	inputName, outputName string,
) {
	inputChannel := c.getOrInitChannel(inputName)
	outputChannel := c.getOrInitChannel(outputName)

	task := func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputChannel)
	}

	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterMultiplexer(
	handlerFunc func(context.Context, []chan string, chan string) error,
	inputNames []string,
	outputName string,
) {
	inputs := make([]chan string, 0, len(inputNames))
	for _, name := range inputNames {
		inputs = append(inputs, c.getOrInitChannel(name))
	}

	outputChannel := c.getOrInitChannel(outputName)

	task := func(ctx context.Context) error {
		return handlerFunc(ctx, inputs, outputChannel)
	}

	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterSeparator(
	handlerFunc func(context.Context, chan string, []chan string) error,
	inputName string,
	outputNames []string,
) {
	inputChannel := c.getOrInitChannel(inputName)

	outputs := make([]chan string, 0, len(outputNames))
	for _, name := range outputNames {
		outputs = append(outputs, c.getOrInitChannel(name))
	}

	task := func(ctx context.Context) error {
		return handlerFunc(ctx, inputChannel, outputs)
	}

	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var waitGroup sync.WaitGroup

	errorChannel := make(chan error, len(c.tasks))

	runTask := func(task func(context.Context) error) {
		defer waitGroup.Done()

		if err := task(ctx); err != nil {
			select {
			case errorChannel <- err:
				cancel()
			default:
			}
		}
	}

	for _, task := range c.tasks {
		waitGroup.Add(1)

		localTask := task
		go runTask(localTask)
	}

	done := make(chan struct{})

	go func() {
		waitGroup.Wait()
		close(done)
	}()

	select {
	case err := <-errorChannel:
		waitGroup.Wait()

		return err
	case <-done:
		return nil
	case <-ctx.Done():
		waitGroup.Wait()

		return nil
	}
}

func (c *Conveyer) Send(name string, data string) error {
	c.mu.RLock()
	channel, ok := c.channels[name]
	c.mu.RUnlock()

	if !ok {
		return errChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	c.mu.RLock()
	channel, ok := c.channels[name]
	c.mu.RUnlock()

	if !ok {
		return "", errChanNotFound
	}

	value, isOpen := <-channel
	if !isOpen {
		return valUndefined, nil
	}

	return value, nil
}
