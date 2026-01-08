func (c *conveyerImpl) Run(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()

		return ErrAlreadyRunning
	}

	c.running = true
	workersCopy := append([]func(context.Context) error(nil), c.workers...)
	c.mu.Unlock()

	runCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	var waitGroup sync.WaitGroup

	waitGroup.Add(len(workersCopy))

	errorChan := make(chan error, 1)

	for _, worker := range workersCopy {
		currentWorker := worker

		workerFunc := func() {
			defer waitGroup.Done()

			err := currentWorker(runCtx)
			if err != nil {
				select {
				case errorChan <- err:
				default:
				}
			}
		}

		go workerFunc()
	}

	var runErr error
	select {
	case <-ctx.Done():
	case runErr = <-errorChan:
	}

	cancelFunc()
	waitGroup.Wait()

	c.mu.Lock()
	channelsCopy := make([]chan string, 0, len(c.channels))

	for _, channel := range c.channels {
		channelsCopy = append(channelsCopy, channel)
	}

	c.running = false
	c.mu.Unlock()

	for _, channel := range channelsCopy {
		safeClose(channel)
	}

	return runErr
}

