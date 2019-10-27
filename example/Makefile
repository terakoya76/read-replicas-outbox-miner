mysql.setup:
	populator

kinesis.setup:
	aws kinesis create-stream \
		--stream-name test-stream \
		--shard-count 1 \
		--endpoint-url http://localhost:4567 \
		--no-verify-ssl

kinesis.desc:
	aws kinesis list-streams --endpoint-url http://localhost:4567 --no-verify-ssl
	aws kinesis describe-stream --stream-name test-stream --endpoint-url http://localhost:4567 --no-verify-ssl

kinesis.tail:
	./kinesis-tail.sh
