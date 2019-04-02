package sql

import "fmt"
import "context"
import "sync"
import "time"
import "testing"
import "github.com/JoveYu/zgo/log"
import _ "github.com/mattn/go-sqlite3"
import _ "github.com/go-sql-driver/mysql"

func TestInstall(t *testing.T) {
	log.Install("stdout")
	Install(DBConf{
		"sqlite3": []string{"sqlite3", "file::memory:?mode=memory&cache=shared"},
		// "mysql": []string{"mysql", "test:123456@tcp(127.0.0.1:3306)/zgo?charset=utf8mb4"},
	})
	db := GetDB("sqlite3")

	db.Exec("wrong sql test")

	db.Exec("drop table if exists test")
	db.Exec("create table if not exists test(id integer not null primary key, name text, time datetime)")

	for i := 1; i <= 10; i++ {
		db.InsertContext(context.TODO(), "test", Values{
			"id":   i,
			"name": fmt.Sprintf("name %d", i),
			"time": time.Now(),
		})
	}

	rows, err := db.SelectMap("test", Where{
		"_field": "count(*)",
	})
	log.Debug("select count: %s, err: %s", rows, err)

	rows, err = db.SelectMap("test", Where{
		"id in": []int{2, 3},
	})
	log.Debug("select in: %s", rows)

	rows, err = db.SelectMap("test", Where{
		"id between": []int{2, 5},
		"_other":     "order by id desc",
	})
	log.Debug("select between: %s", rows)

	db.Delete("test", Where{
		"id >": 5,
	})

	rows, err = db.SelectMap("test", Where{
		"_field": "count(*)",
	})
	log.Debug("select count: %s", rows)

	db.Update("test", Values{
		"name": "new name",
	}, Where{
		"id <": 3,
	})
	rows, err = db.SelectMap("test", Where{})
	log.Debug("select update: %s", rows)

	db.Update("test", Values{
		"name": "new name",
	}, Where{})

	rows, err = db.SelectMap("test", Where{
		"name": "??",
	})
	log.Debug("select ? %s", rows)

}

func TestTransaction(t *testing.T) {
	log.Install("stdout")
	Install(DBConf{
		"sqlite3": []string{"sqlite3", "file::memory:?mode=memory&cache=shared"},
		"mysql":   []string{"mysql", "test:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4"},
	})
	db := GetDB("sqlite3")

	db.Exec("drop table if exists test")
	db.Exec("create table if not exists test(id integer not null primary key, name text, time datetime)")

	tx, _ := db.Begin()
	tx.Insert("test", Values{
		"id":   1,
		"name": "name",
		"time": time.Now(),
	})

	rows, err := db.SelectMap("test", Where{})
	log.Debug("%s %s", rows, err)

	tx.Commit()

	rows, err = db.SelectMap("test", Where{})
	log.Debug("%s", rows)

}

func TestMulitRun(t *testing.T) {
	log.Install("stdout")
	Install(DBConf{
		"sqlite3": []string{"sqlite3", "file::memory:?mode=memory&cache=shared"},
	})
	db := GetDB("sqlite3")
	db.Exec("drop table if exists test")
	db.Exec("create table if not exists test(id integer not null primary key, name text, time datetime)")

	count := 3

	for i := 0; i < count; i++ {
		db.Insert("test", Values{
			"id":   i,
			"name": fmt.Sprintf("name %d", i),
			"time": time.Now(),
		})
	}
	var wa sync.WaitGroup
	wa.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			db.SelectMap("test", Where{})
			wa.Done()
		}()
	}
	wa.Wait()
}

type User struct {
	Id    int       `zdb:"id"`
	Name  string    `zdb:"name"`
	Time  time.Time `zdb:"time"`
	Other string
}

func TestScan(t *testing.T) {
	log.Install("stdout")
	Install(DBConf{
		"sqlite3": []string{"sqlite3", "file::memory:?mode=memory&cache=shared"},
	})
	db := GetDB("sqlite3")
	db.Exec("drop table if exists test")
	db.Exec("create table if not exists test(id integer not null primary key, name text, time datetime)")

	count := 3

	for i := 0; i < count; i++ {
		db.Insert("test", Values{
			"id":   i,
			"name": fmt.Sprintf("name %d", i),
			"time": time.Now(),
		})
	}

	user := []User{}
	log.Debug(user)
	err := db.SelectScan(&user, "test", Where{})
	log.Debug(err)
	log.Debug(user)

	user1 := User{}
	log.Debug(user1)
	err = db.SelectScan(&user1, "test", Where{"_other": "order by id desc"})
	log.Debug(err)
	log.Debug(user1)

}
