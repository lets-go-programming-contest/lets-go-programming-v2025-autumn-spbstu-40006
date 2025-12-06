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
			case output <- val:
			case <-ctx.Done():
				return nil
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

	counter := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case val, ok := <-input:
			if !ok {
				return nil
			}
			idx := counter % len(outputs)
			counter++

			select {
			case outputs[idx] <- val:
			case <-ctx.Done():
				return nil
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	defer close(output)

	var wg sync.WaitGroup
	for _, ch := range inputs {
		wg.Add(1)
		go func(c chan string) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-c:
					if !ok {
						return
					}
					if strings.Contains(val, "no multiplexer") {
						continue
					}
					select {
					case output <- val:
					case <-ctx.Done():
						return
					}
				}
			}
		}(ch)
	}

	wg.Wait()
	return nil
}
