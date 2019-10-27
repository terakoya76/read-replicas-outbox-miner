package tracker

import (
	"database/sql"
	"fmt"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/terakoya76/read-replicas-outbox-miner/config"
)

// MySQLTrackerStorage implements TrackerStorage for MySQL
type MySQLTrackerStorage struct {
	db sqlxCommon
	tx sqlxCommon
}

// BuildMySQLTrackerStorage builds MySQL specific TrackerStorage
func BuildMySQLTrackerStorage() (*MySQLTrackerStorage, error) {
	ci := buildConnectInfo()
	db, err := sqlx.Connect("mysql", ci)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database %s on mysql: %+v", &config.TrackerMySQL.Name, err)
	}

	return &MySQLTrackerStorage{db: db, tx: nil}, nil
}

func buildConnectInfo() string {
	config := config.TrackerMySQL
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
}

// DB provides access to db handler properly to process tx operation
func (mts *MySQLTrackerStorage) DB() sqlxCommon {
	if mts.tx != nil {
		return mts.tx
	}
	return mts.db
}

// WithTx provides transactional operation for given func
func (mts *MySQLTrackerStorage) WithTx(f func() error) (err error) {
	var emptyDB *sqlx.DB
	var tx *sqlx.Tx
	if _db, ok := mts.db.(sqlxDB); ok && _db != nil && _db != emptyDB {
		tx, err = _db.Beginx()
		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
				panic(err)
			}
		}()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	mts.tx = interface{}(tx).(sqlxCommon)
	defer func() {
		mts.tx = nil
	}()

	var emptyTx *sqlx.Tx
	if err := f(); err != nil {
		fmt.Println(err)

		if _tx, ok := mts.tx.(sqlxTx); ok && _tx != nil && _tx != emptyTx {
			_tx.Rollback()
		}
		return err
	} else {
		if _tx, ok := mts.tx.(sqlxTx); ok && _tx != nil && _tx != emptyTx {
			_tx.Commit()
		}
		return nil
	}
}

func (mts *MySQLTrackerStorage) Prepare() error {
	config := config.Miner
	for _, target := range config.Targets {
		if err := mts.prepareProgress(config.Database, target.Table, target.TrackKey); err != nil {
			return err
		}
	}
	return nil
}

func (mts *MySQLTrackerStorage) prepareProgress(dbName string, tableName string, trackKey string) error {
	_, err := mts.GetProgress(dbName, tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			newProgress := Progress{
				ID:           0,
				DatabaseName: dbName,
				TableName:    tableName,
				TrackKey:     trackKey,
				Position:     0,
			}
			if _, err := mts.DB().NamedExec(
				"INSERT INTO progresses (id, database_name, table_name, track_key, position) VALUES (:id, :database_name, :table_name, :track_key, :position)",
				&newProgress,
			); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

// GetProgress fetches progress associated w/ given source
func (mts *MySQLTrackerStorage) GetProgress(dbName string, tableName string) (*Progress, error) {
	prgs := Progress{}
	if err := mts.DB().Get(
		&prgs,
		fmt.Sprintf(
			"SELECT * FROM progresses WHERE database_name = %q AND table_name = %q",
			dbName,
			tableName,
		),
	); err != nil {
		return nil, err
	}

	return &prgs, nil
}

// UpdateProgress updates progress's position
func (mts *MySQLTrackerStorage) UpdateProgress(prgs *Progress) (*Progress, error) {
	if _, err := mts.DB().NamedExec(
		"UPDATE progresses SET position = :position WHERE database_name = :database_name AND table_name = :table_name", prgs); err != nil {
		return nil, err
	}

	return prgs, nil
}
