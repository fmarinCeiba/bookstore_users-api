package users

import (
	"net/http"
	"strconv"

	"github.com/fmarinCeiba/bookstore_users-api/domain/users"
	"github.com/fmarinCeiba/bookstore_users-api/services"
	"github.com/fmarinCeiba/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
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

func SearchUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "implement me!")
}

func CreateUser(c *gin.Context) {
	var user users.User
	// fmt.Println(user)
	// bytes, err := ioutil.ReadAll(c.Request.Body)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// if err := json.Unmarshal(bytes, &user); err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError(err.Error())
		c.JSON(restErr.Status, restErr)
		return
	}
	// fmt.Println(user)
	result, savErr := services.CreateUser(user)
	if savErr != nil {
		c.JSON(savErr.Status, savErr)
		return
	}
	c.JSON(http.StatusCreated, result)
}
