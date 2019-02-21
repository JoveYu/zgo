package builder

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/JoveYu/zgo/log"
)

func TestAll(t *testing.T) {
	log.Install("stdout")

	db, _ := sql.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	db.Exec("drop table if exists test")
	db.Exec("create table if not exists test(id integer not null primary key, name text, time datetime)")
	for i := 1; i <= 10; i++ {
		sql, args := Insert("test", Values{
			"id":   i,
			"name": fmt.Sprintf("name %d", i),
			"time": time.Now(),
		})
		log.Debug("sql: %s, args: %v", sql, args)
		log.Debug("sql: %s", FormatSql(sql, args...))
		_, err := db.Exec(sql, args...)
		log.Debug("insert err:%v", err)
	}

	sql, args := Select("test", Where{
		"id":           1,
		"id >":         0,
		"id is":        nil,
		"id not in":    []int{3, 4},
		"name between": []string{"name 1", "name 5"},
		"_field":       "count(*)",
		"_groupby":     "name",
		"_having": Where{
			"id >": 0,
		},
		"_other": "limit 1",
	})

	log.Debug("sql: %s, args: %v", sql, args)
	log.Debug("sql: %s", FormatSql(sql, args...))

	_, err := db.Query(sql, args...)
	log.Debug("select err:%v", err)

	sql, args = Update("test", Values{
		"name": "new name",
	}, Where{
		"id >": 3,
	})
	log.Debug("sql: %s, args: %v", sql, args)
	log.Debug("sql: %s", FormatSql(sql, args...))
	_, err = db.Exec(sql, args...)
	log.Debug("update err:%v", err)

	sql, args = Delete("test", Where{
		"id !=": 3,
	})
	log.Debug("sql: %s, args: %v", sql, args)
	log.Debug("sql: %s", FormatSql(sql, args...))
	_, err = db.Exec(sql, args...)
	log.Debug("delete err:%v", err)

	sql, args = Select("test", Where{
		"id >":  0,
		"id > ": 1,
	})
	log.Debug("sql: %s, args: %v", sql, args)

}
