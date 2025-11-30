package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

const noMultiplexerData = "no multiplexer"

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context error: %w", ctx.Err())
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
				return fmt.Errorf("context error: %w", ctx.Err())
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, outputChannel := range outputs {
			close(outputChannel)
		}
	}()

	var counter uint64

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context error: %w", ctx.Err())
		case data, isChannelOpen := <-input:
			if !isChannelOpen {
				return nil
			}

			index := atomic.AddUint64(&counter, 1) % uint64(len(outputs))

			select {
			case <-ctx.Done():
				return fmt.Errorf("context error: %w", ctx.Err())
			case outputs[index] <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	for _, inputChannel := range inputs {
		inChannel := inputChannel

		go func(chLocal chan string) {
			for {
				select {
				case <-ctx.Done():
					return

				case line, okay := <-chLocal:
					if !okay {
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
		}(inChannel)
	}

	<-ctx.Done()

	return nil
}
