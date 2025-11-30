package utils

import (
	"context"
	"fmt"
)

type ChannelResult struct {
	Data string
	Okay bool
}

func StartChannelReaders(
	ctx context.Context,
	inputs []chan string,
	results chan<- ChannelResult,
	readerFunc func(context.Context, <-chan string, chan<- ChannelResult),
) {
	for _, inputChannel := range inputs {
		go readerFunc(ctx, inputChannel, results)
	}
}

func ReadStringFromChannel(ctx context.Context, input <-chan string, results chan<- ChannelResult) {
	for {
		select {
		case <-ctx.Done():
			return
		case data, isChannelOpen := <-input:
			SendChannelResult(ctx, results, ChannelResult{data, isChannelOpen})
			if !isChannelOpen {
				return
			}
		}
	}
}

func SendChannelResult(ctx context.Context, results chan<- ChannelResult, result ChannelResult) {
	select {
	case <-ctx.Done():
		return
	case results <- result:
	}
}

func SendStringToOutput(ctx context.Context, output chan<- string, data string) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled: %w", ctx.Err())
	case output <- data:
		return nil
	}
}

func ProcessChannelResults(
	ctx context.Context,
	results <-chan ChannelResult,
	output chan<- string,
	totalInputs int,
	processor func(string) bool,
) error {
	activeInputs := totalInputs

	for activeInputs > 0 {
		_, done, err := ReceiveAndProcessResult(ctx, results, output, processor)
		if err != nil {
			return err
		}
		if done {
			activeInputs--
		}
	}

	return nil
}

func ReceiveAndProcessResult(
	ctx context.Context,
	results <-chan ChannelResult,
	output chan<- string,
	shouldProcess func(string) bool,
) (ChannelResult, bool, error) {
	select {
	case <-ctx.Done():
		return ChannelResult{}, false, fmt.Errorf("context canceled: %w", ctx.Err())
	case result := <-results:
		if !result.Okay {
			return result, true, nil
		}

		if shouldProcess != nil && !shouldProcess(result.Data) {
			return result, false, nil
		}

		if err := SendStringToOutput(ctx, output, result.Data); err != nil {
			return result, false, err
		}

		return result, false, nil
	}
}
