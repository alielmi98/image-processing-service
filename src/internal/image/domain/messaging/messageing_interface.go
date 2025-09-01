package messaging

import (
	"context"

	"github.com/alielmi98/image-processing-service/pkg/contracts"
)

type MessageConsumer interface {
	Start(topic string) error
	Stop() error
}
type ConsumerHandler interface {
	HandleProcessingResult(ctx context.Context, result *contracts.ProcessingResult) error
}

type MessageSender interface {
	SendMessage(ctx context.Context, message *contracts.ProcessingMessage) error
}
