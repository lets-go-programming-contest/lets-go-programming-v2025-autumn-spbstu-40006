package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChannelNotFound = errors.New("channel_not_found")

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
	defer c.mutex.Unlock()

	ch, exists := c.channelMap[identifier]
	if !exists {
		ch = make(chan string, c.bufferSize)
		c.channelMap[identifier] = ch
	}
	return ch
}

func (c *Conveyer) fetchChannel(identifier string) (chan string, error) {
	c.mutex.Lock()
	ch, exists := c.channelMap[identifier]
	c.mutex.Unlock()

	if !exists {
		return nil, ErrChannelNotFound
	}
	return ch, nil
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

func (c *Conveyer) Send(identifier string, value string) error {
	ch := c.obtainChannel(identifier)
	ch <- value
	return nil
}

func (c *Conveyer) Recv(identifier string) (string, error) {
	ch, err := c.fetchChannel(identifier)
	if err != nil {
		return "", err
	}
	val, open := <-ch
	if !open {
		return "", nil
	}
	return val, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(c.taskSlice))

	for _, task := range c.taskSlice {
		wg.Add(1)
		t := task
		go func() {
			defer wg.Done()
			if err := t(ctx); err != nil {
				errorChannel <- err
			}
		}()
	}

	wg.Wait()
	close(errorChannel)

	for err := range errorChannel {
		return err
	}

	return nil
}
