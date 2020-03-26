package users

import (
	"github.com/fmarinCeiba/bookstore_users-api/datasources/mysql/users_db"
	"github.com/fmarinCeiba/bookstore_users-api/utils/date_utils"
	"github.com/fmarinCeiba/bookstore_users-api/utils/errors"
	"github.com/fmarinCeiba/bookstore_users-api/utils/mysql_utils"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	errorNoRows      = "no rows in result set"
	queryInsertUser  = "INSERT INTO users(first_name, last_name, email, date_created) VALUES (?, ?, ?, ?);"
	queryGetAnUser   = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id = ?"
)

var (
	usersDB = make(map[int64]*User)
)

func (u *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetAnUser)
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	defer stmt.Close()
	iResult := stmt.QueryRow(u.Id)
	if err := iResult.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (u *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	u.DateCreated = date_utils.GetNowString()
	iResult, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.DateCreated)
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	uId, err := iResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	u.Id = uId

	return nil
}
