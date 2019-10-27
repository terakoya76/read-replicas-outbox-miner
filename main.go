package main

import (
	"context"
	"fmt"
	"os"

	"github.com/terakoya76/read-replicas-outbox-miner/config"
	"github.com/terakoya76/read-replicas-outbox-miner/executor"
)

func main() {
	ctx := context.Background()
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan executor.WorkerID)
	for i, target := range config.Miner.Targets {
		w, err := executor.BuildWorker(i, target, done)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// TODO:
		// cctx is not used in Worker currently
		// should be passed to each xxx_client
		go executor.Launch(cctx, w)
	}

	select {
	case wid := <-done:
		// restart worker
		target := config.Miner.Targets[wid]
		w, err := executor.BuildWorker(wid, target, done)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// TODO:
		// cctx is not used in Worker currently
		// should be passed to each xxx_client
		go executor.Launch(cctx, w)
	case <-cctx.Done():
		os.Exit(0)
	}
}
