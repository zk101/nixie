package gopool

import "time"

// workerData is the workers local data
type workerData struct {
	lastWork time.Time
}

// worker is the actual go routine that does all the work
func (c *Client) worker(task func()) {
	c.wg.Add(1)

	// defer recovery() incase a task() panics
	defer func() {
		if r := recover(); r != nil {
			<-c.register
			c.wg.Done()
			return
		}
	}()

	w := &workerData{
		lastWork: time.Now().Add(time.Second * time.Duration(c.config.WorkerExpiry)),
	}
	timeout := make(chan bool, 1)
	loopControl := true

	go func() {
		for {
			if time.Now().After(w.lastWork) {
				timeout <- true
				break
			}
			if c.run == false {
				timeout <- true
				break
			}
			time.Sleep(time.Second)
		}
	}()

	if task != nil {
		task()
	}

	for loopControl {
		select {
		case <-timeout:
			loopControl = false
		case task := <-c.work:
			task()
			w.lastWork = time.Now().Add(time.Second * time.Duration(c.config.WorkerExpiry))
		}
	}

	<-c.register
	c.wg.Done()
}

// EOF
