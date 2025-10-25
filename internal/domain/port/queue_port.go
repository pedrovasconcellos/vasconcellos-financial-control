package port

import "context"

type QueueMessage struct {
	ID         string
	Payload    []byte
	Attributes map[string]string
}

type QueuePublisher interface {
	Publish(ctx context.Context, queueName string, message QueueMessage) error
}
