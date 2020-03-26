package mysql_utils

import (
	"fmt"
	"strings"

	"github.com/fmarinCeiba/bookstore_users-api/utils/errors"
	"github.com/go-sql-driver/mysql"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError("no record matching given id")
		}
		return errors.NewInternalServerError("error parsing database response")
	}
	fmt.Println(sqlErr)
	switch sqlErr.Number {
	case 1062:
		return errors.NewBadRequestError("duplicated key data")
	}
	return errors.NewInternalServerError("error processing request")
}
