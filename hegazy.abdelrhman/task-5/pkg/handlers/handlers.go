package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrPrefixDecoratorFuncCantBeDecorated = errors.New("can't be decorated")

func safelyCloseChannel(channel chan string) {
	defer func() {
		if r := recover(); r != nil {
			_ = r
		}
	}()

	close(channel)
}

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error {
	defer safelyCloseChannel(output)

	prefix := "decorated: "

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrPrefixDecoratorFuncCantBeDecorated
			}

			processedData := data
			if !strings.HasPrefix(data, prefix) {
				processedData = prefix + data
			}

			select {
			case output <- processedData:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func readInputToTransfer(
	ctx context.Context,
	input chan string,
	transfer chan string,
	waitGroup *sync.WaitGroup,
) {
	defer waitGroup.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-input:
			if !ok {
				return
			}

			select {
			case transfer <- data:
			case <-ctx.Done():
				return
			}
		}
	}
}

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error {
	var waitGroup sync.WaitGroup

	transfer := make(chan string)

	defer safelyCloseChannel(output)

	readFn := readInputToTransfer

	for _, inputChan := range inputs {
		waitGroup.Add(1)

		localInput := inputChan
		go readFn(ctx, localInput, transfer, &waitGroup)
	}

	go func() {
		waitGroup.Wait()
		close(transfer)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-transfer:
			if !ok {
				return nil
			}

			if strings.Contains(data, "no multiplexer") {
				continue
			}

			select {
			case output <- data:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func SeparatorFunc(
	ctx context.Context,
	input chan string,
	outputs []chan string,
) error {
	defer func() {
		closedChannels := make(map[chan string]struct{})

		for _, output := range outputs {
			if _, exists := closedChannels[output]; exists {
				continue
			}

			closedChannels[output] = struct{}{}

			safelyCloseChannel(output)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	var counter int

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			targetIndex := counter % len(outputs)
			targetChannel := outputs[targetIndex]

			select {
			case targetChannel <- data:
				counter++
			case <-ctx.Done():
				return nil
			}
		}
	}
}
