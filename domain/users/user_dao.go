package users

import (
	"fmt"
	"strings"

	"github.com/fmarinCeiba/bookstore_users-api/datasources/mysql/users_db"
	"github.com/fmarinCeiba/bookstore_users-api/utils/date_utils"
	"github.com/fmarinCeiba/bookstore_users-api/utils/errors"
	"github.com/go-sql-driver/mysql"
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
	// if err := users_db.Client.Ping(); err != nil {
	// 	panic(err)
	// }
	stmt, err := users_db.Client.Prepare(queryGetAnUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	iResult := stmt.QueryRow(u.Id)
	if err := iResult.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated); err != nil {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("user %d not found", u.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to get user %d: %s", u.Id, err.Error()))
	}
	// LOCAL INSTANT DB
	// result := usersDB[u.Id]
	// if result == nil {
	// 	return errors.NewNotFoundError(fmt.Sprintf("user %d not found", u.Id))
	// }

	// u.Id = result.Id
	// u.FirstName = result.FirstName
	// u.LastName = result.LastName
	// u.Email = result.Email
	// u.DateCreated = result.DateCreated
	return nil
}

func (u *User) Save() *errors.RestErr {
	// if err := users_db.Client.Ping(); err != nil {
	// 	panic(err)
	// }
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	u.DateCreated = date_utils.GetNowString()
	iResult, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.DateCreated)
	if err != nil {
		sqlErr, ok := err.(*mysql.MySQLError)
		if !ok {
			return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
		}
		fmt.Println(sqlErr)
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", u.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}
	uId, err := iResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}
	u.Id = uId

	// LOCAL INSTANT DB
	// current := usersDB[u.Id]
	// if current != nil {
	// 	if current.Email == u.Email {
	// 		return errors.NewBadRequestError(fmt.Sprintf("user %s already registered", u.Email))
	// 	}
	// 	return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", u.Id))
	// }
	// usersDB[u.Id] = u

	return nil
}
