package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Seq          int     `json:"seq"`
	User_Name    string  `json:"user_name" binding:"required"`
	User_Id      string  `json:"user_id" binding:"required"`
	Group_Id     string  `json:"group_id" binding:"required"`
	Trainer_Id   string  `json:"trainer_id" binding:"required"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User *string `json:"created_user"`
	Updated_User *string `json:"updated_user"`
}

func (u User) selectAllUser(db *sql.DB) (users []User, err error) {
	rows, err := db.Query(
		`SELECT seq, 
		user_name,
		user_id,
		group_id,
		trainer_id,
		created_date,
		updated_date,
		created_user,
		updated_user FROM t_user`)
	if err != nil {
		return
	}

	for rows.Next() {
		user := User{}
		rows.Scan(
			&u.Seq,
			&u.User_Name,
			&u.User_Id,
			&u.Group_Id,
			&u.Trainer_Id,
			&u.Created_Date,
			&u.Updated_Date,
			&u.Created_User,
			&u.Updated_User)
		users = append(users, user)
	}
	defer rows.Close()

	return
}

func getAllUserList(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		user := User{}
		users, err := user.selectAllUser(db)
		if err != nil {
			nullUsers := [0]User{}
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"parameters": gin.H{
					"body": nullUsers,
				},
			})
		} else {
			c.JSON(http.StatusOK, users)
		}
	}

	return resultFunc
}

func UserListRouter(router *gin.Engine, db *sql.DB) {
	user := router.Group("/api/user")

	user.GET("/all", getAllUserList(db))
}
