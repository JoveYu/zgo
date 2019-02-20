package scanner

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/JoveYu/zgo/log"
	"github.com/JoveYu/zgo/sql"
)

type Test struct {
	Id   int       `zdb:"id"`
	Name string    `zdb:"name"`
	Time time.Time `zdb:"time"`
}

func TestAll(t *testing.T) {
	log.Install("stdout")
	sql.Install(sql.DBConf{
		"sqlite3": []string{"sqlite3", "file::memory:?mode=memory&cache=shared"},
	})
	db := sql.GetDB("sqlite3")
	db.Exec("drop table if exists test")
	db.Exec("create table if not exists test(id integer not null primary key, name text, time datetime)")
	for i := 1; i <= 3; i++ {
		db.Insert("test", sql.Values{
			"id":   i,
			"name": fmt.Sprintf("name %d", i),
			"time": time.Now(),
		})
	}
	rows, _ := db.Select("test", sql.Where{"id >": 2})

	test := []Test{}
	log.Debug(test)
	err := ScanStruct(rows, &test)
	log.Debug(err)
	log.Debug(test)

	rows, _ = db.Select("test", sql.Where{})

	test2 := Test{}
	log.Debug(test2)
	err = ScanStruct(rows, &test2)
	log.Debug(err)
	log.Debug(test2)

}
