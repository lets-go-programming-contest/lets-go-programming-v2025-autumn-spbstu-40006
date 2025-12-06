package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrNoDecorator  = errors.New("can't be decorated")
	ErrEmptyOutputs = errors.New("empty outputs")
)

const (
	noDecorator     = "no decorator"
	noMultiplexer   = "no multiplexer"
	decoratedPrefix = "decorated: "
)

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case str, status := <-input:
			if !status {
				return nil
			}

			if strings.Contains(str, noDecorator) {
				return ErrNoDecorator
			}

			if !strings.HasPrefix(str, decoratedPrefix) {
				str = decoratedPrefix + str
			}

			select {
			case output <- str:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var waitGroup sync.WaitGroup

	waitGroup.Add(len(inputs))

	for _, channelInput := range inputs {
		go func(channel chan string) {
			defer waitGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case str, status := <-channel:
					if !status {
						return
					}

					if strings.Contains(str, noMultiplexer) {
						continue
					}

					select {
					case output <- str:
					case <-ctx.Done():
						return
					}
				}
			}
		}(channelInput)
	}

	waitGroup.Wait()

	return nil
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return ErrEmptyOutputs
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case str, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case outputs[index] <- str:
			case <-ctx.Done():
				return nil
			}

			index = (index + 1) % len(outputs)
		}
	}
}
