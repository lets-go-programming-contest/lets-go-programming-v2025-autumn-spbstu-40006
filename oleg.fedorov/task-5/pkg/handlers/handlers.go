package handlers

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, isChannelOpen := <-input:
			if !isChannelOpen {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	var counter uint64

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, isChannelOpen := <-input:
			if !isChannelOpen {
				return nil
			}

			index := atomic.AddUint64(&counter, 1) % uint64(len(outputs))

			select {
			case <-ctx.Done():
				return ctx.Err()
			case outputs[index] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		for index := range inputs {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case data, isChannelOpen := <-inputs[index]:
				if !isChannelOpen {
					continue
				}

				if !strings.Contains(data, "no multiplexer") {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case output <- data:
					}
				}
			default:
			}
		}
	}
}
