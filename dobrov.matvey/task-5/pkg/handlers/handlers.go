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
	outputCount := len(outputs)

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			out := outputs[index]

			index = (index + 1) % outputCount

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

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, inputChan := range inputs {
		go func(inputChan chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case data, ok := <-inputChan:
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
		}(inputChan)
	}

	waitGroup.Wait()

	return nil
}
