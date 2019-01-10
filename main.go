package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	gorp "gopkg.in/gorp.v1"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

func main() {
	// initialize echo instance
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Login route
	//e.POST("/login", login)
	//Authorized route
	e.POST("/auth", auth)

	//Unauthenticated route
	e.GET("/", accessible)

	// initialize the DbMap
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
	e.Start(":8080")
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

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

// NewBookmark returns the Bookmark struct with URL, Description and present time
func NewBookmark(URL, Description string) Bookmark {
	return Bookmark{
		URL:         URL,
		Description: Description,
		Created:     time.Now().UnixNano(),
	}
}

//NewUser returns the User struct with username and password.
func NewUser(username, password string) User {
	return User{
		UserID:   1,
		Username: username,
		Password: password,
	}
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
