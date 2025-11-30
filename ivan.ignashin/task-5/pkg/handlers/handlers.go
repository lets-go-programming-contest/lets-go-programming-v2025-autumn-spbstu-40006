package handlers

import (
	"context"
	"errors"
	"strings"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case v, ok := <-input:
			if !ok {
				close(output)
				return nil
			}

			if strings.Contains(v, "no decorator") {
				return errors.New("can't be decorated")
			}

			if !strings.HasPrefix(v, "decorated: ") {
				v = "decorated: " + v
			}

			select {
			case output <- v:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	i := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case v, ok := <-input:
			if !ok {
				for _, ch := range outputs {
					close(ch)
				}
				return nil
			}

			out := outputs[i%len(outputs)]
			i++

			select {
			case out <- v:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	open := len(inputs)

	for open > 0 {
		for _, in := range inputs {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case v, ok := <-in:
				if !ok {
					open--
					continue
				}

				if strings.Contains(v, "no multiplexer") {
					continue
				}

				select {
				case output <- v:
				case <-ctx.Done():
					return ctx.Err()
				}
			default:
			}
		}
	}

	close(output)
	return nil
}
