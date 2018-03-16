package asynctx

import (
	"github.com/streadway/amqp"
)

// task struct holds data for message passed to mq and pulled back out
type task struct {
	queue string
	msg   amqp.Publishing
}

// EOF
