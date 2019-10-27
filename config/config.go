package config

import (
	"errors"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

var once = new(sync.Once)

func init() {
	once.Do(func() {
		if err := initMust(); err != nil {
			panic(err)
		}

		if err := initOpt(); err != nil {
			panic(err)
		}
	})
}

type sourceSpec struct {
	Strategy string `default:"mysql"`
}

type trackerSpec struct {
	Strategy string `default:"mysql"`
}

// TrackKey must be AutoIncremental Key
// Default BatchSize fit to Kinesis#PutRecords API limitation
type minerSpec struct {
	Database  string `required:"true"`
	Table     string `required:"true"`
	TrackKey  string `required:"true"`
	BatchSize int64  `default:"500"`
}

type publisherSpec struct {
	Strategy string `default:"kinesis-data-stream"`
}

// Source represents information about DataSource Strategy
var Source *sourceSpec

// Tracker represents information about MySQL configuration for Tracker persisting source position
var Tracker *trackerSpec

// Miner represents information about Miner and its target
var Miner *minerSpec

// Publisher represents information where the mined events to be published
var Publisher *publisherSpec

// Load Configuration which is necessary in any Source/Publisher Strategy
func initMust() error {
	var ss sourceSpec
	if err := envconfig.Process("source", &ss); err != nil {
		return err
	}
	Source = &ss

	var ts trackerSpec
	if err := envconfig.Process("tracker", &ts); err != nil {
		return err
	}
	Tracker = &ts

	var ms minerSpec
	if err := envconfig.Process("miner", &ms); err != nil {
		return err
	}
	Miner = &ms

	var ps publisherSpec
	if err := envconfig.Process("publisher", &ps); err != nil {
		return err
	}
	Publisher = &ps

	return nil
}

type sourceMysqlSpec struct {
	Host     string `default:"127.0.0.1"`
	Port     int    `default:"3306"`
	User     string `required:"true"`
	Password string `required:"true"`
}

type trackerMysqlSpec struct {
	Host     string `default:"127.0.0.1"`
	Port     int    `default:"3306"`
	User     string `required:"true"`
	Password string `required:"true"`
	Name     string `required:"true"`
}

type kinesisPublisherSpec struct {
	Region       string `required:"true"`
	Endpoint     string `default:"127.0.0.1:4567"`
	StreamName   string `required:"true"`
	PartitionKey string `required:"true"`
}

// SourceMySQL represents information about DataSource MySQL Configuration
var SourceMySQL *sourceMysqlSpec

// TrackerMySQL represents information about MySQL Configuration for Tracker Storage
var TrackerMySQL *trackerMysqlSpec

// KinesisPublisher represents information about Kinesis Configuration where events to be published
var KinesisPublisher *kinesisPublisherSpec

func initOpt() error {
	switch Source.Strategy {
	case "mysql":
		var sms sourceMysqlSpec
		if err := envconfig.Process("source_mysql", &sms); err != nil {
			return err
		}
		SourceMySQL = &sms
	default:
		return errors.New("non-supported data source")
	}

	switch Tracker.Strategy {
	case "mysql":
		var tms trackerMysqlSpec
		if err := envconfig.Process("tracker_mysql", &tms); err != nil {
			return err
		}
		TrackerMySQL = &tms
	default:
		return errors.New("non-supported tracker storage")
	}

	switch Publisher.Strategy {
	case "kinesis-data-streams":
		var kps kinesisPublisherSpec
		if err := envconfig.Process("kinesis_publisher", &kps); err != nil {
			return err
		}
		KinesisPublisher = &kps
	default:
		return errors.New("non-supported publisher")
	}

	return nil
}
