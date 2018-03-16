package pool

import "github.com/mavricknz/ldap"

// action holds
const (
	actionAdd = iota
	actionBind
	actionDelete
	actionMod
	actionModDn
	actionSearch
)

// Task holds an actual task details
type taskRequest struct {
	reply  chan taskReply
	action int
	add    *ldap.AddRequest
	bind   *bindData
	delete *ldap.DeleteRequest
	mod    *ldap.ModifyRequest
	modDn  *ldap.ModDnRequest
	search *ldap.SearchRequest
}

// TaskReply holds a completed tasks data
type taskReply struct {
	err   error
	count int
	sr    *ldap.SearchResult
}

// bindData holds username and password for testing a Bind Task
type bindData struct {
	user string
	pass string
}

// addTask returns a task setup to for as an Add request
func addTask(req *ldap.AddRequest) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionAdd,
		add:    req,
	}
}

// bindTask returns a task setup to for as a Bind request
func bindTask(username, password string) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionBind,
		bind:   &bindData{username, password},
	}
}

// deleteTask returns a task setup to for as a Delete request
func deleteTask(delReq *ldap.DeleteRequest) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionDelete,
		delete: delReq,
	}
}

// modTask returns a task setup to for as a Mod request
func modTask(req *ldap.ModifyRequest) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionMod,
		mod:    req,
	}
}

// modDnTask returns a task setup to for as a ModDN request
func modDnTask(req *ldap.ModDnRequest) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionModDn,
		modDn:  req,
	}
}

// searchTask returns a task setup to for as a Search request
func searchTask(req *ldap.SearchRequest) taskRequest {
	return taskRequest{
		reply:  make(chan taskReply, 1),
		action: actionSearch,
		search: req,
	}
}

// EOF
