package rabbitmq

import (
	"errors"
	"time"

	"github.com/streadway/amqp"
)

// CreateBodyRPC is a helper to return a body ready to go
func CreateBodyRPC(content, id, reply string, persist bool, body []byte) (amqp.Publishing, error) {
	if content == "" {
		return amqp.Publishing{}, errors.New("content type can not be empty")
	}

	if id == "" {
		return amqp.Publishing{}, errors.New("correlation id can not be empty")
	}

	if reply == "" {
		return amqp.Publishing{}, errors.New("reply to can not be empty")
	}

	msg := amqp.Publishing{
		ContentType:   content,
		CorrelationId: id,
		ReplyTo:       reply,
		Body:          body,
		Timestamp:     time.Now(),
	}

	if persist == true {
		msg.DeliveryMode = amqp.Persistent
	}

	return msg, nil
}

// CreateBodyRPCReply is a helper to return a body ready to go
func CreateBodyRPCReply(content, id string, persist bool, body []byte) (amqp.Publishing, error) {
	if content == "" {
		return amqp.Publishing{}, errors.New("content type can not be empty")
	}

	if id == "" {
		return amqp.Publishing{}, errors.New("correlation id can not be empty")
	}

	msg := amqp.Publishing{
		ContentType:   content,
		CorrelationId: id,
		Body:          body,
		Timestamp:     time.Now(),
	}

	if persist == true {
		msg.DeliveryMode = amqp.Persistent
	}

	return msg, nil
}

// EOF
