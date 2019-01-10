package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	gorp "gopkg.in/gorp.v1"
)

func main() {
	// initialize the DbMap
	dbmap := initDb()
	defer dbmap.Db.Close()

	// delete any existing rows
	//err := dbmap.TruncateTables()
	//checkErr(err, "TruncateTables failed")

	//bookmark1 := NewBookmark("google.com", "Useful for searching")
	bookmark1 := &Bookmark{1, "yahoo", "hello"}
	bookmark2 := &Bookmark{2, "google", "useful"}
	err := dbmap.Insert(bookmark1, bookmark2)
	if err != nil {
		log.Fatal(err)
	}

}

// Bookmark represents one bookmark.
type Bookmark struct {
	ID          int64  `db:"bookmark_id, primarykey, autoincrement`
	URL         string `db:"url, not null"`
	Description string `db:"description"`
	//Created     int64  `db:"created"`
}

// User represents one user.
type User struct {
	ID       int64  `db: "user_id, primarykey, autoincrement`
	Username string `db: "username`
	Password string `db: "password"`
}

// NewBookmark returns the Bookmark struct with URL, Description and present time
func NewBookmark(URL, Description string) Bookmark {
	return Bookmark{
		ID:          11,
		URL:         URL,
		Description: Description,
		Created:     time.Now(),
	}
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("sqlite3", "./test.db")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	//dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "ID")
	dbmap.AddTableWithName(Bookmark{}, "bookmarks").SetKeys(true, "ID")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
