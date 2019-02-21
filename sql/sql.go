// use sql not orm
// use simple sql not join

// use go sql just like python dbpool.py
// ref : https://github.com/JoveYu/zpy/blob/master/base/dbpool.py

package sql

import (
	"database/sql"
	"strings"
	"time"

	"github.com/JoveYu/zgo/log"
	"github.com/JoveYu/zgo/sql/builder"
	"github.com/JoveYu/zgo/sql/scanner"
)

var (
	dbMap = make(map[string]DB)
)

type DBTool struct {
	db *DB
	tx *Tx
}

type DB struct {
	*DBTool
	*sql.DB
	name   string
	driver string
	dsn    string
}

type Tx struct {
	*DBTool
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
			DB:     db,
			name:   k,
			driver: v[0],
			dsn:    dsn,
		}
		zdb.DBTool = &DBTool{db: &zdb}

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
			builder.FormatSql(query, args...),
		)
	} else {
		log.Warn("ep=%s|name=%s|use=%d|idle=%d|max=%d|wait=%d|waittime=%d|time=%d|trans=%s|sql=%s|err=%s",
			d.driver, d.name, stat.InUse, stat.Idle, stat.MaxOpenConnections, stat.WaitCount,
			stat.WaitDuration/time.Microsecond, duration/time.Microsecond, trans,
			builder.FormatSql(query, args...), *errstr,
		)
	}
}

func (d *DB) Begin() (*Tx, error) {
	log.Info("ep=%s|name=%s|func=begin", d.driver, d.name)
	tx, err := d.DB.Begin()
	ztx := Tx{
		Tx: tx,
		db: d,
	}
	ztx.DBTool = &DBTool{tx: &ztx}
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

func (d *DBTool) QueryScan(obj interface{}, query string, args ...interface{}) error {
	var rows *sql.Rows
	var err error

	if d.tx != nil {
		rows, err = d.tx.Query(query, args...)
	} else {
		rows, err = d.db.Query(query, args...)
	}
	if err != nil {
		return err
	}
	defer rows.Close()

	err = scanner.ScanStruct(rows, obj)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBTool) SelectScan(obj interface{}, table string, where Where) error {
	sql, args := builder.Select(table, d.escapeWhere(where))
	return d.QueryScan(obj, sql, args...)
}

func (d *DBTool) QueryMap(query string, args ...interface{}) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	var rows *sql.Rows
	var err error

	if d.tx != nil {
		rows, err = d.tx.Query(query, args...)
	} else {
		rows, err = d.db.Query(query, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (d *DBTool) SelectMap(table string, where Where) ([]map[string]interface{}, error) {
	sql, args := builder.Select(table, d.escapeWhere(where))
	return d.QueryMap(sql, args...)
}

func (d *DBTool) Select(table string, where Where) (*sql.Rows, error) {
	sql, args := builder.Select(table, d.escapeWhere(where))
	if d.tx != nil {
		return d.tx.Query(sql, args...)
	} else {
		return d.db.Query(sql, args...)
	}
}
func (d *DBTool) Insert(table string, value Values) (sql.Result, error) {
	sql, args := builder.Insert(table, builder.Values(value))

	if d.tx != nil {
		return d.tx.Exec(sql, args...)
	} else {
		return d.db.Exec(sql, args...)
	}
}

func (d *DBTool) Update(table string, value Values, where Where) (sql.Result, error) {

	sql, args := builder.Update(table, builder.Values(value), d.escapeWhere(where))

	if d.tx != nil {
		return d.tx.Exec(sql, args...)
	} else {
		return d.db.Exec(sql, args...)
	}
}

func (d *DBTool) Delete(table string, where Where) (sql.Result, error) {

	sql, args := builder.Delete(table, d.escapeWhere(where))

	if d.tx != nil {
		return d.tx.Exec(sql, args...)
	} else {
		return d.db.Exec(sql, args...)
	}
}

func (d *DBTool) escapeWhere(where Where) builder.Where {
	if value, ok := where["_having"]; ok {
		where["_having"] = builder.Where(value.(Where))
	}
	return builder.Where(where)
}
