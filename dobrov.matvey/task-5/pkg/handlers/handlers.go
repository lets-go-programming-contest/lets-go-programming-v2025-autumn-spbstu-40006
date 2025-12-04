package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

const (
	Prefix               = "decorated: "
	NoDecoratorMessage   = "no decorator"
	NoMultiplexerMessage = "no multiplexer"
	Undefined            = "undefined"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, NoDecoratorMessage) {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(data, Prefix) {
				data = Prefix + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, ch := range outputs {
			close(ch)
		}
	}()

	index := 0
	n := len(outputs)

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			out := outputs[index]

			index = (index + 1) % n

			select {
			case <-ctx.Done():
				return nil
			case out <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, in := range inputs {
		in := in

		go func(ch chan string) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case data, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(data, NoMultiplexerMessage) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- data:
					}
				}
			}
		}(in)
	}

	wg.Wait()
	return nil
}
