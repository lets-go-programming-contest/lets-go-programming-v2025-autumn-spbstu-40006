package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

const (
	Prefix               = "decorated: "
	NoDecoratorMessage   = "no decorator"
	NoMultiplexerMessage = "no multiplexer"
	Undefined            = "undefined"
)

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
			case output <- data:
			case <-ctx.Done():
				return nil
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

	if len(outputs) == 0 {
		return nil
	}

	outputIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			currentOutput := outputs[outputIndex%len(outputs)]
			outputIndex++

			select {
			case currentOutput <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	inputWaitGroup := sync.WaitGroup{}

	for _, inputChannel := range inputs {
		inputWaitGroup.Add(1)

		go func(currentInput chan string) {
			defer inputWaitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-currentInput:
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
		}(inputChannel)
	}

	inputWaitGroup.Wait()

	return nil
}
