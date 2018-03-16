### Couchbase Lib Pool

This extension is a wrapper around the underlying Couchbase Lib.  It provides a pool of connections.

#### Usage

##### Setup Couchbase Config

    cfg := couchbase.DefaultConfig()

    cfg := couchbase.Config{
      Cluster: "couchbase://localhost",
      Bucket:  "bucket",
      Pass:    "password",
    }

##### Setup Couchbase Config

    poolCfg := pool.DefaultConfig()

    poolCfg := pool.Config{
      Min: 1, # Minimum number of workers.  Must be 1 or greater.
      Max: 10, # Maximum number of workers.  Must be equal to or greater than the Min.
      QueueSize: 10, # Channel queue size.
      Timeout: 50, # The number of Milliseconds before a Schedule call should timeout.
      Expiry: 30, # The number of Seconds since the last task before a worker shutsdown.
    }

##### Create a new Client

    cbpool, err := pool.InitClient(&poolCfg, &cfg)

##### Helper functions use similar signatures as found in gocb library, named differently as they return a TaskReply struct ptr

    cas, err := cbpool.Get(key, valuePtr)

    cas, err := cbpool.Insert(key, value, expiry)

    cas, err := cbpool.Remove(key, cas)

    cas, err := cbpool.Replace(key, value, cas, expiry)

    cas, err := cbpool.Touch(key, cas, expiry)

    cas, err := cbpool.Upsert(key, value, expiry)

##### Shutdown a running pool

    err := cbpool.Stop()

##### Functions to support Metrics

    numWorkers := cbpool.GetWorkerCount()
    queueSize := cbpool.GetQueueCount()

##### Task and TaskReply

The pool provides a channel on which a Task is placed.  This contains variables for each of the cbgo methods input requirements, plus a task type and reply channel.  There are helper functions to set these up as well.  The calling function should then block on the reply channel, over which a TaskReply that contains the reply details and an err.

#### Task Breakdown

##### Create a Task, there there three methods.

    task := pool.Task{
      Reply: make(chan pool.TaskReply, 1),
      Action:  pool.ActionGet,
      Key:  "keystring",
      Data:  &getData,
    }

    task := pool.DefaultTask
    task.Action = pool.ActionGet
    task.Key = "keystring"
    task.Data = &getData

    task := pool.GetTask("keystring", &getData)

##### Schedule the Task

    err := cbpool.Schedule(task)

##### Block on Reply Channel

    reply := <-task.Reply

##### Process TaskReply

    if reply.Err != nil {
      log.Fatalln(reply.Err.Error())
    }

    fmt.Println(task.Cas)
    fmt.Printf("%+v\n", getData)

