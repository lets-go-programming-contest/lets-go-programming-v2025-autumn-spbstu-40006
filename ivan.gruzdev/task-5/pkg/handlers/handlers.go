package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotDecorate = errors.New("unacceptable for decoration")

const (
	noDecorateMark   = "no decorator"
	decoratePrefix   = "decorated: "
	muxSkipIndicator = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case line, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(line, noDecorateMark) {
				return ErrCannotDecorate
			}

			if !strings.HasPrefix(line, decoratePrefix) {
				line = decoratePrefix + line
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- line:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, o := range outputs {
			close(o)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	pos := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[pos%len(outputs)] <- msg:
				pos++
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	var wg sync.WaitGroup

	for _, ch := range inputs {
		local := ch
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-local:
					if !ok {
						return
					}

					if strings.Contains(msg, muxSkipIndicator) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- msg:
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
