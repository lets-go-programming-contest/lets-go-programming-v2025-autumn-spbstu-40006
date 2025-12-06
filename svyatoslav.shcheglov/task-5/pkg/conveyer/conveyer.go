package conveyer

import (
	"context"
	"errors"
	"sync"
)

const defaultChannelBufferSize = 16

var ErrChanNotFound = errors.New("channel not found")

type Conveyer struct {
	channels map[string]chan string
	tasks    []func(ctx context.Context) error
	mutex    sync.Mutex
}

func New() *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		tasks:    make([]func(ctx context.Context) error, 0),
		mutex:    sync.Mutex{},
	}
}

func (c *Conveyer) obtainChannel(channelIdentifier string) chan string {
	c.mutex.Lock()
	channel, exists := c.channels[channelIdentifier]
	if !exists {
		channel = make(chan string, defaultChannelBufferSize)
		c.channels[channelIdentifier] = channel
	}
	c.mutex.Unlock()
	return channel
}

func (c *Conveyer) RegisterDecorator(
	processor func(ctx context.Context, src chan string, dst chan string) error,
	sourceChannelID string,
	destinationChannelID string,
) {
	task := func(ctx context.Context) error {
		sourceChannel := c.obtainChannel(sourceChannelID)
		destinationChannel := c.obtainChannel(destinationChannelID)
		return processor(ctx, sourceChannel, destinationChannel)
	}
	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterSeparator(
	processor func(
		ctx context.Context,
		inputChannel chan string,
		outputTrueChannel chan string,
		outputFalseChannel chan string,
	) error,
	inputChannelID string,
	outputTrueChannelID string,
	outputFalseChannelID string,
) {
	task := func(ctx context.Context) error {
		inputChannel := c.obtainChannel(inputChannelID)
		outputTrueChannel := c.obtainChannel(outputTrueChannelID)
		outputFalseChannel := c.obtainChannel(outputFalseChannelID)
		return processor(ctx, inputChannel, outputTrueChannel, outputFalseChannel)
	}
	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterMultiplexer(
	processor func(
		ctx context.Context,
		inputLeftChannel chan string,
		inputRightChannel chan string,
		outputChannel chan string,
	) error,
	leftInputChannelID string,
	rightInputChannelID string,
	outputChannelID string,
) {
	task := func(ctx context.Context) error {
		leftInputChannel := c.obtainChannel(leftInputChannelID)
		rightInputChannel := c.obtainChannel(rightInputChannelID)
		outputChannel := c.obtainChannel(outputChannelID)
		return processor(ctx, leftInputChannel, rightInputChannel, outputChannel)
	}
	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) Send(channelID string, value string) error {
	targetChannel := c.obtainChannel(channelID)
	targetChannel <- value
	return nil
}

func (c *Conveyer) Recv(channelID string) (string, error) {
	targetChannel := c.obtainChannel(channelID)
	value, ok := <-targetChannel
	if !ok {
		return "", nil
	}
	return value, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	var taskGroup sync.WaitGroup
	errorsChannel := make(chan error, len(c.tasks))

	for _, task := range c.tasks {
		taskGroup.Add(1)
		currentTask := task
		go func() {
			defer taskGroup.Done()
			if err := currentTask(ctx); err != nil {
				errorsChannel <- err
			}
		}()
	}

	taskGroup.Wait()
	close(errorsChannel)

	for err := range errorsChannel {
		return err
	}

	return nil
}
