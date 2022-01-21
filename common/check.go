package common

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

// This function is to Check Duplicated Name in category or exercise table's title
func DuplicatedTitleCheck(table, title string, db *sql.DB) (count int, err error) {
	query := fmt.Sprintf("SELECT COUNT(title) FROM %s WHERE title = '%s'", table, title)
	row := db.QueryRow(query)
	err = row.Scan(&count)
	if err != nil {
		return 1, err
	}

	return count, nil
}

func DuplicatedUserIdCheck(table, user_id string, db *sql.DB) (count int, err error) {
	query := fmt.Sprintf("SELECT COUNT(user_id) FROM %s WHERE user_id = '%s'", table, user_id)
	row := db.QueryRow(query)
	err = row.Scan(&count)
	if err != nil {
		return 1, err
	}

	return count, nil
}

func GetQueryString(c *gin.Context) (trainer_id string, group_name string) {
	trainer_id = c.Query("trainer_id")
	group_name = c.Query("group_name")

	return
}
