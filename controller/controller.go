package controller

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	//using mattn/go-sqlite3 with blank import
	_ "github.com/mattn/go-sqlite3"
	gorp "gopkg.in/gorp.v1"

	"github.com/labstack/echo"
)

// Bookmark represents one bookmark.
type Bookmark struct {
	BookmarkID  int64  `db:"bookmark_id, primarykey, autoincrement`
	URL         string `db:"url, not null"`
	Description string `db:"description"`
	Created     int64  `db:"created"`
}

// User represents one user.
type User struct {
	UserID   int64  `db: "user_id, primarykey, autoincrement`
	Username string `db: "username`
	Password string `db: "password"`
}

type Response struct {
	Status  string
	Content interface{}
}

func initDB() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("sqlite3", "./test.db")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	//dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "ID")
	dbmap.AddTableWithName(Bookmark{}, "bookmarks").SetKeys(true, "BookmarkID")
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "UserID")
	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}
func TestDB() error {
	dbmap := initDB()
	defer dbmap.Db.Close()
	// delete any existing rows
	err := dbmap.TruncateTables()
	checkErr(err, "TruncateTables failed")

	bookmark1 := &Bookmark{1, "yahoo", "hello", time.Now().UnixNano()}
	bookmark2 := &Bookmark{1, "google", "useful", time.Now().UnixNano()}
	err = dbmap.Insert(bookmark1, bookmark2)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
func Accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}
func Auth(c echo.Context) error {
	dbmap := initDB()
	defer dbmap.Db.Close()

	username := c.FormValue("username")
	password := c.FormValue("password")

	count, err := dbmap.SelectInt("select count(*) from users where username=?", username)
	checkErr(err, "select count(*) failed")
	user := newUser(username, password)
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

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

//NewUser returns the User struct with username and password.
func newUser(username, password string) User {
	return User{
		UserID:   1,
		Username: username,
		Password: password,
	}
}

func Login(c echo.Context) error {
	dbmap := initDB()
	defer dbmap.Db.Close()

	username := c.FormValue("username")
	password := c.FormValue("password")

	count, err := dbmap.SelectInt("select count(*) from users where username=? and password=?", username, password)
	checkErr(err, "select count(*) failed")
	user := newUser(username, password)
	if count == 0 {
		//already registerd or same username
		res := &Response{
			Status:  "Not registrated",
			Content: user,
		}
		return c.JSON(http.StatusUnauthorized, res)
	}

	//Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = username
	claims["password"] = password
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte(os.Getenv("SECRETKEY")))
	if err != nil {
		checkErr(err, "token sign error")
	}
	res := &Response{
		Status:  "authorized",
		Content: t,
	}
	return c.JSON(http.StatusOK, res)
}

func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	password := claims["password"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!"+"\nYour password is "+password)
}
