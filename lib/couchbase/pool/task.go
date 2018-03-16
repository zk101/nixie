package pool

import (
	"github.com/couchbase/gocb"
)

// action holds
const (
	actionGet = iota
	actionInsert
	actionRemove
	actionReplace
	actionTouch
	actionUpsert
)

// task holds an actual task details
type taskRequest struct {
	reply  chan taskReply
	action int
	key    string
	cas    gocb.Cas
	data   interface{}
	expiry uint32
}

// TaskReply holds a completed tasks data
type taskReply struct {
	err error
	cas gocb.Cas
}

// getTask returns a Task setup for a get
func getTask(key string, valuePtr interface{}) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionGet,
		key:    key,
		data:   valuePtr,
	}
}

// insertTask returns a Task setup for an insert
func insertTask(key string, value interface{}, expiry uint32) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionInsert,
		key:    key,
		data:   value,
		expiry: expiry,
	}
}

// removeTask returns a Task setup for remove
func removeTask(key string, cas gocb.Cas) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionRemove,
		key:    key,
		cas:    cas,
	}
}

// replaceTask returns a Task setup for replace
func replaceTask(key string, value interface{}, cas gocb.Cas, expiry uint32) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionReplace,
		key:    key,
		data:   value,
		cas:    cas,
		expiry: expiry,
	}
}

// touchTask returns a Task setup for a touch
func touchTask(key string, cas gocb.Cas, expiry uint32) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionTouch,
		key:    key,
		cas:    cas,
		expiry: expiry,
	}
}

// upsertTask returns a Task setup for an upsert
func upsertTask(key string, value interface{}, expiry uint32) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionUpsert,
		key:    key,
		data:   value,
		expiry: expiry,
	}
}

// EOF
