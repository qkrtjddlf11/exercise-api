package common

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var DuplicatedTitle = errors.New("Duplicated Title")

// This function is to Check Duplicated Name in category or exercise table's title
func IsDuplicatedTitle(table, title string, db *sql.DB) (int, error) {
	var count int
	row := db.QueryRow(`SELECT COUNT(title) FROM ` + table + ` WHERE title like '%` + title + `%'`)
	err := row.Scan(&count)
	if err != nil {
		return count + 1, errors.Wrap(err, "Failed to select from database")
	} else {
		switch count {
		case 1:
			return count, DuplicatedTitle

		default:
			return count, nil
		}
	}
}

func IsDuplicatedUserId(table, user_id string, db *sql.DB) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(user_id) FROM %s WHERE user_id = '%s'", table, user_id)
	row := db.QueryRow(query)
	err := row.Scan(&count)
	if err != nil {
		return 1, errors.Wrap(err, "Duplicated User ID")
	}

	return count, nil
}

func GetQueryString(c *gin.Context) (string, string) {
	trainer_id := c.Query("trainer_id")
	group_name := c.Query("group_name")

	return trainer_id, group_name
}
