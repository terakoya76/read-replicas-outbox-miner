package converters

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"

	"github.com/terakoya76/read-replicas-outbox-miner/tracker"
)

// InternalRow represents inner-data-type mediates *sqlx.Row and the data-format to be published
type InternalRow = map[string]interface{}

type jsonNullInt32 struct {
	sql.NullInt32
}

func (v jsonNullInt32) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Int32)
}

type jsonNullInt64 struct {
	sql.NullInt64
}

func (v jsonNullInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Int64)
}

type jsonNullFloat64 struct {
	sql.NullFloat64
}

func (v jsonNullFloat64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Float64)
}

type jsonNullTime struct {
	sql.NullTime
}

func (v jsonNullTime) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Time)
}

var jsonNullInt32Type = reflect.TypeOf(jsonNullInt32{})
var jsonNullInt64Type = reflect.TypeOf(jsonNullInt64{})
var jsonNullFloat64Type = reflect.TypeOf(jsonNullFloat64{})
var jsonNullTimeType = reflect.TypeOf(jsonNullTime{})
var nullInt32Type = reflect.TypeOf(sql.NullInt32{})
var nullInt64Type = reflect.TypeOf(sql.NullInt64{})
var nullFloat64Type = reflect.TypeOf(sql.NullFloat64{})
var nullTimeType = reflect.TypeOf(sql.NullTime{})

// SQL2ListMap converts *sqlx.Rows to intermediate-data-type
func SQL2ListMap(rows *sqlx.Rows) ([]InternalRow, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("Column error: %v", err)
	}

	tt, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("Column type error: %v", err)
	}

	types := make([]reflect.Type, len(tt))
	for i, tp := range tt {
		st := tp.ScanType()
		if st == nil {
			return nil, fmt.Errorf("Scantype is null for column: %v", err)
		}
		switch st {
		case nullInt32Type:
			types[i] = jsonNullInt32Type
		case nullInt64Type:
			types[i] = jsonNullInt64Type
		case nullFloat64Type:
			types[i] = jsonNullFloat64Type
		case nullTimeType:
			types[i] = jsonNullTimeType
		default:
			types[i] = st
		}
	}

	values := make([]interface{}, len(tt))
	data := make([]InternalRow, 0)

	for rows.Next() {
		for i := range values {
			values[i] = reflect.New(types[i]).Interface()
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan values: %v", err)
		}

		row := InternalRow{}
		for i := range values {
			row[columns[i]] = values[i]
		}
		data = append(data, row)
	}
	return data, nil
}

// Map2JSON converts intermediate-data-type to JSON
func Map2JSON(row InternalRow) ([]byte, error) {
	return json.Marshal(row)
}

// Cast2Position converts numeric data type to Position
func Cast2Position(val interface{}) tracker.Position {
	switch v := val.(type) {
	case *int32:
		return tracker.Position(int64(*v))
	case *int64:
		return tracker.Position(*v)
	default:
		panic(errors.New("non-supported track key DataType"))
	}
}
