package users

import (
	"net/http"
	"strconv"

	"github.com/fmarinCeiba/bookstore_oauth-go/oauth"
	"github.com/fmarinCeiba/bookstore_users-api/domain/users"
	"github.com/fmarinCeiba/bookstore_users-api/services"
	"github.com/fmarinCeiba/bookstore_utils-go/rest_errors"
	"github.com/gin-gonic/gin"
)

func TestServiceInterface() {

}

func getUserID(uIDParam string) (int64, rest_errors.RestErr) {
	uID, uErr := strconv.ParseInt(uIDParam, 10, 64)
	if uErr != nil {
		return 0, rest_errors.NewBadRequestError("invalid user id")
	}
	return uID, nil
}

func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	// if callerID := oauth.GetCallerID(c.Request); callerID == 0 {
	// 	err := errors.RestErr{
	// 		Status:  http.StatusUnauthorized,
	// 		Message: "resource not available",
	// 	}
	// 	c.JSON(err.Status, err)
	// 	return
	// }

	uID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	result, getErr := services.UserService.Get(uID)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}
	if oauth.GetCallerID(c.Request) == result.Id {
		c.JSON(http.StatusOK, result.Marshall(false))
		return
	}
	c.JSON(http.StatusOK, result.Marshall(oauth.IsPublic(c.Request)))
}

func Search(c *gin.Context) {
	status := c.Query("status")

	users, err := services.UserService.Search(status)
	if err != nil {
		c.JSON(err.Status(), err)
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
		restErr := rest_errors.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	result, savErr := services.UserService.Create(user)
	if savErr != nil {
		c.JSON(savErr.Status(), savErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Update(c *gin.Context) {
	uID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	isPartial := c.Request.Method == http.MethodPatch
	user.Id = uID
	result, updErr := services.UserService.Update(isPartial, user)
	if updErr != nil {
		c.JSON(updErr.Status(), updErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context) {
	uID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}
	if delErr := services.UserService.Delete(uID); delErr != nil {
		c.JSON(delErr.Status(), delErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func LogIn(c *gin.Context) {
	var req users.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		restErr := rest_errors.NewBadRequestError(err.Error())
		c.JSON(restErr.Status(), restErr)
		return
	}
	user, err := services.UserService.LogIn(req)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}
