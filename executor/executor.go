package executor

import (
	"fmt"

	"github.com/terakoya76/read-replicas-outbox-miner/config"
	"github.com/terakoya76/read-replicas-outbox-miner/publisher"
	"github.com/terakoya76/read-replicas-outbox-miner/source"
	"github.com/terakoya76/read-replicas-outbox-miner/tracker"
)

// Exec provides the operation for querying to DB and publishing evnets
func Exec(sc source.SourceClient, tc *tracker.TrackerClient, publisher publisher.Publisher) error {
	config := config.Miner
	var startPosition tracker.Position

	startPosition, err := tc.GetNextPosition(config.Database, config.Table)
	if err != nil {
		return err
	}

	records, err := sc.Mine(startPosition)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		fmt.Print("No records mined.\n")
		return nil
	}

	if err = tc.WithTx(func() error {
		for _, record := range records {
			published, err := publisher.Publish(record)
			if err != nil {
				return err
			}

			if published {
				tc.UpdatePosition(config.Database, config.Table, publisher.GetPosition())

				// TODO: not break and continue to process until the last record for performance
				break
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return err
}
