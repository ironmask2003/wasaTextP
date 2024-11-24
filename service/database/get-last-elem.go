package database

import (
	"database/sql"
	"errors"
)

// Function used to take the last element of a table
func (db *appdbimpl) GetLastElem(query string) (int, error) {
	var _maxID = sql.NullInt64{Int64: 0, Valid: false}
	row, err := db.c.Query(query)
	if err != nil {
		return 0, err
	}

	var maxID int
	for row.Next() {
		if row.Err() != nil {
			return 0, err
		}

		err = row.Scan(&_maxID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return 0, err
		}

		if !_maxID.Valid {
			maxID = 0
		} else {
			maxID = int(_maxID.Int64)
		}
	}

	return maxID, nil
}
