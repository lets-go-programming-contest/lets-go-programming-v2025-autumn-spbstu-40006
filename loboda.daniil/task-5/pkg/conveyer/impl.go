package conveyer

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrAlreadyRunning = errors.New("conveyer already running")
	ErrChanNotFound   = errors.New("chan not found")
)

const UndefinedData = "undefined"

type conveyerImpl struct {
	size int

	mu       sync.RWMutex
	channels map[string]chan string
	workers  []func(ctx context.Context) error

	running bool
}

func newConveyerImpl(size int) *conveyerImpl {
	if size < 0 {
		size = 0
	}

	return &conveyerImpl{
		size:     size,
		mu:       sync.RWMutex{},
		channels: make(map[string]chan string),
		workers:  make([]func(ctx context.Context) error, 0),
		running:  false,
	}
}

func (c *conveyerImpl) ensureChannel(channelName string) chan string {
	if existingChannel, exists := c.channels[channelName]; exists {
		return existingChannel
	}

	createdChannel := make(chan string, c.size)
	c.channels[channelName] = createdChannel

	return createdChannel
}

func (c *conveyerImpl) RegisterDecorator(handler DecoratorFunc, inputName string, outputName string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChannel := c.ensureChannel(inputName)
	outputChannel := c.ensureChannel(outputName)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return handler(ctx, inputChannel, outputChannel)
	})
}

func (c *conveyerImpl) RegisterMultiplexer(handler MultiplexerFunc, inputNames []string, outputName string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChans := make([]chan string, 0, len(inputNames))
	for _, inputName := range inputNames {
		inputChans = append(inputChans, c.ensureChannel(inputName))
	}

	outputChan := c.ensureChannel(outputName)

	c.workers = append(c.workers, func(ctx context.Context) error {
		return handler(ctx, inputChans, outputChan)
	})
}

func (c *conveyerImpl) RegisterSeparator(handler SeparatorFunc, inputName string, outputNames []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	inputChan := c.ensureChannel(inputName)

	outputChans := make([]chan string, 0, len(outputNames))
	for _, outputName := range outputNames {
		outputChans = append(outputChans, c.ensureChannel(outputName))
	}

	c.workers = append(c.workers, func(ctx context.Context) error {
		return handler(ctx, inputChan, outputChans)
	})
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()

		return ErrAlreadyRunning
	}

	c.running = true
	workersCopy := append([]func(context.Context) error(nil), c.workers...)
	c.mu.Unlock()

	runCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(workersCopy))

	errorChan := make(chan error, 1)

	for _, worker := range workersCopy {
		currentWorker := worker

		workerFunc := func() {
			defer waitGroup.Done()

			err := currentWorker(runCtx)
			if err != nil {
				select {
				case errorChan <- err:
				default:
				}
			}
		}

		go workerFunc()
	}

	var runErr error
	select {
	case <-ctx.Done():
	case runErr = <-errorChan:
	}

	cancelFunc()
	waitGroup.Wait()

	c.mu.Lock()
	channelsCopy := make([]chan string, 0, len(c.channels))
	for _, channel := range c.channels {
		channelsCopy = append(channelsCopy, channel)
	}

	c.running = false
	c.mu.Unlock()

	for _, channel := range channelsCopy {
		safeClose(channel)
	}

	return runErr
}

func (c *conveyerImpl) Send(inputName string, data string) error {
	c.mu.RLock()
	inputChan, exists := c.channels[inputName]
	c.mu.RUnlock()

	if !exists {
		return ErrChanNotFound
	}

	inputChan <- data

	return nil
}

func (c *conveyerImpl) Recv(outputName string) (string, error) {
	c.mu.RLock()
	outputChan, exists := c.channels[outputName]
	c.mu.RUnlock()

	if !exists {
		return "", ErrChanNotFound
	}

	data, channelOpen := <-outputChan
	if !channelOpen {
		return UndefinedData, nil
	}

	return data, nil
}

func safeClose(channel chan string) {
	defer func() { _ = recover() }()
	close(channel)
}
