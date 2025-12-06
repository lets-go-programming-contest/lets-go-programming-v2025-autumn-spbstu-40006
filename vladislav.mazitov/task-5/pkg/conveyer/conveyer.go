package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

const UndefinedData = "undefined"

type Conveyer struct {
	buffer   int
	channels map[string]chan string
	handlers []func(ctx context.Context) error
}

func New(size int) *Conveyer {
	newConveyer := &Conveyer{
		buffer:   size,
		channels: make(map[string]chan string),
		handlers: make([]func(context.Context) error, 0),
	}

	return newConveyer
}

func (c *Conveyer) makeChannel(name string) {
	_, ok := c.channels[name]
	if !ok {
		c.channels[name] = make(chan string, c.buffer)
	}
}

func (c *Conveyer) makeChannels(names ...string) {
	for _, name := range names {
		c.makeChannel(name)
	}
}

func (c *Conveyer) RegisterDecorator(
	handler func(ctx context.Context, input chan string, output chan string) error,
	inputName string,
	outputName string,
) {
	c.makeChannels(inputName)
	c.makeChannels(outputName)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		return handler(ctx, c.channels[inputName], c.channels[outputName])
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(ctx context.Context, inputs []chan string, output chan string) error,
	inputNames []string,
	outputName string,
) {
	c.makeChannels(inputNames...)
	c.makeChannels(outputName)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		inputChannels := make([]chan string, 0, len(inputNames))
		for _, input := range inputNames {
			inputChannels = append(inputChannels, c.channels[input])
		}

		return handler(ctx, inputChannels, c.channels[outputName])
	})
}

func (c *Conveyer) RegisterSeparator(
	handler func(ctx context.Context, input chan string, outputs []chan string) error,
	inputName string,
	outputNames []string,
) {
	c.makeChannels(inputName)
	c.makeChannels(outputNames...)

	c.handlers = append(c.handlers, func(ctx context.Context) error {
		outputChannels := make([]chan string, 0, len(outputNames))
		for _, output := range outputNames {
			outputChannels = append(outputChannels, c.channels[output])
		}

		return handler(ctx, c.channels[inputName], outputChannels)
	})
}

func (c *Conveyer) Send(input string, data string) error {
	channel, ok := c.channels[input]
	if !ok {
		return ErrChanNotFound
	}

	channel <- data

	return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
	channel, found := c.channels[output]
	if !found {
		return "", ErrChanNotFound
	}

	data, ok := <-channel
	if !ok {
		return UndefinedData, nil
	}

	return data, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	waitGroup := sync.WaitGroup{}
	errorsChannel := make(chan error, len(c.handlers))

	for _, handler := range c.handlers {
		waitGroup.Add(1)

		tempHandler := handler
		go func() {
			defer waitGroup.Done()

			if err := tempHandler(ctx); err != nil {
				select {
				case errorsChannel <- err:
				default:
				}
			}
		}()
	}

	go func() {
		waitGroup.Wait()
		close(errorsChannel)
	}()

	select {
	case err := <-errorsChannel:
		return err
	case <-ctx.Done():
		waitGroup.Wait()

		return nil
	}
}
