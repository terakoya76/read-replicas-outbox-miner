package source

import (
	"fmt"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/terakoya76/read-replicas-outbox-miner/config"
	"github.com/terakoya76/read-replicas-outbox-miner/converters"
	"github.com/terakoya76/read-replicas-outbox-miner/tracker"
)

// MySQLClient implements SourceClient for MySQL
type MySQLClient struct {
	*sqlx.DB
}

// BuildMySQLClient builds MySQL specific SourceClient
func BuildMySQLClient() (*MySQLClient, error) {
	ci := buildConnectInfo()
	db, err := sqlx.Connect("mysql", ci)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database %s on mysql: %+v", config.Miner.Database, err)
	}

	return &MySQLClient{db}, nil
}

func buildConnectInfo() string {
	dbName := config.Miner.Database
	config := config.SourceMySQL
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		dbName,
	)
}

// Fetch fetches batch-sized number of records from the outbox table
func (c *MySQLClient) Fetch(startPos tracker.Position, target *config.MinerTarget) ([]converters.InternalRow, error) {
	startID := int64(startPos)
	endID := startID + target.BatchSize
	rows, err := c.Queryx(fmt.Sprintf(
		"SELECT * FROM %s WHERE %s BETWEEN %d AND %d;",
		target.Table,
		target.TrackKey,
		startID,
		endID,
	))
	if err != nil {
		return nil, err
	}

	return converters.SQL2ListMap(rows)
}
