// use sql not orm
// use simple sql not join

package sql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/JoveYu/zgo/log"
	"github.com/JoveYu/zgo/sql/builder"
)

var (
	dbMap = make(map[string]DB)
)

type Base struct {
	db    DB
	tx    Tx
	trans bool
}

type DB struct {
	*Base
	*sql.DB
	name   string
	driver string
	dsn    string
}

type Tx struct {
	*Base
	*sql.Tx
	db *DB
}

type Where builder.Where
type Values builder.Values
type DBConf map[string][]string

func Install(conf DBConf) map[string]DB {
	log.Debug("available sql driver: %s", sql.Drivers())
	for k, v := range conf {
		if len(v) != 2 {
			log.Fatal("parse db config error")
		}
		db, err := sql.Open(v[0], v[1])
		if err != nil {
			log.Fatal("%s", err)
		}

		// escape password
		dsn := v[1]
		start := strings.IndexByte(dsn, ':')
		end := strings.IndexByte(dsn, '@')
		if start > 0 && end > 0 {
			dsn = dsn[:start+1] + "***" + dsn[end:]
		}

		zdb := DB{
			Base: &Base{
				trans: false,
			},
			DB:     db,
			name:   k,
			driver: v[0],
			dsn:    dsn,
		}
		zdb.Base.db = zdb

		dbMap[k] = zdb
		log.Info("ep=%s|func=install|name=%s|conf=%s", zdb.driver, zdb.name, zdb.dsn)
	}
	return dbMap
}

func GetDB(name string) *DB {
	if db, ok := dbMap[name]; ok {
		return &db
	} else {
		log.Error("can not get db [%s]", name)
		return nil
	}
}

func (t *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	var errstr = ""
	defer t.db.timeit(time.Now(), &errstr, "1", query, args...)

	result, err := t.Tx.Exec(query, args...)
	if err != nil {
		errstr = err.Error()
	}
	return result, err
}

func (t *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var errstr = ""
	defer t.db.timeit(time.Now(), &errstr, "1", query, args...)

	result, err := t.Tx.Query(query, args...)
	if err != nil {
		errstr = err.Error()
	}
	return result, err
}

func (t *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	var errstr = ""
	defer t.db.timeit(time.Now(), &errstr, "1", query, args...)

	return t.Tx.QueryRow(query, args...)
}

func (t *Tx) Commit() error {
	d := t.db
	log.Info("ep=%s|name=%s|func=commit", d.driver, d.name)
	return t.Tx.Commit()
}

func (t *Tx) Rollback() error {
	d := t.db
	log.Info("ep=%s|name=%s|func=rollback", d.driver, d.name)
	return t.Tx.Rollback()
}

func (d *DB) timeit(start time.Time, errstr *string, trans string, query string, args ...interface{}) {
	stat := d.DB.Stats()
	duration := time.Now().Sub(start)
	if len(*errstr) == 0 {
		log.Info("ep=%s|name=%s|use=%d|idle=%d|max=%d|wait=%d|waittime=%d|time=%d|trans=%s|sql=%s",
			d.driver, d.name, stat.InUse, stat.Idle, stat.MaxOpenConnections, stat.WaitCount,
			stat.WaitDuration/time.Microsecond, duration/time.Microsecond, trans,
			d.FormatSql(query, args...),
		)
	} else {
		log.Warn("ep=%s|name=%s|use=%d|idle=%d|max=%d|wait=%d|waittime=%d|time=%d|trans=%s|sql=%s|err=%s",
			d.driver, d.name, stat.InUse, stat.Idle, stat.MaxOpenConnections, stat.WaitCount,
			stat.WaitDuration/time.Microsecond, duration/time.Microsecond, trans,
			d.FormatSql(query, args...), *errstr,
		)
	}
}

func (d *DB) Begin() (*Tx, error) {
	log.Info("ep=%s|name=%s|func=begin", d.driver, d.name)
	tx, err := d.DB.Begin()
	ztx := Tx{
		Base: &Base{
			trans: true,
		},
		Tx: tx,
		db: d,
	}
	ztx.Base.tx = ztx
	return &ztx, err
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	var errstr = ""
	defer d.timeit(time.Now(), &errstr, "0", query, args...)

	result, err := d.DB.Exec(query, args...)
	if err != nil {
		errstr = err.Error()
	}
	return result, err
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var errstr = ""
	defer d.timeit(time.Now(), &errstr, "0", query, args...)

	result, err := d.DB.Query(query, args...)
	if err != nil {
		errstr = err.Error()
	}
	return result, err
}
func (d *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	var errstr = ""
	defer d.timeit(time.Now(), &errstr, "0", query, args...)

	return d.DB.QueryRow(query, args...)
}

func (b *Base) QueryMap(query string, args ...interface{}) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	var rows *sql.Rows
	var err error

	if b.trans {
		rows, err = b.tx.Query(query, args...)
	} else {
		rows, err = b.db.Query(query, args...)
	}
	if err != nil {
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		values := make([]interface{}, len(cols))
		for i := range values {
			values[i] = new(interface{})
		}

		err := rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, col := range cols {
			m[col] = *(values[i].((*interface{})))
		}
		data = append(data, m)
	}

	return data, nil
}

func (b *Base) FormatSql(query string, args ...interface{}) string {
	if len(args) == 0 {
		return query
	}
	if strings.Count(query, "?") != len(args) {
		log.Warn("format sql error %s %s", query, args)
		return query
	}
	// XXX for logging only, not real sql
	// TODO pg is not '?'
	query = strings.Replace(query, "?", "[%+v]", -1)
	return fmt.Sprintf(query, args...)
}

func (b *Base) Select(table string, where Where) ([]map[string]interface{}, error) {
	sql, args := builder.Select(table, builder.Where(where))
	return b.QueryMap(sql, args...)
}
func (b *Base) Insert(table string, value Values) (sql.Result, error) {
	sql, args := builder.Insert(table, builder.Values(value))

	if b.trans {
		return b.tx.Exec(sql, args...)
	} else {
		return b.db.Exec(sql, args...)
	}
}

func (b *Base) Update(table string, value Values, where Where) (sql.Result, error) {

	sql, args := builder.Update(table, builder.Values(value), builder.Where(where))

	if b.trans {
		return b.tx.Exec(sql, args...)
	} else {
		return b.db.Exec(sql, args...)
	}
}

func (b *Base) Delete(table string, where Where) (sql.Result, error) {

	sql, args := builder.Delete(table, builder.Where(where))

	if b.trans {
		return b.tx.Exec(sql, args...)
	} else {
		return b.db.Exec(sql, args...)
	}
}
