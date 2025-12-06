package conveyer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrChanNotFound = errors.New("chan not found")

type Conveyer struct {
	chans map[string]chan string
	cap   int
	funcs []func(ctx context.Context) error
}

func New(capacity int) *Conveyer {
	return &Conveyer{
		chans: make(map[string]chan string),
		cap:   capacity,
		funcs: make([]func(ctx context.Context) error, 0),
	}
}

func (c *Conveyer) obtainChan(identifier string) chan string {
	if ch, ok := c.chans[identifier]; ok {
		return ch
	}

	channel := make(chan string, c.cap)
	c.chans[identifier] = channel

	return channel
}

func (c *Conveyer) fetchChan(identifier string) (chan string, error) {
	channel, ok := c.chans[identifier]
	if !ok {
		return nil, ErrChanNotFound
	}

	return channel, nil
}

func (c *Conveyer) RegisterDecorator(
	processor func(ctx context.Context, src chan string, dst chan string) error,
	srcID string,
	dstID string,
) {
	c.obtainChan(srcID)
	c.obtainChan(dstID)

	task := func(ctx context.Context) error {
		src := c.obtainChan(srcID)
		dst := c.obtainChan(dstID)

		return processor(ctx, src, dst)
	}

	c.funcs = append(c.funcs, task)
}

func (c *Conveyer) RegisterMultiplexer(
	merger func(ctx context.Context, srcs []chan string, dst chan string) error,
	srcIDs []string,
	dstID string,
) {
	for _, id := range srcIDs {
		c.obtainChan(id)
	}

	c.obtainChan(dstID)

	task := func(ctx context.Context) error {
		srcs := make([]chan string, len(srcIDs))
		for i, id := range srcIDs {
			srcs[i] = c.obtainChan(id)
		}

		dst := c.obtainChan(dstID)

		return merger(ctx, srcs, dst)
	}

	c.funcs = append(c.funcs, task)
}

func (c *Conveyer) RegisterSeparator(
	splitter func(ctx context.Context, src chan string, dsts []chan string) error,
	srcID string,
	dstIDs []string,
) {
	c.obtainChan(srcID)

	for _, id := range dstIDs {
		c.obtainChan(id)
	}

	task := func(ctx context.Context) error {
		src := c.obtainChan(srcID)
		dsts := make([]chan string, len(dstIDs))

		for i, id := range dstIDs {
			dsts[i] = c.obtainChan(id)
		}

		return splitter(ctx, src, dsts)
	}

	c.funcs = append(c.funcs, task)
}

func (c *Conveyer) Send(id string, val string) error {
	channel, err := c.fetchChan(id)
	if err != nil {
		return fmt.Errorf("send: %w", ErrChanNotFound)
	}

	channel <- val

	return nil
}

func (c *Conveyer) Recv(id string) (string, error) {
	channel, err := c.fetchChan(id)
	if err != nil {
		return "", fmt.Errorf("recv: %w", ErrChanNotFound)
	}

	value, ok := <-channel
	if !ok {
		return "", nil
	}

	return value, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	waitGroup := &sync.WaitGroup{}
	errChan := make(chan error, len(c.funcs))
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	for _, proc := range c.funcs {
		waitGroup.Add(1)

		go func(p func(ctx context.Context) error) {
			defer waitGroup.Done()

			if err := p(ctx); err != nil {
				select {
				case errChan <- err:
					cancel()
				default:
				}
			}
		}(proc)
	}

	done := make(chan struct{})
	go func() {
		waitGroup.Wait()
		close(done)
	}()

	select {
	case err := <-errChan:
		<-done
		return err
	case <-ctx.Done():
		<-done
		return fmt.Errorf("run: %w", ctx.Err())
	case <-done:
		return nil
	}
}
