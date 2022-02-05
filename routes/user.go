package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/qkrtjddlf11/exercise-api/common"
)

type User struct {
	Seq          int     `json:"seq"`
	Name         string  `json:"name"`
	Id           string  `json:"id"`
	Group_Name   string  `json:"group_name"`
	Trainer_Id   string  `json:"trainer_id"`
	Created_Date *string `json:"created_date"`
	Updated_Date *string `json:"updated_date"`
	Created_User string  `json:"created_user"`
	Updated_User string  `json:"updated_user"`
	Password     string  `json:"password"`
	Email        *string `json:"email"`
	Use_Yn       string  `json:"use_yn"`
}

type UserList struct {
	Seq          int    `json:"seq"`
	Name         string `json:"name"`
	Id           string `json:"id"`
	Group_Name   string `json:"group_name"`
	Trainer_Id   string `json:"trainer_id"`
	Trainer_Name string `json:"trainer_name"`
	Email        string `json:"email"`
	Use_Yn       string `json:"use_yn"`
}

type ManageUserList struct {
	Name         string `json:"name"`
	Id           string `json:"id"`
	Trainer_Name string `json:"trainer_name"`
	Group_Name   string `json:"group_name"`
}

func (u UserList) selectAllUser(db *sql.DB, trainer_id, group_name string) ([]UserList, error) {
	var users []UserList
	nullUsers := [0]UserList{}
	rows, err := db.Query(
		`SELECT 
			seq, 
			name,
			id,
			trainer_id,
			group_name,
			email FROM t_user
		WHERE trainer_id = ? AND group_name = ? AND use_yn = 'Y'`, trainer_id, group_name)
	if err != nil {
		return nullUsers[:], errors.Wrap(err, "Failed to select From Database")
	}

	for rows.Next() {
		var user UserList
		rows.Scan(
			&user.Seq,
			&user.Name,
			&user.Id,
			&user.Trainer_Id,
			&user.Group_Name,
			&user.Email)
		users = append(users, user)
	}
	defer rows.Close()

	if len(users) > 0 {
		return users, nil
	} else {
		return nullUsers[:], errors.Wrap(err, "Users less than 1")
	}
}

func selectManageUser(db *sql.DB, trainer_id, group_name string) ([]ManageUserList, error) {
	var users []ManageUserList
	nullUsers := [0]ManageUserList{}
	rows, err := db.Query(
		`SELECT 
			u.name,
			u.id,
			t.name,
			u.group_name FROM t_user u left join t_trainer t on u.trainer_id = t.id
		WHERE u.trainer_id = ? AND u.group_name = ? AND u.use_yn = 'Y'`, trainer_id, group_name)
	if err != nil {
		return nullUsers[:], errors.Wrap(err, "Failed to select From Database")
	}

	for rows.Next() {
		var user ManageUserList
		rows.Scan(
			&user.Name,
			&user.Id,
			&user.Trainer_Name,
			&user.Group_Name)
		users = append(users, user)
	}
	defer rows.Close()

	if len(users) > 0 {
		return users, nil
	} else {
		return nullUsers[:], errors.Wrap(err, "Users less than 1")
	}
}

func selectDetailUser(db *sql.DB, user_seq string) (User, error) {
	var user User
	Seq, err := strconv.Atoi(user_seq)
	if err != nil {
		return user, errors.Wrap(err, "Failed to convert user_seq to int")
	}

	row := db.QueryRow(
		`SELECT * FROM t_user
		WHERE seq = ?`, Seq)

	err = row.Scan(
		&user.Seq,
		&user.Name,
		&user.Id,
		&user.Group_Name,
		&user.Trainer_Id,
		&user.Created_Date,
		&user.Updated_Date,
		&user.Created_User,
		&user.Updated_User,
		&user.Password,
		&user.Email,
		&user.Use_Yn)
	if err != nil {
		return user, errors.Wrap(err, "Failed to select from Database")
	}

	return user, nil
}

