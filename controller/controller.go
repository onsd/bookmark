package controller

import (
	"net/http"

	"github.com/labstack/echo"
)

func auth(c echo.Context) error {
	dbmap := initDB()
	defer dbmap.Db.Close()

	username := c.FormValue("username")
	password := c.FormValue("password")

	count, err := dbmap.SelectInt("select count(*) from users where username=?", username)
	checkErr(err, "select count(*) failed")
	user := NewUser(username, password)
	if count != 0 {
		//already registerd or same username
		res := &Response{
			Status:  "Conflict or just Error",
			Content: user,
		}
		return c.JSON(http.StatusConflict, res)
	}
	//not registerd

	err = dbmap.Insert(&user)
	checkErr(err, "insert new user failed")
	res := &Response{
		Status:  "Successed authorize new user",
		Content: user,
	}
	return c.JSON(http.StatusOK, res)
}
