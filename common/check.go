package common

import (
	"database/sql"
	"fmt"
)

// This function is to Check Duplicated Name in category or exercise table's title
func DuplicatedTitleCheck(table, title string, db *sql.DB) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(title) FROM %s WHERE title = '%s'", table, title)
	row := db.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		return 1, err
	}

	return count, nil
}
