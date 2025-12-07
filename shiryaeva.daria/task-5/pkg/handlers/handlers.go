package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

const (
	prefix            = "decorated: "
	noDecoratorText   = "no decorator"
	noMultiplexerText = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecoratorText) {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
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

	if len(outputs) == 0 {
		return nil
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			outputIdx := index % len(outputs)
			index++

			select {
			case <-ctx.Done():
				return nil
			case outputs[outputIdx] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}

	for _, inputChan := range inputs {
		wg.Add(1)

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

					if strings.Contains(data, noMultiplexerText) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- data:
					}
				}
			}
		}(inputChan)
	}

	wg.Wait()

	return nil
}
