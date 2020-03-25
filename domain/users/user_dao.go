package users

import (
	"fmt"

	"github.com/fmarinCeiba/bookstore_users-api/utils/errors"
)

var (
	usersDB = make(map[int64]*User)
)

func (u *User) Get() *errors.RestErr {
	result := usersDB[u.Id]
	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", u.Id))
	}
	u.Id = result.Id
	u.FirstName = result.FirstName
	u.LastName = result.LastName
	u.Email = result.Email
	u.DateCreated = result.DateCreated
	return nil
}

func (u *User) Save() *errors.RestErr {
	current := usersDB[u.Id]
	if current != nil {
		if current.Email == u.Email {
			return errors.NewBadRequestError(fmt.Sprintf("user %s already registered", u.Email))
		}
		return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", u.Id))
	}
	usersDB[u.Id] = u
	return nil
}
