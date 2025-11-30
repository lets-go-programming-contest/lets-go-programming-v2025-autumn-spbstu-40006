package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

const undefinedData = "undefined"

var ErrChanNotFound = errors.New("chan not found")

type Conveyer interface {
	RegisterDecorator(
		fn func(ctx context.Context, input chan string, output chan string) error,
		input string,
		output string,
	)

	RegisterMultiplexer(
		fn func(ctx context.Context, inputs []chan string, output chan string) error,
		inputs []string,
		output string,
	)

	RegisterSeparator(
		fn func(ctx context.Context, input chan string, outputs []chan string) error,
		input string,
		outputs []string,
	)

	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

type decoratorReg struct {
	fn     func(context.Context, chan string, chan string) error
	input  string
	output string
}

type multiplexerReg struct {
	fn     func(context.Context, []chan string, chan string) error
	inputs []string
	output string
}

type separatorReg struct {
	fn      func(context.Context, chan string, []chan string) error
	input   string
	outputs []string
}

type conveyer struct {
	size         int
	chans        map[string]chan string
	decorators   []decoratorReg
	multiplexers []multiplexerReg
	separators   []separatorReg
	mu           sync.RWMutex
}

func New(size int) Conveyer {
	return &conveyer{
		size:  size,
		chans: make(map[string]chan string),
	}
}

func (c *conveyer) RegisterDecorator(
	fn func(ctx context.Context, input chan string, output chan string) error,
	input string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.chans[input]; !ok {
		c.chans[input] = make(chan string, c.size)
	}
	if _, ok := c.chans[output]; !ok {
		c.chans[output] = make(chan string, c.size)
	}

	c.decorators = append(c.decorators, decoratorReg{fn: fn, input: input, output: output})
}

func (c *conveyer) RegisterMultiplexer(
	fn func(ctx context.Context, inputs []chan string, output chan string) error,
	inputs []string,
	output string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, input := range inputs {
		if _, ok := c.chans[input]; !ok {
			c.chans[input] = make(chan string, c.size)
		}
	}
	if _, ok := c.chans[output]; !ok {
		c.chans[output] = make(chan string, c.size)
	}

	c.multiplexers = append(c.multiplexers, multiplexerReg{fn: fn, inputs: inputs, output: output})
}

func (c *conveyer) RegisterSeparator(
	fn func(ctx context.Context, input chan string, outputs []chan string) error,
	input string,
	outputs []string,
) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.chans[input]; !ok {
		c.chans[input] = make(chan string, c.size)
	}
	for _, output := range outputs {
		if _, ok := c.chans[output]; !ok {
			c.chans[output] = make(chan string, c.size)
		}
	}

	c.separators = append(c.separators, separatorReg{fn: fn, input: input, outputs: outputs})
}

func (c *conveyer) Send(input string, data string) error {
	c.mu.RLock()
	ch, ok := c.chans[input]
	c.mu.RUnlock()

	if !ok {
		return ErrChanNotFound
	}
	ch <- data
	return nil
}

func (c *conveyer) Recv(output string) (string, error) {
	c.mu.RLock()
	ch, ok := c.chans[output]
	c.mu.RUnlock()

	if !ok {
		return "", ErrChanNotFound
	}

	val, ok := <-ch
	if !ok {
		return undefinedData, nil
	}
	return val, nil
}

func (c *conveyer) startDecorators(ctx context.Context, start func(fn func() error)) {
	for _, d := range c.decorators {
		in := c.chans[d.input]
		out := c.chans[d.output]
		start(func() error {
			return d.fn(ctx, in, out)
		})
	}
}

func (c *conveyer) startMultiplexers(ctx context.Context, start func(fn func() error)) {
	for _, m := range c.multiplexers {
		var ins []chan string
		for _, id := range m.inputs {
			ins = append(ins, c.chans[id])
		}
		out := c.chans[m.output]
		start(func() error {
			return m.fn(ctx, ins, out)
		})
	}
}

func (c *conveyer) startSeparators(ctx context.Context, start func(fn func() error)) {
	for _, s := range c.separators {
		in := c.chans[s.input]
		var outs []chan string
		for _, id := range s.outputs {
			outs = append(outs, c.chans[id])
		}
		start(func() error {
			return s.fn(ctx, in, outs)
		})
	}
}

func (c *conveyer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	start := func(fn func() error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := fn(); err != nil {
				select {
				case errCh <- err:
					cancel()
				default:
				}
			}
		}()
	}

	c.startDecorators(ctx, start)
	c.startMultiplexers(ctx, start)
	c.startSeparators(ctx, start)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("handler error: %w", err)
		}
	case <-done:
	}

	c.mu.Lock()
	for _, ch := range c.chans {
		close(ch)
	}
	c.mu.Unlock()

	return nil
}
