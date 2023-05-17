package user

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/google/uuid"
)

type EventStream struct {
	streamName    *string
	kinesisClient *kinesis.Client
}

func NewEventStream(cfg aws.Config, streamName string) *EventStream {
	kinesisClient := kinesis.NewFromConfig(cfg)

	return &EventStream{
		streamName:    aws.String(streamName),
		kinesisClient: kinesisClient,
	}
}

func (e *EventStream) publish(ctx context.Context, data []byte) error {
	input := &kinesis.PutRecordInput{
		Data:         data,
		PartitionKey: aws.String(uuid.NewString()),
		StreamName:   e.streamName,
	}

	_, err := e.kinesisClient.PutRecord(ctx, input)

	return err
}
