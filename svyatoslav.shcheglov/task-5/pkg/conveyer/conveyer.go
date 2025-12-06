package conveyer

import (
	"context"
	"fmt"
	"sync"
)

var ErrChanNotFound = fmt.Errorf("channel not found")

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

func (c *Conveyer) obtainChannel(id string) chan string {
	c.mutex.Lock()
	channel, exists := c.channels[id]
	if !exists {
		channel = make(chan string, 16)
		c.channels[id] = channel
	}
	c.mutex.Unlock()
	return channel
}

func (c *Conveyer) RegisterDecorator(
	processor func(ctx context.Context, src chan string, dst chan string) error,
	srcID string,
	dstID string,
) {
	task := func(ctx context.Context) error {
		src := c.obtainChannel(srcID)
		dst := c.obtainChannel(dstID)
		return processor(ctx, src, dst)
	}
	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterSeparator(
	processor func(ctx context.Context, input chan string, outputTrue chan string, outputFalse chan string) error,
	inputID string,
	outputTrueID string,
	outputFalseID string,
) {
	task := func(ctx context.Context) error {
		in := c.obtainChannel(inputID)
		outTrue := c.obtainChannel(outputTrueID)
		outFalse := c.obtainChannel(outputFalseID)
		return processor(ctx, in, outTrue, outFalse)
	}
	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) RegisterMultiplexer(
	processor func(ctx context.Context, inputLeft chan string, inputRight chan string, output chan string) error,
	leftID string,
	rightID string,
	outputID string,
) {
	task := func(ctx context.Context) error {
		left := c.obtainChannel(leftID)
		right := c.obtainChannel(rightID)
		out := c.obtainChannel(outputID)
		return processor(ctx, left, right, out)
	}
	c.tasks = append(c.tasks, task)
}

func (c *Conveyer) Send(id string, value string) error {
	channel := c.obtainChannel(id)
	channel <- value
	return nil
}

func (c *Conveyer) Recv(id string) (string, error) {
	channel := c.obtainChannel(id)
	value, ok := <-channel
	if !ok {
		return "", nil
	}
	return value, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	var taskGroup sync.WaitGroup
	errorsChan := make(chan error, len(c.tasks))

	for _, task := range c.tasks {
		taskGroup.Add(1)
		go func(fn func(context.Context) error) {
			defer taskGroup.Done()
			if err := fn(ctx); err != nil {
				errorsChan <- err
			}
		}(task)
	}

	taskGroup.Wait()
	close(errorsChan)

	for err := range errorsChan {
		return err
	}

	return nil
}
