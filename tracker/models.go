package tracker

// Position represents the AutoIncrement value coupled w/ TrackKey
type Position = int64

// Progress represents mining progress on DataSource
// Source has_one Progress
type Progress struct {
	ID           int64    `db:"id"`
	DatabaseName string   `db:"database_name"`
	TableName    string   `db:"table_name"`
	TrackKey     string   `db:"track_key"`
	Position     Position `db:"position"`
}
