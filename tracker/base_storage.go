package tracker

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/terakoya76/read-replicas-outbox-miner/config"
)

// TrackerStorage provide access to persisted Progress
type TrackerStorage interface {
	WithTx(f func() error) error
	Prepare() error
	GetProgress(dbName string, tableName string) (*Progress, error)
	UpdateProgress(progress *Progress) (*Progress, error)
}

// BuildTrackerStorage builds TrackerStorage for abstraction
func BuildTrackerStorage() (TrackerStorage, error) {
	var ts TrackerStorage
	switch config.Tracker.Strategy {
	case "mysql":
		mts, err := BuildMySQLTrackerStorage()
		if err != nil {
			return nil, err
		}
		ts = mts
	default:
		return nil, errors.New("not supported data source")
	}

	if err := ts.Prepare(); err != nil {
		return nil, err
	}

	return ts, nil
}

type sqlxCommon interface {
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type sqlxDB interface {
	Beginx() (*sqlx.Tx, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

type sqlxTx interface {
	Commit() error
	Rollback() error
}
