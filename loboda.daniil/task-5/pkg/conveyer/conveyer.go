package conveyer

import "context"

type DecoratorFunc func(ctx context.Context, inputChan chan string, outputChan chan string) error

type MultiplexerFunc func(ctx context.Context, inputChans []chan string, outputChan chan string) error

type SeparatorFunc func(ctx context.Context, inputChan chan string, outputChans []chan string) error

type conveyer interface {
	RegisterDecorator(handler DecoratorFunc, inputName string, outputName string)
	RegisterMultiplexer(handler MultiplexerFunc, inputNames []string, outputName string)
	RegisterSeparator(handler SeparatorFunc, inputName string, outputNames []string)

	Run(ctx context.Context) error
	Send(inputName string, data string) error
	Recv(outputName string) (string, error)
}

func New(size int) *conveyerImpl {
	return newConveyerImpl(size)
}

var _ conveyer = (*conveyerImpl)(nil)
