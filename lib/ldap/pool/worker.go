package pool

import (
	"errors"
	"time"

	"github.com/zk101/nixie/lib/ldap"
)

// workerData is the workers local data
type workerData struct {
	ldapbind *ldap.Client
	ldapro   *ldap.Client
	ldaprw   *ldap.Client
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

	/*
		TODO: Remove ldapbind, the bind function should be(has) been updated to first disconnect, try a bind and leave the conn disconnected.actionAdd
		This will first the next task to connect, however, this is a small overhead for a long running worker, and simplifies the code
	*/

	if c.ldapbind != nil {
		w.ldapbind = ldap.NewClient(c.ldapbind, c.cacert)
		if err := w.ldapbind.Connect(); err != nil {
			loopControl = false
		}
	}

	if c.ldapro != nil {
		w.ldapro = ldap.NewClient(c.ldapro, c.cacert)
		if err := w.ldapro.Connect(); err != nil {
			loopControl = false
		}
	}

	if c.ldaprw != nil {
		w.ldaprw = ldap.NewClient(c.ldaprw, c.cacert)
		if err := w.ldaprw.Connect(); err != nil {
			loopControl = false
		}
	}

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

	if w.ldapbind != nil {
		w.ldapbind.Close()
	}

	if w.ldapro != nil {
		w.ldapro.Close()
	}

	if w.ldaprw != nil {
		w.ldaprw.Close()
	}

	<-c.register
	c.wg.Done()
}

// process does the actual task processing
func (w *workerData) process(task *taskRequest) error {
	var err error
	reply := taskReply{}

	switch task.action {
	case actionAdd:
		err = w.actionAdd(task, &reply)

	case actionBind:
		err = w.actionBind(task, &reply)

	case actionDelete:
		err = w.actionDelete(task, &reply)

	case actionMod:
		err = w.actionMod(task, &reply)

	case actionModDn:
		err = w.actionModDN(task, &reply)

	case actionSearch:
		err = w.actionSearch(task, &reply)

	default:
		reply.err = errors.New("unsupported request")
		err = errors.New("unsupported request")
	}

	task.reply <- reply
	return err
}

// actionAdd does an Add task
func (w *workerData) actionAdd(task *taskRequest, reply *taskReply) error {
	if task.add == nil {
		reply.err = errors.New("missing add request")
		return errors.New("missing add request")
	}
	if w.ldaprw == nil {
		reply.err = errors.New("add request requires ldaprw")
		return errors.New("add request requires ldaprw")
	}
	if err := w.ldaprw.Add(task.add); err != nil {
		reply.err = err
		return err
	}
	return nil
}

// actionBind does a Bind task
func (w *workerData) actionBind(task *taskRequest, reply *taskReply) error {
	if task.bind == nil {
		reply.err = errors.New("missing bind request")
		return errors.New("missing bind request")
	}
	if err := w.ldapbind.Bind(task.bind.user, task.bind.pass); err != nil {
		reply.err = errors.New("ldap bind error")
		return errors.New("ldap bind error")
	}
	return nil
}

// actionDelete does a Delete task
func (w *workerData) actionDelete(task *taskRequest, reply *taskReply) error {
	if task.delete == nil {
		reply.err = errors.New("missing delete request")
		return errors.New("missing delete request")
	}
	if w.ldaprw == nil {
		reply.err = errors.New("delete request requires ldaprw")
		return errors.New("delete request requires ldaprw")
	}
	if err := w.ldaprw.Delete(task.delete); err != nil {
		reply.err = err
		return err
	}
	return nil
}

// actionMod does an Mod task
func (w *workerData) actionMod(task *taskRequest, reply *taskReply) error {
	if task.mod == nil {
		reply.err = errors.New("missing mod request")
		return errors.New("missing mod request")
	}
	if w.ldaprw == nil {
		reply.err = errors.New("mod request requires ldaprw")
		return errors.New("mod request requires ldaprw")
	}
	if err := w.ldaprw.Modify(task.mod); err != nil {
		reply.err = err
		return err
	}
	return nil
}

// actionModDN does an ModDN task
func (w *workerData) actionModDN(task *taskRequest, reply *taskReply) error {
	if task.modDn == nil {
		reply.err = errors.New("missing moddn request")
		return errors.New("missing moddn request")
	}
	if w.ldaprw == nil {
		reply.err = errors.New("moddn request requires ldaprw")
		return errors.New("moddn request requires ldaprw")
	}
	if err := w.ldaprw.ModDn(task.modDn); err != nil {
		reply.err = err
		return err
	}
	return nil
}

// actionSearch does a Search task
func (w *workerData) actionSearch(task *taskRequest, reply *taskReply) error {
	if task.search == nil {
		reply.err = errors.New("missing search request")
		return errors.New("missing search request")
	}

	var ldapConn *ldap.Client
	if w.ldapro == nil {
		ldapConn = w.ldaprw
	} else {
		ldapConn = w.ldapro
	}

	reply.sr, reply.err = ldapConn.Search(task.search)
	if reply.err != nil {
		return reply.err
	}
	return nil
}

// EOF
