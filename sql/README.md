
## zgo/sql

整体思路和使用习惯与 `github.com/JoveYu/zpy/base/dbpool.py` 一致，主要根据go语言静态强类型，以及不能使用可选参数的特性，调整了下使用方式

整体使用了一段时间还是比较好用的，符合我的风格

1. 方便的初始化数据库连接池
2. 方便利用结构化数据组装SQL
3. 提供方便的Scan，可以直接查询结果到struct
4. 统一日志输出，打印连接池状态

## TODO

1. sqlbuilder 支持selectjoin 比较容易 现在没需求

## Example

```go
package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"

	"github.com/JoveYu/zgo/log"
	"github.com/JoveYu/zgo/sql"
)

func main() {
	log.Install("stdout")

	sql.Install(sql.DBConf{
		"testdb": []string{"sqlite3", "file::memory:?mode=memory&cache=shared"},
	})
	db := sql.GetDB("testdb")

	db.Exec("drop table if exists test")
	db.Exec("create table if not exists test(id integer not null primary key, name text, time datetime)")

	// INSERT INTO `test`(`id`,`name`,`time`) VALUES(1,'name 1','2019-02-20 16:20:07')
	for i := 1; i <= 10; i++ {
		db.Insert("test", sql.Values{
			"id":   i,
			"name": fmt.Sprintf("name %d", i),
			"time": time.Now(),
		})
	}

	// SELECT * FROM `test` WHERE (`id` = 2) order by id desc limit 1
	db.Select("test", sql.Where{
		"id":     2,
		"_other": "order by id desc limit 1",
	})

	// SELECT count(1) FROM `test` WHERE (`id` > 2) and (`id` between 4 and 6) and (`name` in ('foo','bar'))
	db.Select("test", sql.Where{
		"id >":       2,
		"id between": []int{4, 6},
		"name in":    []string{"foo", "bar"},
		"_field":     "count(1)",
	})

	// SELECT * FROM `test` WHERE (`id` > 2) GROUP BY name HAVING (`id` > 3)
	db.Select("test", sql.Where{
		"id >":     2,
		"_groupby": "name",
		"_having": sql.Where{
			"id >": 3,
		},
	})

	// UPDATE `test` SET `id`=-1 WHERE (`id` > 9)
	db.Update("test", sql.Values{
		"id": -1,
	}, sql.Where{
		"id >": 9,
	})

	// UPDATE `test` SET `name`='jove'
	db.Update("test", sql.Values{
		"name": "jove",
	}, sql.Where{})

	// DELETE FROM `test` WHERE (`name` != 'jove')
	db.Delete("test", sql.Where{
		"name !=": "jove",
	})

	// select scan to struct
	type User struct {
		Id   int       `zdb:"id"`
		Name string    `zdb:"name"`
		Time time.Time `zdb:"time"`
	}
	user := []User{}
	db.SelectScan(&user, "test", sql.Where{})
	log.Debug(user)
	// [{Id:-1 Name:jove Time:2019-02-20 17:20:25.03967 +0800 +0800}......

}

```

