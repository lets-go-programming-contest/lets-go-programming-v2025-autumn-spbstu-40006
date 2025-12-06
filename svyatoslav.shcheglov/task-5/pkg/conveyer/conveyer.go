package conveyer

import (
	"context"
	"errors"
	"sync"
)

var ErrChannelNotFound = errors.New("channel not found")

type Conveyer struct {
	channelMap map[string]chan string
	taskList   []func(context.Context) error
	mutex      sync.Mutex
	bufferSize int
}

func NewConveyer(bufferSize int) *Conveyer {
	return &Conveyer{
		channelMap: make(map[string]chan string),
		taskList:   make([]func(context.Context) error, 0),
		bufferSize: bufferSize,
	}
}

func New(bufferSize int) *Conveyer {
	return NewConveyer(bufferSize)
}

func (conveyer *Conveyer) obtainChannel(channelIdentifier string) chan string {
	conveyer.mutex.Lock()
	channelObject, exists := conveyer.channelMap[channelIdentifier]
	if !exists {
		channelObject = make(chan string, conveyer.bufferSize)
		conveyer.channelMap[channelIdentifier] = channelObject
	}
	conveyer.mutex.Unlock()
	return channelObject
}

func (conveyer *Conveyer) fetchChannel(channelIdentifier string) (chan string, error) {
	conveyer.mutex.Lock()
	channelObject, exists := conveyer.channelMap[channelIdentifier]
	conveyer.mutex.Unlock()
	if !exists {
		return nil, ErrChannelNotFound
	}
	return channelObject, nil
}

func (conveyer *Conveyer) RegisterDecorator(
	processor func(context.Context, chan string, chan string) error,
	sourceIdentifier string,
	destinationIdentifier string,
) {
	task := func(applicationContext context.Context) error {
		sourceChannel := conveyer.obtainChannel(sourceIdentifier)
		destinationChannel := conveyer.obtainChannel(destinationIdentifier)
		return processor(applicationContext, sourceChannel, destinationChannel)
	}
	conveyer.taskList = append(conveyer.taskList, task)
}

func (conveyer *Conveyer) RegisterSeparator(
	processor func(context.Context, chan string, chan string, chan string) error,
	inputIdentifier string,
	trueIdentifier string,
	falseIdentifier string,
) {
	task := func(applicationContext context.Context) error {
		inputChannel := conveyer.obtainChannel(inputIdentifier)
		trueChannel := conveyer.obtainChannel(trueIdentifier)
		falseChannel := conveyer.obtainChannel(falseIdentifier)
		return processor(applicationContext, inputChannel, trueChannel, falseChannel)
	}
	conveyer.taskList = append(conveyer.taskList, task)
}

func (conveyer *Conveyer) RegisterMultiplexer(
	processor func(context.Context, chan string, chan string, chan string) error,
	leftIdentifier string,
	rightIdentifier string,
	outputIdentifier string,
) {
	task := func(applicationContext context.Context) error {
		leftChannel := conveyer.obtainChannel(leftIdentifier)
		rightChannel := conveyer.obtainChannel(rightIdentifier)
		outputChannel := conveyer.obtainChannel(outputIdentifier)
		return processor(applicationContext, leftChannel, rightChannel, outputChannel)
	}
	conveyer.taskList = append(conveyer.taskList, task)
}

func (conveyer *Conveyer) Send(channelIdentifier string, value string) error {
	channelObject := conveyer.obtainChannel(channelIdentifier)
	channelObject <- value
	return nil
}

func (conveyer *Conveyer) Recv(channelIdentifier string) (string, error) {
	channelObject, errorValue := conveyer.fetchChannel(channelIdentifier)
	if errorValue != nil {
		return "", errorValue
	}
	value, isOpen := <-channelObject
	if !isOpen {
		return "", nil
	}
	return value, nil
}

func (conveyer *Conveyer) Run(applicationContext context.Context) error {
	var waitGroup sync.WaitGroup
	errorChannel := make(chan error, len(conveyer.taskList))

	for _, task := range conveyer.taskList {
		waitGroup.Add(1)
		go func(taskFunction func(context.Context) error) {
			defer waitGroup.Done()
			errorValue := taskFunction(applicationContext)
			if errorValue != nil {
				errorChannel <- errorValue
			}
		}(task)
	}

	waitGroup.Wait()
	close(errorChannel)

	for errorValue := range errorChannel {
		return errorValue
	}

	return nil
}
