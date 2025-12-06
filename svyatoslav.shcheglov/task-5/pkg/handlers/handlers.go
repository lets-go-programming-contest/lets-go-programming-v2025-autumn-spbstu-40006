package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, inputChannel chan string, outputChannel chan string) error {
	defer close(outputChannel)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("prefix decorator: %w", ctx.Err())
		case messageFromInput, ok := <-inputChannel:
			if !ok {
				return nil
			}

			if strings.Contains(messageFromInput, "no decorator") {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(messageFromInput, "decorated: ") {
				messageFromInput = "decorated: " + messageFromInput
			}

			select {
			case outputChannel <- messageFromInput:
			case <-ctx.Done():
				return fmt.Errorf("prefix decorator: %w", ctx.Err())
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, inputChannel chan string, outputChannels []chan string) error {
	defer func() {
		for _, outputChannel := range outputChannels {
			close(outputChannel)
		}
	}()

	positionCounter := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("separator: %w", ctx.Err())
		case messageFromInput, ok := <-inputChannel:
			if !ok {
				return nil
			}

			if len(outputChannels) == 0 {
				continue
			}

			index := positionCounter % len(outputChannels)
			positionCounter++

			select {
			case outputChannels[index] <- messageFromInput:
			case <-ctx.Done():
				return fmt.Errorf("separator: %w", ctx.Err())
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChannels []chan string, outputChannel chan string) error {
	defer close(outputChannel)

	if len(inputChannels) == 0 {
		return nil
	}

	var sourceGroup sync.WaitGroup
	sourceGroup.Add(len(inputChannels))

	for _, currentInputChannel := range inputChannels {
		currentInputChannel := currentInputChannel

		go func() {
			defer sourceGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case messageFromInput, ok := <-currentInputChannel:
					if !ok {
						return
					}

					if strings.Contains(messageFromInput, "no multiplexer") {
						continue
					}

					select {
					case outputChannel <- messageFromInput:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	sourceGroup.Wait()

	return nil
}
