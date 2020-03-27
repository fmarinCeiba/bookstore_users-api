package users

import (
	"fmt"

	"github.com/fmarinCeiba/bookstore_users-api/datasources/mysql/users_db"
	"github.com/fmarinCeiba/bookstore_users-api/utils/crypto_utils"
	"github.com/fmarinCeiba/bookstore_users-api/utils/date_utils"
	"github.com/fmarinCeiba/bookstore_users-api/utils/errors"
	"github.com/fmarinCeiba/bookstore_users-api/utils/mysql_utils"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);"
	queryGetAnUser        = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?"
	queryUpdateUser       = "UPDATE users SET first_name = ?, last_name = ?, email = ?, status = ? WHERE id = ?"
	queryDeleteUser       = "DELETE FROM users WHERE id = ?"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = ?"
)

var (
	usersDB = make(map[int64]*User)
)

func (u *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetAnUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	iResult := stmt.QueryRow(u.Id)
	if err := iResult.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status); err != nil {
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

	u.Status = StatusActive
	u.DateCreated = date_utils.GetNowDBFormat()
	u.Password = crypto_utils.GetMd5(u.Password)
	iResult, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.DateCreated, u.Status, u.Password)
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

func (u *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(u.FirstName, u.LastName, u.Email, u.Id); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (u *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(u.Id); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (u *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}
