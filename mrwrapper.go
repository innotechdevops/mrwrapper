package mrwrapper

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func Count(conn *sqlx.DB, query string, args ...any) int64 {
	var id int64 = 0
	err := conn.Get(&id, query, args...)
	if err == nil {
		return id
	}
	fmt.Println("Count error:", err)
	return 0
}

func SelectOne[R any](conn *sqlx.DB, query string, args ...any) R {
	var row R
	var rows []R
	if err := conn.Select(&rows, query, args...); err != nil {
		fmt.Println("SelectOne error:", err)
		return row
	}
	if len(rows) > 0 {
		return rows[0]
	}
	return row
}

func SelectList[R any](conn *sqlx.DB, query string, args ...any) []R {
	var rows []R
	if err := conn.Select(&rows, query, args...); err != nil {
		fmt.Println("SelectList error:", err)
		return []R{}
	}
	if len(rows) > 0 {
		return rows
	}
	return []R{}
}

func Create(conn *sqlx.DB, query string, dest []any, args ...any) (*sql.Tx, error) {
	tx, err := conn.Begin()
	if err != nil {
		return tx, err
	}

	_, err = tx.Exec(query, args...)
	_ = tx.QueryRow("SELECT LAST_INSERT_ID()").Scan(dest...)

	if err == nil {
		return tx, nil
	}

	// Roll back the transaction in error case of foreign key violation
	_ = tx.Rollback()

	return tx, err
}

func Update(conn *sqlx.DB, query string, set string, params map[string]any) (*sql.Tx, error) {
	tx, err := conn.Begin()
	if err != nil {
		return tx, err
	}

	if set == "" {
		return tx, errors.New("is not data to update")
	}
	query = fmt.Sprintf(query, set[1:])

	query, args, err := sqlx.Named(query, params)
	if err != nil {
		return tx, err
	}
	query = conn.Rebind(query)

	rs, err := tx.Exec(query, args...)
	if err != nil {

		// Roll back the transaction in error case of foreign key violation
		_ = tx.Rollback()

		return tx, err
	}

	_, err = rs.RowsAffected()
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func Delete(conn *sqlx.DB, query string, args ...any) (*sql.Tx, error) {
	tx, err := conn.Begin()
	if err != nil {
		return tx, err
	}

	rs, err := tx.Exec(query, args...)
	if err != nil {

		// Roll back the transaction in error case of foreign key violation
		_ = tx.Rollback()

		return tx, err
	}

	row, err := rs.RowsAffected()
	if err != nil {
		return tx, err
	}

	if row <= 0 {
		return tx, err
	}

	return tx, nil
}
