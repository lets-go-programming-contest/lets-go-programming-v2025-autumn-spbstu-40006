package handlers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrCannotBeDecorated = errors.New("can't be decorated")

const (
	DecoratorPrefix = "decorated: "
	NoDecoratorMsg  = "no decorator"
	NoMultiplexMsg  = "no multiplexer"
)

func PrefixDecoratorFunc(ctx context.Context, src chan string, dst chan string) error {
	defer close(dst)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("prefix decorator: %w", ctx.Err())
		case message, ok := <-src:
			if !ok {
				return nil
			}

			if strings.Contains(message, NoDecoratorMsg) {
				return ErrCannotBeDecorated
			}

			if !strings.HasPrefix(message, DecoratorPrefix) {
				message = DecoratorPrefix + message
			}

			select {
			case dst <- message:
			case <-ctx.Done():
				return fmt.Errorf("prefix decorator: %w", ctx.Err())
			}
		}
	}
}

func SeparatorFunc(ctx context.Context, src chan string, dstList []chan string) error {
	defer func() {
		for _, dst := range dstList {
			close(dst)
		}
	}()

	position := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("separator: %w", ctx.Err())
		case message, ok := <-src:
			if !ok {
				return nil
			}

			if len(dstList) == 0 {
				continue
			}

			index := position % len(dstList)
			position++

			select {
			case dstList[index] <- message:
			case <-ctx.Done():
				return fmt.Errorf("separator: %w", ctx.Err())
			}
		}
	}
}

func MultiplexerFunc(ctx context.Context, srcList []chan string, dst chan string) error {
	defer close(dst)

	if len(srcList) == 0 {
		return nil
	}

	var sourceGroup sync.WaitGroup

	sourceGroup.Add(len(srcList))

	for _, source := range srcList {
		currentSource := source

		go func() {
			defer sourceGroup.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case message, ok := <-currentSource:
					if !ok {
						return
					}

					if strings.Contains(message, NoMultiplexMsg) {
						continue
					}

					select {
					case dst <- message:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}

	sourceGroup.Wait()

	return nil
}
