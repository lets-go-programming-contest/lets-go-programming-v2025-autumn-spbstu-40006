package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("cannot_be_decorated")

func PrefixDecoratorFunc(
	applicationContext context.Context,
	inputChannel chan string,
	outputChannel chan string,
) error {
	defer close(outputChannel)

	for {
		select {
		case <-applicationContext.Done():
			return fmt.Errorf("prefix_decorator: %w", applicationContext.Err())

		case message, opened := <-inputChannel:
			if !opened {
				return nil
			}

			if strings.Contains(message, "no decorator") {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(message, "decorated: ") {
				message = "decorated: " + message
			}

			select {
			case outputChannel <- message:
			case <-applicationContext.Done():
				return fmt.Errorf("prefix_decorator: %w", applicationContext.Err())
			}
		}
	}
}

func SeparatorFunc(
	applicationContext context.Context,
	inputChannel chan string,
	outputChannelSlice []chan string,
) error {
	defer func() {
		for _, out := range outputChannelSlice {
			close(out)
		}
	}()

	currentPosition := 0

	for {
		select {
		case <-applicationContext.Done():
			return fmt.Errorf("separator: %w", applicationContext.Err())

		case message, opened := <-inputChannel:
			if !opened {
				return nil
			}

			if len(outputChannelSlice) == 0 {
				continue
			}

			targetIndex := currentPosition % len(outputChannelSlice)
			currentPosition = currentPosition + 1

			targetChannel := outputChannelSlice[targetIndex]

			select {
			case targetChannel <- message:
			case <-applicationContext.Done():
				return fmt.Errorf("separator: %w", applicationContext.Err())
			}
		}
	}
}

func MultiplexerFunc(
	applicationContext context.Context,
	inputChannelSlice []chan string,
	outputChannel chan string,
) error {
	defer close(outputChannel)

	if len(inputChannelSlice) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputChannelSlice))

	for _, source := range inputChannelSlice {
		sourceChannel := source

		go func() {
			defer waitGroup.Done()

			for {
				select {
				case <-applicationContext.Done():
					return

				case message, opened := <-sourceChannel:
					if !opened {
						return
					}

					if strings.Contains(message, "no multiplexer") {
						continue
					}

					select {
					case outputChannel <- message:
					case <-applicationContext.Done():
						return
					}
				}
			}
		}()
	}

	waitGroup.Wait()

	return nil
}
