package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("cant_be_decorated")

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

		case receivedMessage, channelOpen := <-inputChannel:
			if !channelOpen {
				return nil
			}

			if strings.Contains(receivedMessage, "no decorator") {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(receivedMessage, "decorated: ") {
				receivedMessage = "decorated: " + receivedMessage
			}

			select {
			case outputChannel <- receivedMessage:
			case <-applicationContext.Done():
				return fmt.Errorf("prefix_decorator: %w", applicationContext.Err())
			}
		}
	}
}

func SeparatorFunc(
	applicationContext context.Context,
	inputChannel chan string,
	outputChannelList []chan string,
) error {
	defer func() {
		for _, outputChannel := range outputChannelList {
			close(outputChannel)
		}
	}()

	currentPositionIndex := 0

	for {
		select {
		case <-applicationContext.Done():
			return fmt.Errorf("separator: %w", applicationContext.Err())

		case receivedMessage, channelOpen := <-inputChannel:
			if !channelOpen {
				return nil
			}

			if len(outputChannelList) == 0 {
				continue
			}

			targetIndex := currentPositionIndex % len(outputChannelList)
			currentPositionIndex++

			select {
			case outputChannelList[targetIndex] <- receivedMessage:
			case <-applicationContext.Done():
				return fmt.Errorf("separator: %w", applicationContext.Err())
			}
		}
	}
}

func MultiplexerFunc(
	applicationContext context.Context,
	inputChannelList []chan string,
	outputChannel chan string,
) error {
	defer close(outputChannel)

	if len(inputChannelList) == 0 {
		return nil
	}

	var readerGroup sync.WaitGroup
	readerGroup.Add(len(inputChannelList))

	for _, sourceChannel := range inputChannelList {
		currentSourceChannel := sourceChannel

		go func(channelToRead chan string) {
			defer readerGroup.Done()

			for {
				select {
				case <-applicationContext.Done():
					return

				case receivedMessage, channelOpen := <-channelToRead:
					if !channelOpen {
						return
					}

					if strings.Contains(receivedMessage, "no multiplexer") {
						continue
					}

					select {
					case outputChannel <- receivedMessage:
					case <-applicationContext.Done():
						return
					}
				}
			}
		}(currentSourceChannel)
	}

	readerGroup.Wait()

	return nil
}
