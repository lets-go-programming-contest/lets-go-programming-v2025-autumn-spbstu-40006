package conveyer

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

const UndefinedData = "undefined"

var ErrChanNotFound = errors.New("chan not found")

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
		for _, id := range inputs {
			inStreams = append(inStreams, c.pipes[id])
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
		for _, id := range outputs {
			outStreams = append(outStreams, c.pipes[id])
		}

		return handler(ctx, c.pipes[input], outStreams)
	})
}

func (c *Conveyer) Run(ctx context.Context) error {
	group, gCtx := errgroup.WithContext(ctx)

	for _, proc := range c.processors {
		task := proc
		group.Go(func() error {
			return task(gCtx)
		})
	}

	if err := group.Wait(); err != nil {
		return fmt.Errorf("conveyer error: %w", err)
	}

	return nil
}

func (c *Conveyer) Send(name string, msg string) error {
	ch, ok := c.pipes[name]
	if !ok {
		return ErrChanNotFound
	}

	ch <- msg
	return nil
}

func (c *Conveyer) Recv(name string) (string, error) {
	ch, ok := c.pipes[name]
	if !ok {
		return "", ErrChanNotFound
	}

	val, open := <-ch
	if !open {
		return UndefinedData, nil
	}

	return val, nil
}
