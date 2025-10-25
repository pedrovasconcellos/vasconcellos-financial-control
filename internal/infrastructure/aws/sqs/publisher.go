package sqs

import (
	"context"
	"encoding/json"
	"sync"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"github.com/vasconcellos/finance-control/internal/domain/port"
)

type Publisher struct {
	client          *sqs.Client
	defaultQueueURL string
	queueCache      map[string]string
	mu              sync.RWMutex
}

var _ port.QueuePublisher = (*Publisher)(nil)

func NewPublisher(cfg aws.Config, queueURL string) *Publisher {
	client := sqs.NewFromConfig(cfg)
	return &Publisher{client: client, defaultQueueURL: queueURL, queueCache: map[string]string{}}
}

func (p *Publisher) Publish(ctx context.Context, queueName string, message port.QueueMessage) error {
	body := message.Payload
	if len(body) == 0 && len(message.Attributes) > 0 {
		encoded, err := json.Marshal(message.Attributes)
		if err != nil {
			return err
		}
		body = encoded
	}

	queueURL := p.resolveQueueURL(ctx, queueName)
	if queueURL == "" {
		queueURL = p.defaultQueueURL
	}

	attrs := map[string]types.MessageAttributeValue{}
	for key, value := range message.Attributes {
		attrs[key] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}

	_, err := p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:          aws.String(queueURL),
		MessageBody:       aws.String(string(body)),
		MessageAttributes: attrs,
	})
	return err
}

func (p *Publisher) resolveQueueURL(ctx context.Context, queueName string) string {
	if queueName == "" {
		return ""
	}

	p.mu.RLock()
	if url, ok := p.queueCache[queueName]; ok {
		p.mu.RUnlock()
		return url
	}
	p.mu.RUnlock()

	output, err := p.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{QueueName: aws.String(queueName)})
	if err != nil || output.QueueUrl == nil {
		return ""
	}

	url := aws.ToString(output.QueueUrl)
	p.mu.Lock()
	p.queueCache[queueName] = url
	p.mu.Unlock()
	return url
}
