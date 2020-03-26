package users

import (
	"net/http"
	"strconv"

	"github.com/fmarinCeiba/bookstore_users-api/domain/users"
	"github.com/fmarinCeiba/bookstore_users-api/services"
	"github.com/fmarinCeiba/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
)

func getUserId(uIdParam string) (int64, *errors.RestErr) {
	uId, uErr := strconv.ParseInt(uIdParam, 10, 64)
	if uErr != nil {
		return 0, errors.NewBadRequestError("invalid user id")
	}
	return uId, nil
}

func Get(c *gin.Context) {
	uId, uErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if uErr != nil {
		err := errors.NewBadRequestError("invalid user id")
		c.JSON(err.Status, err)
		return
	}

	result, getErr := services.GetUser(uId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func Search(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")
}

func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	result, savErr := services.CreateUser(user)
	if savErr != nil {
		c.JSON(savErr.Status, savErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func Update(c *gin.Context) {
	uId, uErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if uErr != nil {
		err := errors.NewBadRequestError("invalid user id")
		c.JSON(err.Status, err)
		return
	}
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	isPartial := c.Request.Method == http.MethodPatch
	user.Id = uId
	result, updErr := services.UpdateUser(isPartial, user)
	if updErr != nil {
		c.JSON(updErr.Status, updErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func Delete(c *gin.Context) {

}
