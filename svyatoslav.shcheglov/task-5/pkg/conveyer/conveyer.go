package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChannelNotFound = errors.New("channel not found")

type Conveyer struct {
	channels map[string]chan string
	tasks    []func(ctx context.Context) error
	mutex    sync.Mutex
}

func NewConveyer() *Conveyer {
	return &Conveyer{
		channels: make(map[string]chan string),
		tasks:    make([]func(ctx context.Context) error, 0),
	}
}

func (c *Conveyer) obtainChannel(channelID string) chan string {
	c.mutex.Lock()
	channel, exists := c.channels[channelID]
	if !exists {
		channel = make(chan string, 32)
		c.channels[channelID] = channel
	}
	c.mutex.Unlock()
	return channel
}

func (c *Conveyer) fetchChannel(channelID string) (chan string, error) {
	c.mutex.Lock()
	channel, exists := c.channels[channelID]
	c.mutex.Unlock()
	if !exists {
		return nil, ErrChannelNotFound
	}
	return channel, nil
}

func (c *Conveyer) RegisterDecorator(
	processor func(ctx context.Context, src chan string, dst chan string) error,
	srcID string,
	dstID string,
) {
	task := func(ctx context.Context) error {
		sourceChannel := c.obtainChannel(srcID)
		destinationChannel := c.obtainChannel(dstID)
		return processor(ctx, sourceChannel, destinationChannel)
	}
	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterSeparator(
	processor func(ctx context.Context, inputChannel chan string, outputTrueChannel chan string, outputFalseChannel chan string) error,
	inputID string,
	outputTrueID string,
	outputFalseID string,
) {
	task := func(ctx context.Context) error {
		inputChannel := c.obtainChannel(inputID)
		outputTrueChannel := c.obtainChannel(outputTrueID)
		outputFalseChannel := c.obtainChannel(outputFalseID)
		return processor(ctx, inputChannel, outputTrueChannel, outputFalseChannel)
	}
	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterMultiplexer(
	processor func(ctx context.Context, leftInputChannel chan string, rightInputChannel chan string, outputChannel chan string) error,
	leftID string,
	rightID string,
	outputID string,
) {
	task := func(ctx context.Context) error {
		leftInputChannel := c.obtainChannel(leftID)
		rightInputChannel := c.obtainChannel(rightID)
		outputChannel := c.obtainChannel(outputID)
		return processor(ctx, leftInputChannel, rightInputChannel, outputChannel)
	}
	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) Send(channelID string, value string) error {
	channel := c.obtainChannel(channelID)
	channel <- value
	return nil
}

func (c *Conveyer) Recv(channelID string) (string, error) {
	channel, err := c.fetchChannel(channelID)
	if err != nil {
		return "", err
	}
	value, ok := <-channel
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
		go func(taskFunc func(ctx context.Context) error) {
			defer taskGroup.Done()
			err := taskFunc(ctx)
			if err != nil {
				errorsChannel <- err
			}
		}(currentTask)
	}

	taskGroup.Wait()
	close(errorsChannel)

	for err := range errorsChannel {
		return err
	}

	return nil
}
