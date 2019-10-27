#!/bin/bash
ENDPOINT="http://localhost:4567"
SHARDID="shardId-000000000000"
ITER=$(aws kinesis get-shard-iterator \
    --stream-name test-stream \
    --shard-id $SHARDID \
    --shard-iterator-type LATEST \
    --endpoint-url $ENDPOINT \
    --no-verify-ssl \
    | jq -r .ShardIterator)

while :
do
    aws kinesis get-records \
        --shard-iterator $ITER \
        --endpoint-url $ENDPOINT \
        --no-verify-ssl \
        | jq '.Records | length' | xargs -I % echo "Get % records from shards"
    ITER=$(aws kinesis get-records \
        --shard-iterator $ITER \
        --endpoint-url $ENDPOINT \
        --no-verify-ssl \
        | jq -r .NextShardIterator)
    sleep 0.3
done
