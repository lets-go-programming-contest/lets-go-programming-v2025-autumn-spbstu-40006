package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(val, "no decorator") {
				return errors.New("can't be decorated")
			}

			if !strings.HasPrefix(val, "decorated: ") {
				val = "decorated: " + val
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- val:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	index := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case val, ok := <-input:
			if !ok {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case outputs[index%len(outputs)] <- val:
				index++
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(inputs))

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
					case <-ctx.Done():
						return
					case output <- val:
					}
				}
			}
		}(ch)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	case err := <-errCh:
		return err
	}
}
