package executor

import (
	"context"
	"fmt"

	"github.com/terakoya76/read-replicas-outbox-miner/config"
	"github.com/terakoya76/read-replicas-outbox-miner/publisher"
	"github.com/terakoya76/read-replicas-outbox-miner/source"
	"github.com/terakoya76/read-replicas-outbox-miner/tracker"
)

type WorkerID = int

type Worker struct {
	ID            WorkerID
	Target        *config.MinerTarget
	SourceClient  source.SourceClient
	TrackerClient *tracker.TrackerClient
	Publisher     publisher.Publisher
	Done          chan WorkerID
}

func BuildWorker(id WorkerID, target *config.MinerTarget, done chan WorkerID) (*Worker, error) {
	sc, err := source.BuildClient()
	if err != nil {
		return nil, err
	}

	tc, err := tracker.BuildTrackerClient()
	if err != nil {
		return nil, err
	}

	pub, err := publisher.BuildPublisher(target)
	if err != nil {
		return nil, err
	}

	w := Worker{
		ID:            id,
		Target:        target,
		SourceClient:  sc,
		TrackerClient: tc,
		Publisher:     pub,
		Done:          done,
	}
	return &w, nil
}

func (w *Worker) run(ctx context.Context) {
	for {
		if err := w.mine(); err != nil {
			fmt.Println(err)
		}
	}
}

func (w *Worker) mine() error {
	dbName := config.Miner.Database
	startPosition, err := w.TrackerClient.GetNextPosition(dbName, w.Target.Table)
	if err != nil {
		return err
	}

	records, err := w.SourceClient.Fetch(startPosition, w.Target)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		fmt.Print("No records mined.\n")
		return nil
	}

	if err = w.TrackerClient.WithTx(func() error {
		for _, record := range records {
			published, err := w.Publisher.Publish(record)
			if err != nil {
				return err
			}

			if published {
				w.TrackerClient.UpdatePosition(dbName, w.Target.Table, w.Publisher.GetPosition())

				// TODO: not break and continue to process until the last record for performance
				break
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func Launch(ctx context.Context, w *Worker) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			w.Done <- w.ID
		}
	}()
	w.run(ctx)
}
