package conveyer

import "context"

type DecoratorFunc func(ctx context.Context, input chan string, output chan string) error

type MultiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error

type SeparatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type Conveyer interface {
	RegisterDecorator(fn DecoratorFunc, input string, output string)
	RegisterMultiplexer(fn MultiplexerFunc, inputs []string, output string)
	RegisterSeparator(fn SeparatorFunc, input string, outputs []string)

	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

// New создает новый конвейер с указанным размером буфера каналов
//
//nolint:ireturn
func New(size int) Conveyer {
	return newConveyerImpl(size)
}
