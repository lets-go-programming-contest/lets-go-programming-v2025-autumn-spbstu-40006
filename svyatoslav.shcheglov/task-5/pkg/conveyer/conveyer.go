package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChannelNotFound = errors.New("channel_not_found")

type Conveyer struct {
	channelMap map[string]chan string
	taskSlice  []func(context.Context) error
	mutex      sync.Mutex
	bufferSize int
}

func NewConveyer(size int) *Conveyer {
	newConveyer := &Conveyer{
		channelMap: make(map[string]chan string),
		taskSlice:  make([]func(context.Context) error, 0),
		mutex:      sync.Mutex{},
		bufferSize: size,
	}

	return newConveyer
}

func New(size int) *Conveyer {
	return NewConveyer(size)
}

func (conveyer *Conveyer) obtainChannel(identifier string) chan string {
	conveyer.mutex.Lock()

	channelObject, exists := conveyer.channelMap[identifier]
	if !exists {
		createdChannel := make(chan string, conveyer.bufferSize)
		channelObject = createdChannel
		conveyer.channelMap[identifier] = createdChannel
	}

	conveyer.mutex.Unlock()

	return channelObject
}

func (conveyer *Conveyer) fetchChannel(identifier string) (chan string, error) {
	conveyer.mutex.Lock()

	channelObject, exists := conveyer.channelMap[identifier]

	conveyer.mutex.Unlock()

	if !exists {
		return nil, ErrChannelNotFound
	}

	return channelObject, nil
}

func (conveyer *Conveyer) RegisterDecorator(
	processor func(context.Context, chan string, chan string) error,
	sourceID string,
	targetID string,
) {
	newTask := func(ctx context.Context) error {
		sourceChannel := conveyer.obtainChannel(sourceID)
		targetChannel := conveyer.obtainChannel(targetID)

		err := processor(ctx, sourceChannel, targetChannel)

		return err
	}

	conveyer.taskSlice = append(conveyer.taskSlice, newTask)
}

func (conveyer *Conveyer) RegisterSeparator(
	processor func(context.Context, chan string, chan string, chan string) error,
	inputID string,
	trueID string,
	falseID string,
) {
	newTask := func(ctx context.Context) error {
		inputChannel := conveyer.obtainChannel(inputID)
		trueChannel := conveyer.obtainChannel(trueID)
		falseChannel := conveyer.obtainChannel(falseID)

		err := processor(ctx, inputChannel, trueChannel, falseChannel)

		return err
	}

	conveyer.taskSlice = append(conveyer.taskSlice, newTask)
}

func (conveyer *Conveyer) RegisterMultiplexer(
	processor func(context.Context, chan string, chan string, chan string) error,
	leftID string,
	rightID string,
	outputID string,
) {
	newTask := func(ctx context.Context) error {
		leftChannel := conveyer.obtainChannel(leftID)
		rightChannel := conveyer.obtainChannel(rightID)
		outputChannel := conveyer.obtainChannel(outputID)

		err := processor(ctx, leftChannel, rightChannel, outputChannel)

		return err
	}

	conveyer.taskSlice = append(conveyer.taskSlice, newTask)
}

func (conveyer *Conveyer) Send(identifier string, value string) error {
	channelObject := conveyer.obtainChannel(identifier)

	channelObject <- value

	return nil
}

func (conveyer *Conveyer) Recv(identifier string) (string, error) {
	channelObject, fetchError := conveyer.fetchChannel(identifier)
	if fetchError != nil {
		return "", fetchError
	}

	value, open := <-channelObject
	if !open {
		return "", nil
	}

	return value, nil
}

func (conveyer *Conveyer) Run(ctx context.Context) error {
	var waitGroup sync.WaitGroup

	errorChannel := make(chan error, len(conveyer.taskSlice))

	for _, task := range conveyer.taskSlice {
		taskFunction := task

		waitGroup.Add(1)

		go func() {
			defer waitGroup.Done()

			err := taskFunction(ctx)
			if err != nil {
				errorChannel <- err
			}
		}()
	}

	waitGroup.Wait()

	close(errorChannel)

	for err := range errorChannel {
		return err
	}

	return nil
}
