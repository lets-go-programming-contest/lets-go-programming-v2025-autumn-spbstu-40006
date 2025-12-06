package conveyer

import (
	"context"
	"errors"
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

	ch := make(chan string, c.cap)
	c.chans[identifier] = ch

	return ch
}

func (c *Conveyer) fetchChan(identifier string) (chan string, error) {
	ch, ok := c.chans[identifier]
	if !ok {

		return nil, ErrChanNotFound
	}

	return ch, nil
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
	ch, err := c.fetchChan(id)
	if err != nil {

		return ErrChanNotFound
	}

	ch <- val

	return nil
}

func (c *Conveyer) Recv(id string) (string, error) {
	ch, err := c.fetchChan(id)
	if err != nil {

		return "", ErrChanNotFound
	}

	v, ok := <-ch
	if !ok {

		return "", nil

	}

	return v, nil
}

func (c *Conveyer) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	errChan := make(chan error, len(c.funcs))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, proc := range c.funcs {
		wg.Add(1)

		go func(p func(ctx context.Context) error) {
			defer wg.Done()

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
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errChan:
		<-done
		return err
	case <-ctx.Done():
		<-done
		return ctx.Err()
	case <-done:
		return nil
	}
}
