package publisher

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"

	"github.com/terakoya76/read-replicas-outbox-miner/config"
	"github.com/terakoya76/read-replicas-outbox-miner/converters"
	"github.com/terakoya76/read-replicas-outbox-miner/tracker"
)

const (
	// 500 records
	LIMIT_RECORDS_PER_REQUEST = 500
	// 1 MB
	LIMIT_SIZE_PER_RECORD = 1000000
	// 5 MB
	LIMIT_SIZE_PER_REQUEST = 5000000
	// 4 MB
	PUBLISH_READINESS_THRESHOLD = 4000000
)

// KinesisDataStreamsPublisher implements Publisher for KinesisDataStreams
type KinesisDataStreamsPublisher struct {
	*kinesis.Kinesis
	streamName   *string
	partitionKey *string
	buffer       []*kinesis.PutRecordsRequestEntry
	counter      int64
	position     tracker.Position
}

// BuildKinesisDataStreamsPublisher builds KinesisDataStreams specific Publisher
func BuildKinesisDataStreamsPublisher() (*KinesisDataStreamsPublisher, error) {
	err := validateClient()
	if err != nil {
		return nil, err
	}

	sess := session.Must(session.NewSession())
	config := config.KinesisPublisher
	kc := kinesis.New(
		sess,
		aws.NewConfig().
			WithRegion(config.Region).
			WithEndpoint(config.Endpoint),
	)

	streamName := aws.String(config.StreamName)
	partitionKey := aws.String(config.PartitionKey)
	_, err = kc.DescribeStream(&kinesis.DescribeStreamInput{StreamName: streamName})
	if err != nil {
		panic(err)
	}

	kp := KinesisDataStreamsPublisher{
		Kinesis:      kc,
		streamName:   streamName,
		partitionKey: partitionKey,
		buffer:       make([]*kinesis.PutRecordsRequestEntry, 0),
		counter:      0,
		position:     0,
	}
	return &kp, nil
}

func validateClient() error {
	if config.Miner.BatchSize > LIMIT_RECORDS_PER_REQUEST {
		return errors.New("KinesisDataStreams#PutRecords accepts 500 records per request")
	}
	return nil
}

// Publish publishes events to the destination.
// KinesisDataStreamsPublisher works like events buffer for performance.
// it returns bool whether actually call Kinesis PutRecords API or just buffered event in its buffer.
func (kp *KinesisDataStreamsPublisher) Publish(event converters.InternalRow) (bool, error) {
	published := false

	ev, err := converters.Map2JSON(event)
	if err != nil {
		return published, err
	}

	if err := kp.stuff(ev); err != nil {
		return published, err
	}

	kp.track(event)
	if kp.ready() {
		if err := kp.publish(); err != nil {
			return published, err
		}

		published = true
		return published, nil
	}

	return published, nil
}

// stuffs events in the buffer as data type for KinesisDataStreams.
// Returns error if buffer already reached the limitation of PutRecords API call.
func (kp *KinesisDataStreamsPublisher) stuff(event []byte) error {
	entry := &kinesis.PutRecordsRequestEntry{
		Data:         event,
		PartitionKey: kp.partitionKey,
	}

	// We chose not to stop execution when over 1MB record are passed, and just ignore it.
	err := kp.acceptable(entry)
	if err != nil {
		fmt.Println(err)
	} else {
		kp.buffer = append(kp.buffer, entry)
	}

	kp.counter++
	return nil
}

func (kp *KinesisDataStreamsPublisher) acceptable(entry *kinesis.PutRecordsRequestEntry) error {
	bytes, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	if binary.Size(bytes) > LIMIT_SIZE_PER_RECORD {
		return errors.New("record's size cannot be over 1MB")
	}

	bytes, err = json.Marshal(kp.buffer)
	if err != nil {
		return err
	}
	if binary.Size(bytes) > LIMIT_SIZE_PER_REQUEST {
		// This case is not expected.
		// Before reached entire request size limitation, events should be pulished
		// by a little bit tight readiness threshold(4MB).
		panic(errors.New("record's size cannot be over 5MB"))
	}

	return nil
}

// Save last-stuffed event's position to make usable from tracker
func (kp *KinesisDataStreamsPublisher) track(event converters.InternalRow) {
	// TODO: check if event does not have trackKey attr
	kp.position = converters.Cast2Position(event[config.Miner.TrackKey])
}

func (kp *KinesisDataStreamsPublisher) ready() bool {
	bytes, err := json.Marshal(kp.buffer)
	if err != nil {
		panic(err)
	}
	return kp.counter == config.Miner.BatchSize || binary.Size(bytes) >= PUBLISH_READINESS_THRESHOLD
}

func (kp *KinesisDataStreamsPublisher) publish() error {
	if _, err := kp.PutRecords(&kinesis.PutRecordsInput{
		Records:    kp.buffer,
		StreamName: kp.streamName,
	}); err != nil {
		return err
	}

	kp.buffer = make([]*kinesis.PutRecordsRequestEntry, 0)
	kp.counter = 0
	return nil
}

// GetPosition returns the position of last event in buffer
func (kp *KinesisDataStreamsPublisher) GetPosition() tracker.Position {
	return kp.position
}
