
package sql

import "fmt"
import "time"
import "testing"
import "github.com/JoveYu/zgo/log"
import _ "github.com/mattn/go-sqlite3"
import _ "github.com/go-sql-driver/mysql"

func TestInstall(t *testing.T) {
    log.Install("stdout")
    Install(map[string][]string{
        "sqlite3": []string{"sqlite3", "file::memory:?mode=memory&cache=shared"},
        "mysql": []string{"mysql", "test:123456@tcp(127.0.0.1:3306)/cmdb?charset=utf8mb4"},
    })
    db := GetDB("sqlite3")

    db.Exec("drop table if exists test")
    db.Exec("create table if not exists test(id integer not null primary key, name text, time datetime)")

    var id int64
    for i:=1; i<=10; i++ {
        id = db.Insert("test", map[string]interface{}{
            "id": i,
            "name": fmt.Sprintf("name %d", i),
            "time": time.Now(),
        })
        log.Debug("insert get id %d", id)
    }

    rows := db.Select("test", map[string]interface{}{
        "_field": "count(*)",
    })
    log.Debug("select count: %s", rows)

    rows = db.Select("test", map[string]interface{}{
        "id in": []int{2,3},
    })
    log.Debug("select in: %s", rows)

    rows = db.Select("test", map[string]interface{}{
        "id between": []int{2,5},
        "_other": "order by id desc",
    })
    log.Debug("select between: %s", rows)

    id = db.Delete("test", map[string]interface{}{
        "id >": 5,
    })
    log.Debug("delete ret: %d", id)
    rows = db.Select("test", map[string]interface{}{
        "_field": "count(*)",
    })
    log.Debug("select count: %s", rows)

    db.Update("test", map[string]interface{}{
        "name": "new name",
    }, map[string]interface{}{
        "id <": 3,
    })
    rows = db.Select("test", map[string]interface{}{})
    log.Debug("select update: %s", rows)

    db.Update("test", map[string]interface{}{
        "name": "new name",
    }, map[string]interface{}{
    })
}

func TestTransaction(t *testing.T) {
    log.Install("stdout")
    Install(map[string][]string{
        "sqlite3": []string{"sqlite3", "file::memory:?mode=memory&cache=shared"},
        "mysql": []string{"mysql", "test:123456@tcp(127.0.0.1:3306)/cmdb?charset=utf8mb4"},
    })
    db := GetDB("sqlite3")

    db.Exec("drop table if exists test")
    db.Exec("create table if not exists test(id integer not null primary key, name text, time datetime)")

    tx,_ := db.Begin()
    tx.Exec("select 1")
    tx.Rollback()

}
