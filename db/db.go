package db

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

type DB struct {
    *sql.DB
    name string
    driver string
    dsn string
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
            DB: db,
            name: k,
            driver: v[0],
            dsn: dsn,
        }
        dbMap[k] = zdb
        log.Info("ep=%s|func=install|name=%s|conf=%s", zdb.driver, zdb.name, zdb.dsn)
    }
    return dbMap
}

func GetDB(name string) *DB {
    if db, ok:=dbMap[name]; ok {
        return &db
    }else{
        log.Errord(1, "can not get db [%s]", name)
        return nil
    }
}

func (d *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
    var errstr = ""
    start_time := time.Now()

    result, err := d.DB.Exec(query, args...)
    if err != nil {
        errstr = err.Error()
    }

    stat := d.DB.Stats()
    end_time := time.Now()
    duration := end_time.Sub(start_time)

    if err != nil {
        log.Warn("eq=%s|name=%s|use=%d|idle=%d|max=%d|wait=%d|waittime=%d|time=%d|sql=%s|err=%s",
            d.driver, d.name, stat.InUse, stat.Idle, stat.MaxOpenConnections, stat.WaitCount,
            stat.WaitDuration/time.Microsecond, duration/time.Microsecond,
            d.FormatSql(query, args...), errstr,
        )
    } else {
        log.Info("eq=%s|name=%s|use=%d|idle=%d|max=%d|wait=%d|waittime=%d|time=%d|sql=%s",
            d.driver, d.name, stat.InUse, stat.Idle, stat.MaxOpenConnections, stat.WaitCount,
            stat.WaitDuration/time.Microsecond, duration/time.Microsecond,
            d.FormatSql(query, args...),
        )
    }
    return result, err
}

func (d *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
    var errstr = ""
    start_time := time.Now()

    result, err := d.DB.Query(query, args...)
    if err != nil {
        errstr = err.Error()
    }

    stat := d.DB.Stats()
    end_time := time.Now()
    duration := end_time.Sub(start_time)

    if err != nil {
        log.Warn("eq=%s|name=%s|use=%d|idle=%d|max=%d|wait=%d|waittime=%d|time=%d|sql=%s|err=%s",
            d.driver, d.name, stat.InUse, stat.Idle, stat.MaxOpenConnections, stat.WaitCount,
            stat.WaitDuration/time.Microsecond, duration/time.Microsecond,
            d.FormatSql(query, args...), errstr,
        )
    } else {
        log.Info("eq=%s|name=%s|use=%d|idle=%d|max=%d|wait=%d|waittime=%d|time=%d|sql=%s",
            d.driver, d.name, stat.InUse, stat.Idle, stat.MaxOpenConnections, stat.WaitCount,
            stat.WaitDuration/time.Microsecond, duration/time.Microsecond,
            d.FormatSql(query, args...),
        )
    }
    return result, err
}
func (d *DB) QueryRow(query string, args ...interface{}) (*sql.Row) {
    var errstr = ""
    start_time := time.Now()

    result:= d.DB.QueryRow(query, args...)

    stat := d.DB.Stats()
    end_time := time.Now()
    duration := end_time.Sub(start_time)

    log.Info("eq=%s|name=%s|use=%d|idle=%d|max=%d|wait=%d|waittime=%d|time=%d|sql=%s|err=%s",
        d.driver, d.name, stat.InUse, stat.Idle, stat.MaxOpenConnections, stat.WaitCount,
        stat.WaitDuration/time.Microsecond, duration/time.Microsecond,
        d.FormatSql(query, args...), errstr,
    )
    return result
}

func (d *DB) QueryMap(query string, args ...interface{}) ([]map[string]interface{}, error) {
    var data []map[string]interface{}

    rows, err := d.Query(query, args...)
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

func (d *DB) FormatSql(query string, args ...interface{}) string {
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

func (d *DB) Select(table string, where map[string]interface{}) []map[string]interface{} {
    sql := "select %s from `%s`%s%s %s"
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
    k, v := d.where2sql(where)
    if len(v) == 0 {
        sql = fmt.Sprintf(sql, field, table, "", "", other)
    } else {
        sql = fmt.Sprintf(sql, field, table, " where ", k, other)
    }
    result, err := d.QueryMap(sql, v...)
    if err != nil {
        return nil
    }
    return result
}
func (d *DB) Insert(table string, values map[string]interface{}) int64 {
    sql := "insert into %s(%s) values(%s)"
    k,v,i := d.insert2sql(values)

    sql = fmt.Sprintf(sql, table, k, v)
    result, err := d.Exec(sql, i...)
    if err != nil {
        return int64(-1)
    }
    r, err := result.LastInsertId()
    if err != nil {
        return int64(-1)
    }
    return r
}

func (d *DB) Update(table string, values map[string]interface{}, where map[string]interface{}) int64 {
    sql := "update %s set %s%s%s %s"
    other := ""
    var i []interface{}

    if value, ok := where["_other"]; ok {
        other = value.(string)
        delete(where, "_other")
    }
    set_k, set_v := d.update2sql(values)
    i = append(i, set_v...)
    where_k, where_v := d.where2sql(where)
    i = append(i, where_v...)

    if len(where_v) == 0 {
        sql = fmt.Sprintf(sql, table, set_k, "", "", other)
    } else {
        sql = fmt.Sprintf(sql, table, set_k, " where ", where_k, other)
    }
    result, err := d.Exec(sql, i...)
    if err != nil {
        return int64(-1)
    }

    r, err := result.RowsAffected()
    if err != nil {
        return int64(-1)
    }
    return r
}

func (d *DB) Delete(table string, where map[string]interface{}) int64 {
    sql := "delete from `%s`%s%s %s"
    other := ""

    if value, ok := where["_other"]; ok {
        other = value.(string)
        delete(where, "_other")
    }
    k, v := d.where2sql(where)
    if len(v) == 0 {
        sql = fmt.Sprintf(sql, table, "", "", other)
    } else {
        sql = fmt.Sprintf(sql, table, " where ", k, other)
    }
    result, err := d.Exec(sql, v...)
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
func (d *DB) exp2sql(key string, op string, value interface{}) (string, []interface{} ){
    var i []interface{}

    builder := strings.Builder{}
    builder.WriteString(fmt.Sprintf("(`%s` %s ", key, op))

    if op == "in" {
        builder.WriteString("(")
        for idx, v := range d.interface2slice(value) {
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
        v := d.interface2slice(value)
        i = append(i, v[0])
        i = append(i, v[1])
    }else{
        builder.WriteString("?)")
        i = append(i, value)
    }
    return builder.String(), i
}

// from zpy/base/dbpool.py
func (d *DB) where2sql(where map[string]interface{}) (string, []interface{}) {
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
        ss, is := d.exp2sql(key, op, v)
        s = append(s, ss)
        i = append(i, is...)
        j++
    }
    return strings.Join(s, " and "), i
}

// from zpy/base/dbpool.py
func (d *DB) update2sql(values map[string]interface{}) (string, []interface{}) {
    var i []interface{}
    var s []string
    for k,v := range values {
        s = append(s, fmt.Sprintf("`%s`=?", k))
        i = append(i, v)
    }
    return strings.Join(s, ","), i
}

// from zpy/base/dbpool.py
func (d *DB) insert2sql(values map[string]interface{}) (string, string, []interface{}) {
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

func (d *DB) interface2slice(value interface{}) []interface{} {
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

