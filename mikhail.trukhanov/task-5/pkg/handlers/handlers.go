package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input, output chan string) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(val, "no decorator") {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated: " + val
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- val:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	defer func() {
		for _, ch := range outputs {
			close(ch)
		}
	}()

	if len(outputs) == 0 {
		return nil
	}

	idx := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case outputs[idx%len(outputs)] <- val:
				idx++
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	var wg sync.WaitGroup

	for _, ch := range inputs {
		wg.Add(1)
		channel := ch

		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-channel:
					if !ok {
						return
					}

					if strings.Contains(val, "no multiplexer") {
						continue
					}

					select {
					case <-ctx.Done():
						return
					case output <- val:
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
