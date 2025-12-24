package handlers

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var errCantBeDecorated = errors.New("can’t be decorated")

const (
	prefix              = "decorated: "
	noDecoratorSubstr   = "no decorator"
	noMultiplexerSubstr = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			if strings.Contains(data, noDecoratorSubstr) {
				return errCantBeDecorated
			}

			if !strings.HasPrefix(data, prefix) {
				data = prefix + data
			}

			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
	if len(outputs) == 0 {
		for {
			select {
			case <-ctx.Done():
				return nil
			case _, ok := <-input:
				if !ok {
					return nil
				}
			}
		}
	}

	i := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-input:
			if !ok {
				return nil
			}

			out := outputs[i%len(outputs)]
			i++

			select {
			case <-ctx.Done():
				return nil
			case out <- data:
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
	if len(inputs) == 0 {
		return nil
	}

	merged := make(chan string)

	var wg sync.WaitGroup
	wg.Add(len(inputs))

	for _, in := range inputs {
		inputChan := in
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case data, ok := <-inputChan:
					if !ok {
						return
					}
					if strings.Contains(data, noMultiplexerSubstr) {
						continue
					}
					select {
					case <-ctx.Done():
						return
					case merged <- data:
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(merged)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data, ok := <-merged:
			if !ok {
				return nil
			}
			select {
			case <-ctx.Done():
				return nil
			case output <- data:
			}
		}
	}
}
