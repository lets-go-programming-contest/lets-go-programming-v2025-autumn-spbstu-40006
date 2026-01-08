package handlers

import (
	"context"
	"errors"
	"strings"
)

var ErrCantDecorate = errors.New("can't be decorated")

const (
	decorPrefix = "decorated: "
	blockDecor  = "no decorator"
	blockMux    = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil

		case str, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(str, blockDecor) {
				return ErrCantDecorate
			}

			if !strings.HasPrefix(str, decorPrefix) {
				str = decorPrefix + str
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- str:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, out := range outputs {
			close(out)
		}
	}()

	outCnt := len(outputs)
	pos := 0

	for {
		select {
		case <-ctx.Done():
			return nil

		case data, ok := <-input:
			if !ok {
				return nil
			}

			if outCnt == 0 {
				continue
			}

			dstData := outputs[pos%outCnt]
			pos++

			select {
			case <-ctx.Done():
				return nil
			case dstData <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	if len(inputs) == 0 {
		return nil
	}

	for _, src := range inputs {
		channel := src

		go func(chLocal chan string) {
			for {
				select {
				case <-ctx.Done():
					return

				case data, ok := <-chLocal:
					if !ok {
						return
					}

					if strings.Contains(data, blockMux) {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- data:
					}
				}
			}
		}(channel)
	}

	<-ctx.Done()

	return nil
}
