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
		case message, ok := <-inputChannel:
			if !ok {
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
		case message, ok := <-inputChannel:
			if !ok {
				return nil
			}

			if len(outputChannels) == 0 {
				continue
			}

			index := positionCounter % len(outputChannels)
			positionCounter++

			select {
			case outputChannels[index] <- message:
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

	for i := range inputChannels {
		currentChannel := inputChannels[i]

		go func(channelToRead chan string) {
			defer sourceGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case message, ok := <-channelToRead:
					if !ok {
						return
					}

					if strings.Contains(message, "no multiplexer") {
						continue
					}

					select {
					case outputChannel <- message:
					case <-ctx.Done():
						return
					}
				}
			}
		}(currentChannel)
	}

	sourceGroup.Wait()

	return nil
}
