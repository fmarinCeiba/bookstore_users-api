package users

import (
	"net/http"
	"strconv"

	"github.com/fmarinCeiba/bookstore_users-api/domain/users"
	"github.com/fmarinCeiba/bookstore_users-api/services"
	"github.com/fmarinCeiba/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
)

func TestServiceInterface() {

}

func getUserID(uIDParam string) (int64, *errors.RestErr) {
	uID, uErr := strconv.ParseInt(uIDParam, 10, 64)
	if uErr != nil {
		return 0, errors.NewBadRequestError("invalid user id")
	}
	return uID, nil
}

func Get(c *gin.Context) {
	uID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	result, getErr := services.UserService.Get(uID)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Search(c *gin.Context) {
	status := c.Query("status")

	users, err := services.UserService.Search(status)
	if err != nil {
		c.JSON(err.Status, err)
	}
	// result := make([]interface{}, len(users))
	// for index, user := range users {
	// 	result[index] = user.Marshall(c.GetHeader("X-Public") == "true")
	// }
	c.JSON(http.StatusOK, users.Marshal(c.GetHeader("X-Public") == "true"))
}

func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	result, savErr := services.UserService.Create(user)
	if savErr != nil {
		c.JSON(savErr.Status, savErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Update(c *gin.Context) {
	uID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	isPartial := c.Request.Method == http.MethodPatch
	user.Id = uID
	result, updErr := services.UserService.Update(isPartial, user)
	if updErr != nil {
		c.JSON(updErr.Status, updErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	uID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}
	if delErr := services.UserService.Delete(uID); delErr != nil {
		c.JSON(delErr.Status, delErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}
