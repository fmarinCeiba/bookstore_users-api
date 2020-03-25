package services

import (
	"github.com/fmarinCeiba/bookstore_users-api/domain/users"
	"github.com/fmarinCeiba/bookstore_users-api/utils/errors"
)

func GetUser(uId int64) (*users.User, *errors.RestErr) {
	u := users.User{Id: uId}
	if err := u.Get(); err != nil {
		return nil, err
	}

	return &u, nil
}

func CreateUser(u users.User) (*users.User, *errors.RestErr) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := u.Save(); err != nil {
		return nil, err
	}

	return &u, nil
}
