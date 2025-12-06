package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChannelNotFound = errors.New("chan not found")
var ErrCannotBeDecorated = errors.New("can't be decorated")

type Conveyer struct {
	channelMap map[string]chan string
	taskSlice  []func(context.Context) error
	mutex      sync.Mutex
	bufferSize int
}

func NewConveyer(size int) *Conveyer {
	return &Conveyer{
		channelMap: make(map[string]chan string),
		taskSlice:  make([]func(context.Context) error, 0),
		mutex:      sync.Mutex{},
		bufferSize: size,
	}
}

func New(size int) *Conveyer {
	return NewConveyer(size)
}

func (c *Conveyer) obtainChannel(identifier string) chan string {
	c.mutex.Lock()
	channel, exists := c.channelMap[identifier]
	if !exists {
		channel = make(chan string, c.bufferSize)
		c.channelMap[identifier] = channel
	}
	c.mutex.Unlock()

	return channel
}

func (c *Conveyer) fetchChannel(identifier string) (chan string, error) {
	c.mutex.Lock()
	channel, exists := c.channelMap[identifier]
	c.mutex.Unlock()

	if !exists {
		return nil, ErrChannelNotFound
	}

	return channel, nil
}

func (c *Conveyer) RegisterDecorator(
	processor func(context.Context, chan string, chan string) error,
	sourceID string,
	targetID string,
) {
	task := func(ctx context.Context) error {
		sourceChannel := c.obtainChannel(sourceID)
		targetChannel := c.obtainChannel(targetID)

		return processor(ctx, sourceChannel, targetChannel)
	}

	c.taskSlice = append(c.taskSlice, task)
}

func (c *Conveyer) RegisterSeparator(
	processor func(context.Context, chan string, []chan string) error,
	inputID string,
	outputIDs []string,
) {
	task := func(ctx context.Context) error {
		inputChannel := c.obtainChannel(inputID)

		outputChannels := make([]chan string, 0, len(outputIDs))
		for _, id := range outputIDs {
			outputChannels = append(outputChannels, c.obtainChannel(id))
		}

		return processor(ctx, inputChannel, outputChannels)
	}

	c.taskSlice = append(c.taskSlice, task)
}

func (c *Conveyer) RegisterMultiplexer(
	processor func(context.Context, []chan string, chan string) error,
	inputIDs []string,
	outputID string,
) {
	task := func(ctx context.Context) error {
		inputChannels := make([]chan string, 0, len(inputIDs))
		for _, id := range inputIDs {
			inputChannels = append(inputChannels, c.obtainChannel(id))
		}

		outputChannel := c.obtainChannel(outputID)

		return processor(ctx, inputChannels, outputChannel)
	}

	c.taskSlice = append(c.taskSlice, task)
}

func (c *Conveyer) Send(ctx context.Context, identifier string, value string) error {
	channel := c.obtainChannel(identifier)

	select {
	case channel <- value:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Conveyer) Recv(ctx context.Context, identifier string) (string, error) {
	channel, err := c.fetchChannel(identifier)
	if err != nil {
		return "", err
	}

	select {
	case val, open := <-channel:
		if !open {
			return "", nil
		}
		return val, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (c *Conveyer) Run(ctx context.Context) error {
	var waitGroup sync.WaitGroup
	errorChannel := make(chan error, len(c.taskSlice))

	for _, task := range c.taskSlice {
		waitGroup.Add(1)
		currentTask := task

		go func() {
			defer waitGroup.Done()
			if err := currentTask(ctx); err != nil {
				select {
				case errorChannel <- err:
				case <-ctx.Done():
				}
			}
		}()
	}

	waitGroup.Wait()
	close(errorChannel)

	for err := range errorChannel {
		return err
	}

	return nil
}
