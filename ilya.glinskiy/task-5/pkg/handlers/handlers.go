package handlers

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	const prefix = "decorated: "
	const errorSubstring = "no decorator"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, errorSubstring) {
				return errors.New("can't be decorated")
			}

			var result string
			if strings.HasPrefix(data, prefix) {
				result = data
			} else {
				result = prefix + data
			}

			select {
			case output <- result:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		return nil
	}

	var counter uint64

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case data, ok := <-input:
			if !ok {
				return nil
			}

			index := int(atomic.AddUint64(&counter, 1)-1) % len(outputs)

			select {
			case outputs[index] <- data:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	const filterSubstring = "no multiplexer"

	errChan := make(chan error, 1)
	done := make(chan struct{})

	for i, input := range inputs {
		go func(ch chan string, idx int) {
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					select {
					case errChan <- ctx.Err():
					default:
					}
					return
				case data, ok := <-ch:
					if !ok {
						return
					}

					if strings.Contains(data, filterSubstring) {
						continue
					}

					select {
					case output <- data:
					case <-ctx.Done():
						select {
						case errChan <- ctx.Err():
						default:
						}
						return
					}
				}
			}
		}(input, i)
	}

	select {
	case err := <-errChan:
		close(done)
		return err
	case <-ctx.Done():
		close(done)
		return ctx.Err()
	}
}
