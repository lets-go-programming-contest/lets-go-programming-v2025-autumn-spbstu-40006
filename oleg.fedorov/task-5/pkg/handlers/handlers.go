package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/dizey5k/task-5/pkg/utils"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case data, isChannelOpen := <-input:
			if !isChannelOpen {
				return nil
			}

			if strings.Contains(data, "no decorator") {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(data, "decorated: ") {
				data = "decorated: " + data
			}

			if err := utils.SendStringToOutput(ctx, output, data); err != nil {
				return err
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	var counter uint64

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context canceled: %w", ctx.Err())
		case data, isChannelOpen := <-input:
			if !isChannelOpen {
				return nil
			}

			index := atomic.AddUint64(&counter, 1) % uint64(len(outputs))

			if err := utils.SendStringToOutput(ctx, outputs[index], data); err != nil {
				return err
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	results := make(chan utils.ChannelResult, len(inputs))

	utils.StartChannelReaders(ctx, inputs, results, utils.ReadStringFromChannel)

	processor := func(data string) bool {
		return !strings.Contains(data, "no multiplexer")
	}

	return utils.ProcessChannelResults(ctx, results, output, len(inputs), processor)
}