func (u User) addUser(db *sql.DB) (int, error) {
	var id int

	row := db.QueryRow(
		`SELECT 
			seq FROM t_user 
		ORDER BY seq 
		DESC LIMIT 1`)

	err := row.Scan(&id)
	if err != nil {
		return id, errors.Wrap(err, "Failed to query From Database")
	}

	stmt, err := db.Prepare(
		`INSERT INTO
			t_user(
				name,
				id,
				group_name,
				trainer_id,
				created_user,
				updated_user,
				use_yn)
			VALUE(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return id, errors.Wrap(err, "Failed to prepare to Database")
	}

	result, err := stmt.Exec(
		u.Name, "user"+strconv.Itoa(id), u.Group_Name, u.Trainer_Id, u.Trainer_Id, u.Trainer_Id, u.Use_Yn)
	if err != nil {
		return id, errors.Wrap(err, "Failed to execute to Database")
	}

	seq, err := result.LastInsertId()
	if err != nil {
		return id, errors.Wrap(err, "Failed to insert last id to Database")
	}

	id = int(seq)
	defer stmt.Close()

	return id, nil
}

func deleteUser(db *sql.DB, user_seq string) (int, error) {
	seq, err := strconv.Atoi(user_seq)
	if err != nil {
		return seq, errors.Wrap(err, "Failed to convert String to Int")
	}

	stmt, err := db.Prepare(`
		DELETE FROM t_user WHERE seq = ?`)
	if err != nil {
		return seq, errors.Wrap(err, "Failed to prepare")
	}

	result, err := stmt.Exec(seq)
	if err != nil {
		return seq, errors.Wrap(err, "Failed to delete category")
	}

	row, err := result.RowsAffected()
	if err != nil {
		return seq, errors.Wrap(err, "Failed to receive From Database")
	}
	defer stmt.Close()

	seq = int(row)

	return seq, nil
}

func GetAllUserList(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		trainer_id, group_name := common.GetQueryString(c)
		user := UserList{}
		users, err := user.selectAllUser(db, trainer_id, group_name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, users))
		} else {
			c.JSON(http.StatusOK, common.SucceedResponse(users))
		}
	}

	return resultFunc
}

func GetManageDetailUserList(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		user_seq := c.Param("user_seq")
		user_info, err := selectDetailUser(db, user_seq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, user_info))
		} else {
			c.JSON(http.StatusOK, common.SucceedResponse(user_info))
		}
	}

	return resultFunc
}

func GetManageUserList(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		trainer_id, group_name := common.GetQueryString(c)
		users, err := selectManageUser(db, trainer_id, group_name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, users))
		} else {
			c.JSON(http.StatusOK, common.SucceedResponse(users))
		}
	}

	return resultFunc
}

func PostUser(db *sql.DB) gin.HandlerFunc {
	resultFunc := func(c *gin.Context) {
		user := User{}
		if err := c.Bind(&user); err != nil {
			c.JSON(http.StatusBadRequest, common.FailedResponse(err, user))
		}

		_, err := user.addUser(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, user))
		} else {
			c.JSON(http.StatusCreated, common.SucceedResponse(user))
		}
	}

	return resultFunc
}

func DeleteUser(db *sql.DB) gin.HandlerFunc {
	resultFUnc := func(c *gin.Context) {
		user_seq := c.Param("user_seq")
		seq, err := deleteUser(db, user_seq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.FailedResponse(err, seq))
		} else {
			c.JSON(http.StatusCreated, common.SucceedResponse(seq))
		}
	}

	return resultFUnc
}

func UserListRouter(router *gin.Engine, db *sql.DB) {
	user := router.Group("/api/user")

	user.GET("/all", GetAllUserList(db))
	user.GET("/management", GetManageUserList(db))
	user.GET("/detail/:user_seq", GetManageDetailUserList(db))

	user.POST("/", PostUser(db))

	user.DELETE("/:user_seq", DeleteUser(db))
}
