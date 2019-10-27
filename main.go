package main

import (
	"fmt"

	"github.com/terakoya76/read-replicas-outbox-miner/executor"
	"github.com/terakoya76/read-replicas-outbox-miner/publisher"
	"github.com/terakoya76/read-replicas-outbox-miner/source"
	"github.com/terakoya76/read-replicas-outbox-miner/tracker"
)

func main() {
	source, err := source.BuildClient()
	if err != nil {
		panic(err)
	}

	tracker, err := tracker.BuildTrackerClient()
	if err != nil {
		panic(err)
	}

	publisher, err := publisher.BuildPublisher()
	if err != nil {
		panic(err)
	}

	for {
		if err := executor.Exec(source, tracker, publisher); err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}
