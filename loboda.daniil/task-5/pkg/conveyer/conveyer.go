package conveyer

import "context"

type decoratorFunc func(ctx context.Context, input chan string, output chan string) error
type multiplexerFunc func(ctx context.Context, inputs []chan string, output chan string) error
type separatorFunc func(ctx context.Context, input chan string, outputs []chan string) error

type conveyer interface {
	RegisterDecorator(fn decoratorFunc, input string, output string)
	RegisterMultiplexer(fn multiplexerFunc, inputs []string, output string)
	RegisterSeparator(fn separatorFunc, input string, outputs []string)

	Run(ctx context.Context) error
	Send(input string, data string) error
	Recv(output string) (string, error)
}

func New(size int) conveyer {
	return newImpl(size)
}
