package pool

import (
	"errors"
	"time"

	"github.com/zk101/nixie/lib/couchbase"
)

// workerData is the workers local data
type workerData struct {
	bucket   *couchbase.Client
	lastWork time.Time
}

// worker is the actual go routine that does all the work
func (c *Client) worker(task taskRequest, initialStart bool) {
	c.wg.Add(1)

	w := &workerData{
		lastWork: time.Now().Add(time.Second * time.Duration(c.config.Expiry)),
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

	w.bucket = couchbase.NewClient(c.bucket)
	if err := w.bucket.Connect(); err != nil {
		loopControl = false
	}
	defer w.bucket.Close()

	if initialStart == false {
		w.process(&task)
	}

	for loopControl {
		select {
		case <-timeout:
			loopControl = false
		case task := <-c.work:
			if err := w.process(&task); err != nil {
				loopControl = false
			}
			w.lastWork = time.Now().Add(time.Second * time.Duration(c.config.Expiry))
		}
	}

	<-c.register
	c.wg.Done()
}

// process does the actual task processing
func (w *workerData) process(task *taskRequest) error {
	var err error
	reply := taskReply{}

	switch task.action {
	case actionGet:
		err = w.actionGet(task, &reply)

	case actionInsert:
		err = w.actionInsert(task, &reply)

	case actionRemove:
		err = w.actionRemove(task, &reply)

	case actionReplace:
		err = w.actionReplace(task, &reply)

	case actionTouch:
		err = w.actionTouch(task, &reply)

	case actionUpsert:
		err = w.actionUpsert(task, &reply)

	default:
		reply.err = errors.New("unsupported request")
		err = errors.New("unsupported request")
	}

	task.reply <- reply
	return err
}

// actionGet does a Get task
func (w *workerData) actionGet(task *taskRequest, reply *taskReply) error {
	if task.key == "" {
		reply.err = errors.New("missing key for request")
		return errors.New("missing key for request")
	}
	if task.data == nil {
		reply.err = errors.New("missing data for request")
		return errors.New("missing data for request")
	}

	cas, err := w.bucket.Get(task.key, task.data)
	if err != nil {
		reply.err = err
		return err
	}
	reply.cas = cas

	return nil
}

// actionInsert does a Insert task
func (w *workerData) actionInsert(task *taskRequest, reply *taskReply) error {
	if task.key == "" {
		reply.err = errors.New("missing key for request")
		return errors.New("missing key for request")
	}
	if task.data == nil {
		reply.err = errors.New("missing data for request")
		return errors.New("missing data for request")
	}

	cas, err := w.bucket.Insert(task.key, task.data, task.expiry)
	if err != nil {
		reply.err = err
		return err
	}
	reply.cas = cas

	return nil
}

// actionRemove does a Remove task
func (w *workerData) actionRemove(task *taskRequest, reply *taskReply) error {
	if task.key == "" {
		reply.err = errors.New("missing key for request")
		return errors.New("missing key for request")
	}

	cas, err := w.bucket.Remove(task.key, task.cas)
	if err != nil {
		reply.err = err
		return err
	}
	reply.cas = cas

	return nil
}

// actionReplace does a Replace task
func (w *workerData) actionReplace(task *taskRequest, reply *taskReply) error {
	if task.key == "" {
		reply.err = errors.New("missing key for request")
		return errors.New("missing key for request")
	}
	if task.data == nil {
		reply.err = errors.New("missing data for request")
		return errors.New("missing data for request")
	}

	cas, err := w.bucket.Replace(task.key, task.data, task.cas, task.expiry)
	if err != nil {
		reply.err = err
		return err
	}
	reply.cas = cas

	return nil
}

// actionTouch does a Touch task
func (w *workerData) actionTouch(task *taskRequest, reply *taskReply) error {
	if task.key == "" {
		reply.err = errors.New("missing key for request")
		return errors.New("missing key for request")
	}

	cas, err := w.bucket.Touch(task.key, task.cas, task.expiry)
	if err != nil {
		reply.err = err
		return err
	}
	reply.cas = cas

	return nil
}

// actionUpsert does a Upsert task
func (w *workerData) actionUpsert(task *taskRequest, reply *taskReply) error {
	if task.key == "" {
		reply.err = errors.New("missing key for request")
		return errors.New("missing key for request")
	}
	if task.data == nil {
		reply.err = errors.New("missing data for request")
		return errors.New("missing data for request")
	}

	cas, err := w.bucket.Upsert(task.key, task.data, task.expiry)
	if err != nil {
		reply.err = err
		return err
	}
	reply.cas = cas

	return nil
}

// EOF
