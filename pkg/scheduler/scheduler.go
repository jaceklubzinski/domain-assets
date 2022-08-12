package scheduler

import "time"

type Scheduler struct {
	initialDelay time.Duration
	delay        time.Duration
	ticker       time.Ticker
}

func (s *Scheduler) Executor(task func()) {
	firstRunChan := make(chan bool, 2)
	time.AfterFunc(s.initialDelay, func() {
		firstRunChan <- true
	})
	go func() {
		defer s.ticker.Stop()
		for {
			select {
			case <-s.ticker.C:
				go task()
			case <-firstRunChan:
				go task()
			}
		}
	}()
}

func NewScheduler(initialDelay, delay time.Duration) *Scheduler {
	return &Scheduler{
		initialDelay: initialDelay,
		delay:        delay,
		ticker:       *time.NewTicker(delay),
	}
}
