package handlers

import (
	"context"
	"errors"
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

			return ctx.Err()
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

				return ctx.Err()
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

			return ctx.Err()
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

				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, srcs []chan string, dst chan string) error {
	defer close(dst)

	if len(srcs) == 0 {

		return nil
	}

	wg := &sync.WaitGroup{}
	errChan := make(chan error, 1)

	for _, s := range srcs {
		wg.Add(1)

		go func(src chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():

					return
				case msg, ok := <-src:
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
		}(s)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errChan:

		return err
	case <-ctx.Done():

		return ctx.Err()
	case <-done:

		return nil
	}
}
