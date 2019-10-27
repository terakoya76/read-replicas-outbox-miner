package publisher

import (
	"errors"

	"github.com/terakoya76/read-replicas-outbox-miner/config"
	"github.com/terakoya76/read-replicas-outbox-miner/converters"
	"github.com/terakoya76/read-replicas-outbox-miner/tracker"
)

// Publisher provides event publishing
type Publisher interface {
	Publish(event converters.InternalRow) (bool, error)
	GetPosition() tracker.Position
}

// BuildPublisher builds Publisher for abstraction
func BuildPublisher() (Publisher, error) {
	switch config.Publisher.Strategy {
	case "kinesis-data-streams":
		return BuildKinesisDataStreamsPublisher()
	default:
		return nil, errors.New("not supported Publisher Strategy")
	}
}
