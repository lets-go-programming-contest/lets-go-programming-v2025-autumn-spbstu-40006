package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantBeDecorated = errors.New("can't be decorated")

const (
	noDecoratorData   = "no decorator"
	textForDecorating = "decorated: "
	noMultiplexerData = "no multiplexer"
)

func PrefixDecoratorFunction(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case line, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(line, noDecoratorData) {
				return ErrCantBeDecorated
			}

			if !strings.HasPrefix(line, textForDecorating) {
				line = textForDecorating + line
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- line:
			}
		}
	}
}

func MultiplexerFunction(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	for _, input := range inputs {
		tempInput := input

		go func(in chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case line, ok := <-in:
					if !ok {
						return
					}

					if strings.Contains(line, noMultiplexerData) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- line:
					}
				}
			}
		}(tempInput)
	}

	<-ctx.Done()

	return nil
}

func SeparatorFunction(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, output := range outputs {
			close(output)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	cnt := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case line, ok := <-input:
			if !ok {
				return nil
			}

			index := cnt % len(outputs)
			cnt++

			select {
			case outputs[index] <- line:
			case <-ctx.Done():
				return nil
			}
		}
	}
}
