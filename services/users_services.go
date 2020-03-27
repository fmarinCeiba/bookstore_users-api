package services

import (
	"github.com/fmarinCeiba/bookstore_users-api/domain/users"
	"github.com/fmarinCeiba/bookstore_users-api/utils/errors"
)

var (
	UserService usersServiceInterface = &usersService{}
)

type usersService struct {
}

type usersServiceInterface interface {
	Get(int64) (*users.User, *errors.RestErr)
	Search(string) (users.Users, *errors.RestErr)
	Create(users.User) (*users.User, *errors.RestErr)
	Update(bool, users.User) (*users.User, *errors.RestErr)
	Delete(int64) *errors.RestErr
}

func (s *usersService) Get(uID int64) (*users.User, *errors.RestErr) {
	u := users.User{Id: uID}
	if err := u.Get(); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *usersService) Search(status string) (users.Users, *errors.RestErr) {
	dao := users.User{}
	return dao.FindByStatus(status)
}

func (s *usersService) Create(u users.User) (*users.User, *errors.RestErr) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := u.Save(); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *usersService) Update(isPartial bool, u users.User) (*users.User, *errors.RestErr) {
	c := &users.User{Id: u.Id}
	if err := c.Get(); err != nil {
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

func (s *usersService) Delete(uID int64) *errors.RestErr {
	u := users.User{Id: uID}
	return u.Delete()
}
