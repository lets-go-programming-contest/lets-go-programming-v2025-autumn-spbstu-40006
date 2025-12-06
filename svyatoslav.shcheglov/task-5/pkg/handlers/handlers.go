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
	Undef      = "undefined"
)

func PrefixDecoratorFunc(ctx context.Context, src chan string, dst chan string) error {
	defer close(dst)

	for {
		select {
		case <-ctx.Done():
			return nil
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
				return nil
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

	if len(dsts) == 0 {
		return nil
	}

	pos := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-src:
			if !ok {
				return nil
			}

			idx := pos % len(dsts)
			pos++

			select {
			case dsts[idx] <- msg:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, srcs []chan string, dst chan string) error {
	defer close(dst)

	if len(srcs) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}

	for _, s := range srcs {
		wg.Add(1)
		srcCopy := s

		reader := func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-srcCopy:
					if !ok {
						return
					}

					if strings.Contains(msg, NoMultiMsg) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case dst <- msg:
					}
				}
			}
		}

		go reader()
	}

	wg.Wait()
	return nil
}
