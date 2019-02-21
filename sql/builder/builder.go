// build sql just like use zpy/base/dbpool.py
// not build for all sql

// TODO use $1 instead of ? for pg
// TODO SelectJoin
// TODO InsertMany

package builder

import (
	"fmt"
	"reflect"
	"strings"
)

type Where map[string]interface{}
type Values map[string]interface{}

func Select(table string, where Where) (string, []interface{}) {

	var args []interface{}

	field := "*"
	groupby := ""
	having := Where{}
	other := ""

	if value, ok := where["_field"]; ok {
		field = value.(string)
		delete(where, "_field")
	}
	if value, ok := where["_groupby"]; ok {
		groupby = value.(string)
		delete(where, "_groupby")
	}
	if value, ok := where["_having"]; ok {
		having = value.(Where)
		delete(where, "_having")
	}
	if value, ok := where["_other"]; ok {
		other = value.(string)
		delete(where, "_other")
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("SELECT %s FROM `%s`", field, table))

	// where
	if len(where) > 0 {
		sql, arg := where2sql(where)
		sb.WriteString(" WHERE ")
		sb.WriteString(sql)
		args = append(args, arg...)
	}

	// groupby
	if groupby != "" {
		sb.WriteString(" GROUP BY ")
		sb.WriteString(groupby)
	}

	// having
	if len(having) > 0 {
		sql, arg := where2sql(having)
		sb.WriteString(" HAVING ")
		sb.WriteString(sql)
		args = append(args, arg...)
	}

	// orderby limit offset
	if other != "" {
		sb.WriteString(" ")
		sb.WriteString(other)
	}

	return sb.String(), args
}

func Insert(table string, value Values) (string, []interface{}) {
	k, v, i := values2insert(value)
	sql := fmt.Sprintf("INSERT INTO `%s`(%s) VALUES(%s)", table, k, v)
	return sql, i
}

func Update(table string, value Values, where Where) (string, []interface{}) {
	var args []interface{}

	other := ""
	if value, ok := where["_other"]; ok {
		other = value.(string)
		delete(where, "_other")
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("UPDATE `%s`", table))

	// set
	k, v := values2set(value)
	sb.WriteString(" SET ")
	sb.WriteString(k)
	args = append(args, v...)

	// where
	if len(where) > 0 {
		sql, arg := where2sql(where)
		sb.WriteString(" WHERE ")
		sb.WriteString(sql)
		args = append(args, arg...)
	}

	// orderby limit offset
	if other != "" {
		sb.WriteString(" ")
		sb.WriteString(other)
	}

	return sb.String(), args
}

func Delete(table string, where Where) (string, []interface{}) {
	var args []interface{}

	other := ""
	if value, ok := where["_other"]; ok {
		other = value.(string)
		delete(where, "_other")
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("DELETE FROM `%s`", table))

	// where
	if len(where) > 0 {
		sql, arg := where2sql(where)
		sb.WriteString(" WHERE ")
		sb.WriteString(sql)
		args = append(args, arg...)
	}

	// orderby limit offset
	if other != "" {
		sb.WriteString(" ")
		sb.WriteString(other)
	}

	return sb.String(), args
}

func FormatSql(query string, args ...interface{}) string {
	if len(args) == 0 {
		return query
	}
	if strings.Count(query, "?") != len(args) {
		return query
	}
	// XXX for logging only, not real sql
	// TODO pg is not '?'
	query = strings.Replace(query, "?", "[%+v]", -1)
	return fmt.Sprintf(query, args...)
}

func values2insert(values Values) (string, string, []interface{}) {
	var args []interface{}
	var name []string
	var value []string
	for k, v := range values {
		name = append(name, fmt.Sprintf("`%s`", k))
		value = append(value, "?")
		args = append(args, v)
	}
	return strings.Join(name, ","), strings.Join(value, ","), args
}

func values2set(values Values) (sql string, args []interface{}) {
	var sqls []string
	for k, v := range values {
		sqls = append(sqls, fmt.Sprintf("`%s`=?", k))
		args = append(args, v)
	}
	return strings.Join(sqls, ","), args
}

func where2sql(where Where) (sql string, args []interface{}) {
	var key, op string
	var sqls []string
	for k, v := range where {
		k = strings.Trim(k, " ")
		idx := strings.IndexByte(k, ' ')
		if idx == -1 {
			key = k
			op = "="
		} else {
			key = k[:idx]
			op = k[idx+1:]
		}
		s, i := exp2sql(key, op, v)
		sqls = append(sqls, s)
		args = append(args, i...)
	}
	return strings.Join(sqls, " and "), args
}

func exp2sql(key string, op string, value interface{}) (sql string, args []interface{}) {

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("(`%s` %s ", key, op))

	if strings.Contains(op, "in") {
		builder.WriteString("(")
		for idx, v := range interface2slice(value) {
			if idx == 0 {
				builder.WriteString("?")
			} else {
				builder.WriteString(",?")
			}
			args = append(args, v)
		}
		builder.WriteString("))")
	} else if strings.Contains(op, "between") {
		builder.WriteString("? and ?)")
		v := interface2slice(value)
		args = append(args, v[0])
		args = append(args, v[1])
	} else {
		builder.WriteString("?)")
		args = append(args, value)
	}
	return builder.String(), args
}
func interface2slice(value interface{}) []interface{} {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice {
		return nil
	}
	s := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		s[i] = v.Index(i).Interface()
	}
	return s
}
