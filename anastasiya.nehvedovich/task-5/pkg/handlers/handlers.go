package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

const (
	noDecoratorMessage   = "no decorator"
	decoratorPrefix      = "decorated: "
	noMultiplexerMessage = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil

		case value, isOpen := <-input:
			if !isOpen {
				return nil
			}

			if strings.Contains(value, noDecoratorMessage) {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(value, decoratorPrefix) {
				value = decoratorPrefix + value
			}

			select {
			case output <- value:
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

	currentIndex := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case value, isOpen := <-input:
			if !isOpen {
				return nil
			}

			select {
			case outputs[currentIndex] <- value:
				currentIndex = (currentIndex + 1) % len(outputs)
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

	var waitGroup sync.WaitGroup

	for _, inputChannel := range inputs {
		channel := inputChannel

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case value, isOpen := <-channel:
					if !isOpen {
						return
					}

					if strings.Contains(value, noMultiplexerMessage) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- value:
					}
				}
			}
		}()
	}

	waitGroup.Wait()

	return nil
}
