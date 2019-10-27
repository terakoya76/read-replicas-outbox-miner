# How to use
```bash
$ go get github.com/terakoya76/populator
$ docker-compose up
$ make mysql.setup
$ make kinesis.setup
$ source env.sh
$ cd .. && git checkout
$ go build
$ ./read-replicas-outbox-miner

# on another tab
$ make kinesis.tail
```

# Notice
For generating seed in MySQL, this Makefile uses [populator](https://github.com/terakoya76/populator), which executes INSERT statement along w/ the yaml configuration.

So, before typing `make mysql.setup`, you'll need to install populator itself by `go get github.com/terakoya76/populator`.
