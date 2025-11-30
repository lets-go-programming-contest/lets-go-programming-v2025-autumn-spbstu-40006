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
		case v, ok := <-input:
			if !ok {
				close(output)
				return nil
			}
			if strings.Contains(v, "no decorator") {
				return errors.New("can't be decorated")
			}
			if !strings.HasPrefix(v, "decorated: ") {
				v = "decorated: " + v
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case output <- v:
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

	i := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case v, ok := <-input:
			if !ok {
				return nil
			}
			out := outputs[i%len(outputs)]
			i++
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- v:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		close(output)
		return nil
	}

	var wg sync.WaitGroup
	done := make(chan struct{})

	for _, in := range inputs {
		wg.Add(1)
		go func(ch chan string) {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if !ok {
						return
					}
					if strings.Contains(v, "no multiplexer") {
						continue
					}
					select {
					case <-done:
						return
					case <-ctx.Done():
						return
					case output <- v:
					}
				}
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(done)
		close(output)
	}()

	select {
	case <-ctx.Done():
		close(done)
		return ctx.Err()
	case <-done:
		return nil
	}
}
