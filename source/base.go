package source

import (
	"errors"

	"github.com/terakoya76/read-replicas-outbox-miner/config"
	"github.com/terakoya76/read-replicas-outbox-miner/converters"
	"github.com/terakoya76/read-replicas-outbox-miner/tracker"
)

// SourceClient provide queries for outbox table
type SourceClient interface {
	Mine(startPos tracker.Position) ([]converters.InternalRow, error)
}

// BuildClient builds SourceClient for abstraction
func BuildClient() (SourceClient, error) {
	switch config.Source.Strategy {
	case "mysql":
		return BuildMySQLClient()
	default:
		return nil, errors.New("not supported data source")
	}
}
