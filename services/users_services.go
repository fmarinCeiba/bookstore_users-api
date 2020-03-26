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

func UpdateUser(isPartial bool, u users.User) (*users.User, *errors.RestErr) {
	c, err := GetUser(u.Id)
	if err != nil {
		return nil, err
	}

	if isPartial {
		if u.FirstName != "" {
			c.FirstName = u.FirstName
		}
		if u.LastName != "" {
			c.LastName = u.LastName
		}
		if u.Email != "" {
			c.Email = u.Email
		}
	} else {
		if err := u.Validate(); err != nil {
			return nil, err
		}
		c.FirstName = u.FirstName
		c.LastName = u.LastName
		c.Email = u.Email
	}

	if err := u.Update(); err != nil {
		return nil, err
	}
	return c, nil
}
