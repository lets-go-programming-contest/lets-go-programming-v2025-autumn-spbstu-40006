package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const (
		prefix         = "decorated: "
		errorSubstring = "no decorator"
	)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("decorator context error: %w", ctx.Err())

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, errorSubstring) {
				return errors.New("can't be decorated")
			}

			var result string
			if strings.HasPrefix(data, prefix) {
				result = data
			} else {
				result = prefix + data
			}

			select {
			case output <- result:
				continue
			case <-ctx.Done():
				return fmt.Errorf("decorator output error: %w", ctx.Err())
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	var counter uint64

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("separator context error: %w", ctx.Err())

		case data, ok := <-input:
			if !ok {
				return nil
			}

			current := atomic.AddUint64(&counter, 1)
			index := int((current - 1) % uint64(len(outputs)))

			select {
			case outputs[index] <- data:
				continue
			case <-ctx.Done():
				return fmt.Errorf("separator output error: %w", ctx.Err())
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	const filterSubstring = "no multiplexer"

	processInput := func(inputChan chan string, stopChan <-chan struct{}) error {
		for {
			select {
			case <-stopChan:
				return nil
			case <-ctx.Done():
				return fmt.Errorf("multiplexer context error: %w", ctx.Err())
			case data, ok := <-inputChan:
				if !ok {
					return nil
				}

				if strings.Contains(data, filterSubstring) {
					continue
				}

				select {
				case output <- data:
					continue
				case <-ctx.Done():
					return fmt.Errorf("multiplexer output error: %w", ctx.Err())
				case <-stopChan:
					return nil
				}
			}
		}
	}

	stopChan := make(chan struct{})
	errChan := make(chan error, len(inputs))

	for _, inputChan := range inputs {
		go func(ch chan string) {
			errChan <- processInput(ch, stopChan)
		}(inputChan)
	}

	var firstErr error
	errCount := 0

	for range inputs {
		select {
		case err := <-errChan:
			errCount++

			if err != nil && firstErr == nil {
				firstErr = err
				close(stopChan)
			}

			if errCount == len(inputs) {
				if firstErr != nil {
					return fmt.Errorf("multiplexer error: %w", firstErr)
				}

				return nil
			}

		case <-ctx.Done():
			close(stopChan)

			for i := errCount; i < len(inputs); i++ {
				<-errChan
			}

			return fmt.Errorf("multiplexer cancelled: %w", ctx.Err())
		}
	}

	return nil
}
