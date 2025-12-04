package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

const UndefinedData = "undefined"

type Conveyer struct {
	buffer     int
	pipes      map[string]chan string
	processors []func(context.Context) error
}

func New(size int) *Conveyer {
	return &Conveyer{
		buffer:     size,
		pipes:      make(map[string]chan string),
		processors: make([]func(context.Context) error, 0),
	}
}

func (c *Conveyer) createPipe(name string) {
	if _, ok := c.pipes[name]; !ok {
		c.pipes[name] = make(chan string, c.buffer)
	}
}

func (c *Conveyer) createPipes(names ...string) {
	for _, n := range names {
		c.createPipe(n)
	}
}

func (c *Conveyer) RegisterDecorator(
	handler func(context.Context, chan string, chan string) error,
	input string,
	output string,
) {
	c.createPipes(input, output)

	c.processors = append(c.processors, func(ctx context.Context) error {
		return handler(ctx, c.pipes[input], c.pipes[output])
	})
}

func (c *Conveyer) RegisterMultiplexer(
	handler func(context.Context, []chan string, chan string) error,
	inputs []string,
	output string,
) {
	c.createPipes(inputs...)

	c.createPipes(output)

	c.processors = append(c.processors, func(ctx context.Context) error {
		var inStreams []chan string
		for _, name := range inputs {
			inStreams = append(inStreams, c.pipes[name])
		}

		return handler(ctx, inStreams, c.pipes[output])
	})
}

func (c *Conveyer) RegisterSeparator(
	handler func(context.Context, chan string, []chan string) error,
	input string,
	outputs []string,
) {
	c.createPipes(input)

	c.createPipes(outputs...)

	c.processors = append(c.processors, func(ctx context.Context) error {
		var outStreams []chan string
		for _, name := range outputs {
			outStreams = append(outStreams, c.pipes[name])
		}

		return handler(ctx, c.pipes[input], outStreams)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	var waitGroup sync.WaitGroup

	errCh := make(chan error, len(c.processors))

	for _, processor := range c.processors {
		waitGroup.Add(1)

		go func(proc func(context.Context) error) {
			defer waitGroup.Done()

			if err := proc(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(processor)
	}

	go func() {
		waitGroup.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Conveyer) Send(name string, msg string) error {
	pipe, ok := c.pipes[name]

	if !ok {
		return ErrChanNotFound
	}

	pipe <- msg

	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	pipe, ok := c.pipes[name]
	if !ok {
		return "", ErrChanNotFound
	}

	val, open := <-pipe
	if !open {
		return UndefinedData, nil
	}

	return val, nil
}
