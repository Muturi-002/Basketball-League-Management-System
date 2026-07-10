package database

import (
	"database/sql"
	"fmt"
	"time"
)

func queryList[T any](connectContext, queryContext, iterateContext, query string, scan func(*sql.Rows) (T, error), args ...any) ([]T, error) {
	conn, err := EnsureConnected()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", connectContext, err)
	}

	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", queryContext, err)
	}
	defer rows.Close()

	items := make([]T, 0)
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", iterateContext, err)
	}

	return items, nil
}

func nullStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}

func nullInt64Ptr(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	result := value.Int64
	return &result
}

func nullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	result := value.Time
	return &result
}
