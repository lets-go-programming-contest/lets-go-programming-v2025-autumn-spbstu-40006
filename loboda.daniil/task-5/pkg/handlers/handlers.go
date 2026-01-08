package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCantBeDecorated = errors.New("can't be decorated (can’t be decorated)")

const (
	prefix               = "decorated: "
	noDecoratorSubstring = "no decorator"
	noMultiplexerSubstr  = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, inputChan chan string, outputChan chan string) error {
	defer close(outputChan)

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, open := <-inputChan:
			if !open {
				return nil
			}

			if strings.Contains(data, noDecoratorSubstring) {
				return ErrCantBeDecorated
			}

			result := data
			if !strings.HasPrefix(data, prefix) {
				result = prefix + data
			}

			select {
			case <-ctx.Done():
				return nil

			case outputChan <- result:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, inputChan chan string, outputChans []chan string) error {
	defer closeAll(outputChans)

	if len(outputChans) == 0 {
		return drainInput(ctx, inputChan)
	}

	index := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, open := <-inputChan:
			if !open {
				return nil
			}

			currentOutput := outputChans[index%len(outputChans)]
			index++

			select {
			case <-ctx.Done():
				return nil

			case currentOutput <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputChans []chan string, outputChan chan string) error {
	defer close(outputChan)

	if len(inputChans) == 0 {
		return nil
	}

	mergedChan := make(chan string, len(inputChans)*2)

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(inputChans))

	for _, inputChan := range inputChans {
		currentInput := inputChan

		forwarder := func() {
			defer waitGroup.Done()
			forwardInput(ctx, currentInput, mergedChan)
		}

		go forwarder()
	}

	closer := func() {
		waitGroup.Wait()
		close(mergedChan)
	}

	go closer()

	return forwardMerged(ctx, mergedChan, outputChan)
}

func forwardInput(ctx context.Context, inputChan <-chan string, mergedChan chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return

		case data, open := <-inputChan:
			if !open {
				return
			}

			if strings.Contains(data, noMultiplexerSubstr) {
				continue
			}

			select {
			case <-ctx.Done():
				return

			case mergedChan <- data:
			}
		}
	}
}

func forwardMerged(ctx context.Context, mergedChan <-chan string, outputChan chan<- string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case data, open := <-mergedChan:
			if !open {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil

			case outputChan <- data:
			}
		}
	}
}

func drainInput(ctx context.Context, inputChan <-chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case _, open := <-inputChan:
			if !open {
				return nil
			}
		}
	}
}

func closeAll(chans []chan string) {
	for _, channel := range chans {
		close(channel)
	}
}
