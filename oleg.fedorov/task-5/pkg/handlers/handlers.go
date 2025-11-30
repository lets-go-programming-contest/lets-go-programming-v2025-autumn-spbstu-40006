package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

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

	aggregateChannel := make(chan string, len(inputs))

	var waitGroup sync.WaitGroup

	for index := range inputs {
		waitGroup.Add(1)

		go func(inputIndex int) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, isChannelOpen := <-inputs[inputIndex]:
					if !isChannelOpen {
						return
					}

					if strings.Contains(data, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case aggregateChannel <- data:
					}
				}
			}
		}(index)
	}

	go func() {
		waitGroup.Wait()
		close(aggregateChannel)
	}()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context error: %w", ctx.Err())
		case data, isChannelOpen := <-aggregateChannel:
			if !isChannelOpen {
				return nil
			}

			select {
			case <-ctx.Done():
				return fmt.Errorf("context error: %w", ctx.Err())
			case output <- data:
			}
		}
	}
}
