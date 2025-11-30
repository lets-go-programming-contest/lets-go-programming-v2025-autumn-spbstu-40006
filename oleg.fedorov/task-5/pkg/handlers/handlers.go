package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
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
				return fmt.Errorf("context canceled: %w", ctx.Err())
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
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case data, isChannelOpen := <-input:
			if !isChannelOpen {
				return nil
			}

			index := atomic.AddUint64(&counter, 1) % uint64(len(outputs))

			select {
			case <-ctx.Done():
				return fmt.Errorf("context canceled: %w", ctx.Err())
			case outputs[index] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	type channelResult struct {
		data string
		ok   bool
	}

	results := make(chan channelResult, len(inputs))

	for _, inputChannel := range inputs {
		go func(in chan string) {
			for {
				select {
				case <-ctx.Done():
					return
				case data, isChannelOpen := <-in:
					if !isChannelOpen {
						return
					}

					select {
					case <-ctx.Done():
						return
					case results <- channelResult{data, isChannelOpen}:
					}
				}
			}
		}(inputChannel)
	}

	activeInputs := len(inputs)

	for activeInputs > 0 {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case result := <-results:
			if !result.ok {
				activeInputs--
				continue
			}

			if strings.Contains(result.data, "no multiplexer") {
				continue
			}

			select {
			case <-ctx.Done():
				return fmt.Errorf("context canceled: %w", ctx.Err())
			case output <- result.data:
			}
		}
	}

	return nil
}
