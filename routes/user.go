package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qkrtjddlf11/exercise-api/common"
)

/*
type User struct {
	Seq          int     `json:"seq"`
	Name         string  `json:"name"`
	Id           string  `json:"id"`
	Group_Name   string  `json:"group_name"`
	Trainer_Id   string  `json:"trainer_id"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User *string `json:"created_user"`
	Updated_User *string `json:"updated_user"`
	Use_Yn       *string `json:"use_yn"`
	Password     *string `json:"password"`
	Email        *string `json:"email"`
}
*/

type User struct {
	Seq          int    `json:"seq"`
	Name         string `json:"name"`
	Id           string `json:"id"`
	Group_Name   string `json:"group_name"`
	Trainer_Id   string `json:"trainer_id"`
	Created_Date string `json:"created_date"`
	Updated_Date string `json:"updated_date"`
	Created_User string `json:"created_user"`
	Updated_User string `json:"updated_user"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Use_Yn       string `json:"use_yn"`
}

func (u User) selectAllUser(db *sql.DB) (users []User, err error) {
	rows, err := db.Query(
		`SELECT 
			seq, 
			name,
			trainer_id,
			group_name,
			created_date,
			updated_date,
			created_user,
			updated_user,
			email,
			use_yn FROM t_user
		WHERE trainer_id = ? AND group_name = ?`, u.Trainer_Id, u.Group_Name)
	if err != nil {
		return
	}

	for rows.Next() {
		var user User
		rows.Scan(
			&user.Seq,
			&user.Name,
			&user.Trainer_Id,
			&user.Group_Name,
			&user.Created_Date,
			&user.Updated_Date,
			&user.Created_User,
			&user.Updated_User,
			&user.Email,
			&user.Use_Yn)
		users = append(users, user)
	}
	defer rows.Close()

	return
}

func getAllUserList(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		trainer_id, group_name := common.GetQueryString(c)
		user := User{Trainer_Id: trainer_id, Group_Name: group_name}
		users, err := user.selectAllUser(db)
		if err != nil {
			nullUsers := [0]User{}
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, nullUsers))
		} else {
			c.JSON(http.StatusOK, common.SucceedResponse(users))
		}
	}

	return resultFunc
}

func UserListRouter(router *gin.Engine, db *sql.DB) {
	user := router.Group("/api/user")

	user.GET("/all", getAllUserList(db))
}
