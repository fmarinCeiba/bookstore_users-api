package services

import (
	"github.com/fmarinCeiba/bookstore_users-api/domain/users"
	"github.com/fmarinCeiba/bookstore_utils-go/rest_errors"
)

var (
	UserService usersServiceInterface = &usersService{}
)

type usersService struct {
}

type usersServiceInterface interface {
	Get(int64) (*users.User, rest_errors.RestErr)
	Search(string) (users.Users, rest_errors.RestErr)
	Create(users.User) (*users.User, rest_errors.RestErr)
	Update(bool, users.User) (*users.User, rest_errors.RestErr)
	Delete(int64) rest_errors.RestErr
	LogIn(users.LoginRequest) (*users.User, rest_errors.RestErr)
}

func (s *usersService) Get(uID int64) (*users.User, rest_errors.RestErr) {
	u := users.User{Id: uID}
	if err := u.Get(); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *usersService) Search(status string) (users.Users, rest_errors.RestErr) {
	dao := users.User{}
	return dao.FindByStatus(status)
}

func (s *usersService) Create(u users.User) (*users.User, rest_errors.RestErr) {
	if err := u.Validate(); err != nil {
		return nil, err
	}

	if err := u.Save(); err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *usersService) Update(isPartial bool, u users.User) (*users.User, rest_errors.RestErr) {
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

func (s *usersService) Delete(uID int64) rest_errors.RestErr {
	u := users.User{Id: uID}
	return u.Delete()
}

func (s *usersService) LogIn(lr users.LoginRequest) (*users.User, rest_errors.RestErr) {
	dao := &users.User{
		Email:    lr.Email,
		Password: lr.Password,
	}
	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}

	return dao, nil
}
