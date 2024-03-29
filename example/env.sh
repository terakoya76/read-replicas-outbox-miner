# Source
export SOURCE_STRATEGY=mysql
export SOURCE_MYSQL_HOST=127.0.0.1
export SOURCE_MYSQL_PORT=3306
export SOURCE_MYSQL_USER=root
export SOURCE_MYSQL_PASSWORD=

# Tracker
export TRACKER_STRATEGY=mysql
export TRACKER_MYSQL_HOST=127.0.0.1
export TRACKER_MYSQL_PORT=3306
export TRACKER_MYSQL_USER=root
export TRACKER_MYSQL_PASSWORD=
export TRACKER_MYSQL_NAME=read_replicas_outbox_miner_db

# Publisher
export PUBLISHER_STRATEGY=kinesis-data-streams
export KINESIS_PUBLISHER_REGION=ap-northeast-1
export KINESIS_PUBLISHER_ENDPOINT=http://127.0.0.1:4567
export KINESIS_PUBLISHER_STREAMNAME=test-stream

# other
export AWS_CBOR_DISABLE=1
export MINER_TARGETS_CONFIG_PATH=./example/miner-targets.yaml