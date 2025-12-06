package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrAlreadyRunning  = errors.New("conveyer is already running")
	ErrChanNotFound    = errors.New("chan not found")
	ErrNoDataAvailable = errors.New("no data available")
	ErrChanFull        = errors.New("channel is full")
)

type conveyerImpl struct {
	mu           sync.RWMutex
	channels     map[string]chan string
	decorators   []decoratorHandler
	multiplexers []multiplexerHandler
	separators   []separatorHandler
	size         int
	running      bool
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	errChan      chan error
}

type decoratorHandler struct {
	fn     DecoratorFunc
	input  string
	output string
}

type multiplexerHandler struct {
	fn     MultiplexerFunc
	inputs []string
	output string
}

type separatorHandler struct {
	fn      SeparatorFunc
	input   string
	outputs []string
}

func newConveyerImpl(size int) *conveyerImpl {
	const errChanBufferSize = 10

	return &conveyerImpl{
		channels:     make(map[string]chan string),
		size:         size,
		decorators:   make([]decoratorHandler, 0),
		multiplexers: make([]multiplexerHandler, 0),
		separators:   make([]separatorHandler, 0),
		errChan:      make(chan error, errChanBufferSize),
		mu:           sync.RWMutex{},
		running:      false,
		cancel:       nil,
		wg:           sync.WaitGroup{},
	}
}

func (c *conveyerImpl) RegisterDecorator(function DecoratorFunc, input string, output string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureChannel(input)
	c.ensureChannel(output)

	c.decorators = append(c.decorators, decoratorHandler{
		fn:     function,
		input:  input,
		output: output,
	})
}

func (c *conveyerImpl) RegisterMultiplexer(function MultiplexerFunc, inputs []string, output string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, inputName := range inputs {
		c.ensureChannel(inputName)
	}

	c.ensureChannel(output)

	c.multiplexers = append(c.multiplexers, multiplexerHandler{
		fn:     function,
		inputs: inputs,
		output: output,
	})
}

func (c *conveyerImpl) RegisterSeparator(function SeparatorFunc, input string, outputs []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ensureChannel(input)

	for _, outputName := range outputs {
		c.ensureChannel(outputName)
	}

	c.separators = append(c.separators, separatorHandler{
		fn:      function,
		input:   input,
		outputs: outputs,
	})
}

func (c *conveyerImpl) ensureChannel(name string) {
	if _, exists := c.channels[name]; !exists {
		c.channels[name] = make(chan string, c.size)
	}
}

func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.Lock()

	if c.running {
		c.mu.Unlock()

		return fmt.Errorf("%w", ErrAlreadyRunning)
	}

	c.running = true

	runCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	c.mu.Unlock()

	defer func() {
		cancel()
		c.stop()
	}()

	c.startHandlers(runCtx)

	select {
	case <-runCtx.Done():
		return nil

	case err := <-c.errChan:
		return err
	}
}

func (c *conveyerImpl) startHandlers(ctx context.Context) {
	c.startDecorators(ctx)
	c.startMultiplexers(ctx)
	c.startSeparators(ctx)
}

func (c *conveyerImpl) startDecorators(ctx context.Context) {
	for _, handler := range c.decorators {
		c.wg.Add(1)

		go func(decorator decoratorHandler) {
			defer c.wg.Done()

			inputChan := c.getChannel(decorator.input)
			outputChan := c.getChannel(decorator.output)

			if err := decorator.fn(ctx, inputChan, outputChan); err != nil {
				select {
				case c.errChan <- err:
				default:
				}
			}
		}(handler)
	}
}

func (c *conveyerImpl) startMultiplexers(ctx context.Context) {
	for _, handler := range c.multiplexers {
		c.wg.Add(1)

		go func(multiplexer multiplexerHandler) {
			defer c.wg.Done()

			inputChans := make([]chan string, len(multiplexer.inputs))

			for index, inputName := range multiplexer.inputs {
				inputChans[index] = c.getChannel(inputName)
			}

			outputChan := c.getChannel(multiplexer.output)

			if err := multiplexer.fn(ctx, inputChans, outputChan); err != nil {
				select {
				case c.errChan <- err:
				default:
				}
			}
		}(handler)
	}
}

func (c *conveyerImpl) startSeparators(ctx context.Context) {
	for _, handler := range c.separators {
		c.wg.Add(1)

		go func(separator separatorHandler) {
			defer c.wg.Done()

			inputChan := c.getChannel(separator.input)
			outputChans := make([]chan string, len(separator.outputs))

			for index, outputName := range separator.outputs {
				outputChans[index] = c.getChannel(outputName)
			}

			if err := separator.fn(ctx, inputChan, outputChans); err != nil {
				select {
				case c.errChan <- err:
				default:
				}
			}
		}(handler)
	}
}

func (c *conveyerImpl) stop() {
	c.mu.Lock()
	c.running = false
	channels := make(map[string]chan string)

	for key, value := range c.channels {
		channels[key] = value
	}
	c.mu.Unlock()

	for _, channel := range channels {
		close(channel)
	}

	c.wg.Wait()
}

func (c *conveyerImpl) getChannel(name string) chan string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.channels[name]
}

func (c *conveyerImpl) Send(input string, data string) error {
	channel := c.getChannel(input)
	if channel == nil {
		return fmt.Errorf("%w", ErrChanNotFound)
	}

	select {
	case channel <- data:
		return nil
	default:
		return nil
	}
}

func (c *conveyerImpl) Recv(output string) (string, error) {
	channel := c.getChannel(output)
	if channel == nil {
		return "", fmt.Errorf("%w", ErrChanNotFound)
	}

	select {
	case data, ok := <-channel:
		if !ok {
			return "undefined", nil
		}

		return data, nil

	default:
		return "", fmt.Errorf("%w", ErrNoDataAvailable)
	}
}
