package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrNoDecorator = errors.New("can't be decorated")

const (
	noDecoratorKeyword  = "no decorator"
	decoratorPrefix     = "decorated: "
	noMultiplexerFilter = "no multiplexer"
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

			if strings.Contains(line, noDecoratorKeyword) {
				return ErrNoDecorator
			}

			if !strings.HasPrefix(line, decoratorPrefix) {
				line = decoratorPrefix + line
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
		for _, out := range outputs {
			close(out)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	index := 0

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
			case outputs[index%len(outputs)] <- msg:
				index++
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	for _, ch := range inputs {
		local := ch
		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case msg, ok := <-local:
					if !ok {
						return
					}

					if strings.Contains(msg, noMultiplexerFilter) {
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

	waitGroup.Wait()
	return nil
}
