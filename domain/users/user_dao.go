package users

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fmarinCeiba/bookstore_users-api/datasources/mysql/users_db"
	"github.com/fmarinCeiba/bookstore_users-api/logger"
	"github.com/fmarinCeiba/bookstore_users-api/utils/crypto_utils"
	"github.com/fmarinCeiba/bookstore_users-api/utils/date_utils"
	"github.com/fmarinCeiba/bookstore_users-api/utils/mysql_utils"
	"github.com/fmarinCeiba/bookstore_utils-go/rest_errors"
)

const (
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);"
	queryGetAnUser              = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?;"
	queryUpdateUser             = "UPDATE users SET first_name = ?, last_name = ?, email = ?, status = ? WHERE id = ?;"
	queryDeleteUser             = "DELETE FROM users WHERE id = ?;"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = ?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email = ? AND password = ? AND status = ?"
)

var (
	usersDB = make(map[int64]*User)
)

func (u *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetAnUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("error when trying to prepare get user statement", errors.New("database error"))
	}
	defer stmt.Close()
	iResult := stmt.QueryRow(u.Id)
	if err := iResult.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); err != nil {
		logger.Error("error when trying to get user by id", err)
		return rest_errors.NewInternalServerError("error when trying to prepare get user statement", errors.New("database error"))
	}
	return nil
}

func (u *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("error when trying to prepare save user statement", errors.New("database error"))
	}
	defer stmt.Close()

	u.Status = StatusActive
	u.DateCreated = date_utils.GetNowDBFormat()
	u.Password = crypto_utils.GetMd5(u.Password)
	iResult, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.DateCreated, u.Status, u.Password)
	if err != nil {
		logger.Error("error when trying to save user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", errors.New("database error"))
	}
	uID, err := iResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to get last insert id after creating a new user", err)
		return rest_errors.NewInternalServerError("error when trying to get last insert id after creating a new user", errors.New("database error"))
	}
	u.Id = uID

	return nil
}

func (u *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("error when trying to prepare update user statement", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(u.FirstName, u.LastName, u.Email, u.Id, u.Status); err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("error when trying to update user", errors.New("database error"))
	}
	return nil
}

func (u *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("error when trying to prepare delete user statement", errors.New("database error"))
	}
	defer stmt.Close()

	if _, err = stmt.Exec(u.Id); err != nil {
		logger.Error("error when trying to delete user", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", errors.New("database error"))
	}
	return nil
}

func (u *User) FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to prepare find users by status statement", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find users by status", errors.New("database error"))
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to scan row into user struct", err)
			return nil, rest_errors.NewInternalServerError("error when trying to scan row into user struct", errors.New("database error"))
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}

func (u *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("error when trying to prepare get user by email and password statement", errors.New("database error"))
	}
	defer stmt.Close()
	iResult := stmt.QueryRow(u.Email, crypto_utils.GetMd5(u.Password), StatusActive)
	if err := iResult.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); err != nil {
		if strings.Contains(err.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", err)
		return rest_errors.NewInternalServerError("error when trying to get user by email and password", errors.New("database error"))
	}
	return nil
}
