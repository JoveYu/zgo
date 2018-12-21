// TODO select join
// TODO select on
// TODO select having

package sql

import (
    "fmt"
    "time"
    "strings"
    "reflect"
    "database/sql"
    "github.com/JoveYu/zgo/log"
)

var (
    dbMap = make(map[string]DB)
)

type Base struct {
    db DB
    tx Tx
    trans bool
}

type DB struct {
    *Base
    *sql.DB
    name string
    driver string
    dsn string
}

type Tx struct {
    *Base
    *sql.Tx
    db *DB
}

func Install(conf map[string][]string) map[string]DB {
    log.Debug("available sql driver: %s", sql.Drivers())
    for k,v := range conf {
        if len(v) !=2 {
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
                trans:false,
            },
            DB: db,
            name: k,
            driver: v[0],
            dsn: dsn,
        }
        zdb.Base.db = zdb

        dbMap[k] = zdb
        log.Info("ep=%s|func=install|name=%s|conf=%s", zdb.driver, zdb.name, zdb.dsn)
    }
    return dbMap
}

func GetDB(name string) *DB {
    if db, ok:=dbMap[name]; ok {
        return &db
    }else{
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

func (t *Tx) QueryRow(query string, args ...interface{}) (*sql.Row) {
    var errstr = ""
    defer t.db.timeit(time.Now(), &errstr, "1", query, args...)

    return t.Tx.QueryRow(query, args...)
}

func (t *Tx) Commit() error {
    d := t.db
    log.Info("eq=%s|func=commit|name=%s", d.driver, d.name)
    return t.Tx.Commit()
}

func (t *Tx) Rollback() error {
    d := t.db
    log.Info("eq=%s|func=rollback|name=%s", d.driver, d.name)
    return t.Tx.Rollback()
}

func (d *DB) timeit(start time.Time, errstr *string, trans string, query string, args ...interface{}) {
    stat := d.DB.Stats()
    duration := time.Now().Sub(start)
    if len(*errstr) == 0 {
        log.Info("eq=%s|name=%s|use=%d|idle=%d|max=%d|wait=%d|waittime=%d|time=%d|trans=%s|sql=%s",
            d.driver, d.name, stat.InUse, stat.Idle, stat.MaxOpenConnections, stat.WaitCount,
            stat.WaitDuration/time.Microsecond, duration/time.Microsecond, trans,
            d.FormatSql(query, args...),
        )
    } else {
        log.Warn("eq=%s|name=%s|use=%d|idle=%d|max=%d|wait=%d|waittime=%d|time=%d|trans=%s|sql=%s|err=%s",
            d.driver, d.name, stat.InUse, stat.Idle, stat.MaxOpenConnections, stat.WaitCount,
            stat.WaitDuration/time.Microsecond, duration/time.Microsecond, trans,
            d.FormatSql(query, args...), *errstr,
        )
    }
}

func (d *DB) Begin() (*Tx, error) {
    log.Info("eq=%s|func=begin|name=%s", d.driver, d.name)
    tx,err := d.DB.Begin()
    ztx := Tx{
        Base: &Base{
            trans:true,
        },
        Tx:tx,
        db:d,
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
func (d *DB) QueryRow(query string, args ...interface{}) (*sql.Row) {
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
    // TODO 与实际sql 有差距 仅供参考
    query = strings.Replace(query, "?", "[%+v]", -1)
    return fmt.Sprintf(query, args...)
}

func (b *Base) Select(table string, where map[string]interface{}) []map[string]interface{} {
    query := "select %s from `%s`%s%s %s"
    field := "*"
    other := ""

    if value, ok := where["_field"]; ok {
        field = value.(string)
        delete(where, "_field")
    }
    if value, ok := where["_other"]; ok {
        other = value.(string)
        delete(where, "_other")
    }
    k, v := b.where2sql(where)
    if len(v) == 0 {
        query = fmt.Sprintf(query, field, table, "", "", other)
    } else {
        query = fmt.Sprintf(query, field, table, " where ", k, other)
    }
    result, err := b.QueryMap(query, v...)
    if err != nil {
        return nil
    }
    return result
}
func (b *Base) Insert(table string, values map[string]interface{}) int64 {
    query := "insert into %s(%s) values(%s)"
    k,v,i := b.insert2sql(values)

    query = fmt.Sprintf(query, table, k, v)

    var result sql.Result
    var err error
    if b.trans {
        result, err = b.tx.Exec(query, i...)
    } else {
        result, err = b.db.Exec(query, i...)
    }

    if err != nil {
        return int64(-1)
    }
    r, err := result.LastInsertId()
    if err != nil {
        return int64(-1)
    }
    return r
}

func (b *Base) Update(table string, values map[string]interface{}, where map[string]interface{}) int64 {
    query := "update %s set %s%s%s %s"
    other := ""
    var i []interface{}

    if value, ok := where["_other"]; ok {
        other = value.(string)
        delete(where, "_other")
    }
    set_k, set_v := b.update2sql(values)
    i = append(i, set_v...)
    where_k, where_v := b.where2sql(where)
    i = append(i, where_v...)

    if len(where_v) == 0 {
        query = fmt.Sprintf(query, table, set_k, "", "", other)
    } else {
        query = fmt.Sprintf(query, table, set_k, " where ", where_k, other)
    }

    var result sql.Result
    var err error
    if b.trans {
        result, err = b.tx.Exec(query, i...)
    } else {
        result, err = b.db.Exec(query, i...)
    }
    if err != nil {
        return int64(-1)
    }

    r, err := result.RowsAffected()
    if err != nil {
        return int64(-1)
    }
    return r
}

func (b *Base) Delete(table string, where map[string]interface{}) int64 {
    query := "delete from `%s`%s%s %s"
    other := ""

    if value, ok := where["_other"]; ok {
        other = value.(string)
        delete(where, "_other")
    }
    k, v := b.where2sql(where)
    if len(v) == 0 {
        query = fmt.Sprintf(query, table, "", "", other)
    } else {
        query = fmt.Sprintf(query, table, " where ", k, other)
    }

    var result sql.Result
    var err error
    if b.trans {
        result, err = b.tx.Exec(query, v...)
    } else {
        result, err = b.db.Exec(query, v...)
    }
    if err != nil {
        return int64(-1)
    }

    r, err := result.RowsAffected()
    if err != nil {
        return int64(-1)
    }
    return r
}

// from zpy/base/dbpool.py
func (b *Base) exp2sql(key string, op string, value interface{}) (string, []interface{} ){
    var i []interface{}

    builder := strings.Builder{}
    builder.WriteString(fmt.Sprintf("(`%s` %s ", key, op))

    if op == "in" {
        builder.WriteString("(")
        for idx, v := range b.interface2slice(value) {
            if idx == 0 {
                builder.WriteString("?")
            } else {
                builder.WriteString(",?")
            }
            i = append(i, v)
        }
        builder.WriteString("))")
    }else if op == "between" {
        builder.WriteString("? and ?)")
        v := b.interface2slice(value)
        i = append(i, v[0])
        i = append(i, v[1])
    }else{
        builder.WriteString("?)")
        i = append(i, value)
    }
    return builder.String(), i
}

// from zpy/base/dbpool.py
func (b *Base) where2sql(where map[string]interface{}) (string, []interface{}) {
    var key, op string
    var i []interface{}
    var s []string
    j := 0
    for k,v := range where {
        k = strings.Trim(k, " ")
        idx := strings.IndexByte(k, ' ')
        if idx == -1 {
            key = k
            op = "="
        } else {
            key = k[:idx]
            op = k[idx+1:]
        }
        ss, is := b.exp2sql(key, op, v)
        s = append(s, ss)
        i = append(i, is...)
        j++
    }
    return strings.Join(s, " and "), i
}

// from zpy/base/dbpool.py
func (b *Base) update2sql(values map[string]interface{}) (string, []interface{}) {
    var i []interface{}
    var s []string
    for k,v := range values {
        s = append(s, fmt.Sprintf("`%s`=?", k))
        i = append(i, v)
    }
    return strings.Join(s, ","), i
}

// from zpy/base/dbpool.py
func (b *Base) insert2sql(values map[string]interface{}) (string, string, []interface{}) {
    var i []interface{}
    var s []string
    var ss []string
    for k,v := range values {
        s = append(s, fmt.Sprintf("`%s`", k))
        ss = append(ss, "?")
        i = append(i, v)
    }
    return strings.Join(s, ","), strings.Join(ss, ","), i
}

func (b *Base) interface2slice(value interface{}) []interface{} {
    v := reflect.ValueOf(value)
    if v.Kind() != reflect.Slice {
        log.Error("can not convert interface to slice")
        return nil
    }
    s := make([]interface{}, v.Len())
    for i := 0; i<v.Len(); i++ {
        s[i] = v.Index(i).Interface()
    }
    return s
}

