package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

const (
	DecoPrefix = "decorated: "
	NoDecoMsg  = "no decorator"
	NoMultiMsg = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, src chan string, dst chan string) error {
	defer close(dst)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("prefix decorator: %w", ctx.Err())
		case msg, ok := <-src:
			if !ok {
				return nil
			}

			if strings.Contains(msg, NoDecoMsg) {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(msg, DecoPrefix) {
				msg = DecoPrefix + msg
			}

			select {
			case dst <- msg:
			case <-ctx.Done():
				return fmt.Errorf("prefix decorator: %w", ctx.Err())
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, src chan string, dsts []chan string) error {
	defer func() {
		for _, d := range dsts {
			close(d)
		}
	}()

	pos := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("separator: %w", ctx.Err())
		case msg, ok := <-src:
			if !ok {
				return nil
			}

			if len(dsts) == 0 {
				continue
			}

			idx := pos % len(dsts)
			pos++

			select {
			case dsts[idx] <- msg:
			case <-ctx.Done():
				return fmt.Errorf("separator: %w", ctx.Err())
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, srcs []chan string, dst chan string) error {
	defer close(dst)

	if len(srcs) == 0 {
		return nil
	}

	waitGroup := &sync.WaitGroup{}

	for _, srcChan := range srcs {
		waitGroup.Add(1)

		go func(source chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-source:
					if !ok {
						return
					}

					if strings.Contains(msg, NoMultiMsg) {
						continue
					}

					select {
					case dst <- msg:
					case <-ctx.Done():
						return
					}
				}
			}
		}(srcChan)
	}

	waitGroup.Wait()
	return nil
}
