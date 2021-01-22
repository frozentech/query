package query

import (
	"fmt"
	"reflect"
	"strings"
)

// Contains ....
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// SelectDistinctBuilder ...
func SelectDistinctBuilder(model interface{}, table string, complexQuery bool) string {
	var (
		fields string
	)

	columns, _ := builder(model)

	if !complexQuery {
		fields = strings.Join(prepareColumns(columns, ""), ", ")
	} else {
		fields = strings.Join(prepareColumns(columns, table), ", ")
	}

	return fmt.Sprintf("SELECT DISTINCT %s FROM `%s`",
		fields,
		table)
}

// UpdateBuilder ...
func UpdateBuilder(model interface{}, table string, ignoreFields ...string) (string, []interface{}) {
	columns, values := builder(model)

	fields := []string{}
	items := []interface{}{}
	for i, column := range columns {
		if !Contains(ignoreFields, column) {
			fields = append(fields, fmt.Sprintf(" %s = ? ", column))
			items = append(items, values[i])
		}
	}

	sql := fmt.Sprintf("UPDATE `%s` SET %s ", table, strings.Join(fields, ", "))

	return sql, items
}

// InsertBuilder ...
func InsertBuilder(model interface{}, table string, complexQuery bool, ignoreFields ...string) (string, []interface{}) {
	columns, values := builder(model, ignoreFields...)

	var (
		fields string
	)

	fieldItems := []string{}
	for _, column := range columns {
		if !Contains(ignoreFields, column) {
			fieldItems = append(fieldItems, column)
		}
	}

	if complexQuery {
		fields = strings.Join(prepareColumns(fieldItems, table), ", ")
	} else {
		fields = strings.Join(prepareColumns(fieldItems, ""), ", ")
	}

	sql := fmt.Sprintf("INSERT INTO `%s` ( %s ) VALUES ( %s )",
		table,
		fields,
		strings.Join(prepareValues(fieldItems), ", "))
	return sql, values
}

// SelectBuilder ...
func SelectBuilder(model interface{}, table string, complexQuery bool) string {

	var (
		fields string
	)

	columns, _ := builder(model)

	if !complexQuery {
		fields = strings.Join(prepareColumns(columns, ""), ", ")
	} else {
		fields = strings.Join(prepareColumns(columns, table), ", ")
	}

	return fmt.Sprintf("SELECT %s FROM `%s`",
		fields,
		table)
}

// prepareColumns ...
func prepareColumns(columns []string, table string) []string {
	var formatColumns []string
	for _, value := range columns {
		if table == "" {
			formatColumns = append(formatColumns, fmt.Sprintf("`%s`", value))
		} else {
			formatColumns = append(formatColumns, fmt.Sprintf("`%s`.`%s`", table, value))
		}

	}

	return formatColumns
}

// prepareValues ...
func prepareValues(columns []string) []string {
	var formatColumns []string

	for i := 0; i < len(columns); i++ {
		formatColumns = append(formatColumns, fmt.Sprintf("%s", "?"))
	}

	return formatColumns
}

// builder ...
func builder(model interface{}, ignoreFields ...string) (columns []string, values []interface{}) {
	t := reflect.ValueOf(model).Elem()

	for i := 0; i < t.NumField(); i++ {
		f := t.Type().Field(i)
		tag := f.Tag.Get("db")
		if "" == tag || "-" == tag {
			continue
		}
		tags := strings.Split(tag, ",")

		if Contains(ignoreFields, tags[0]) {
			continue
		}

		columns = append(columns, tags[0])
		values = append(values, t.Field(i).Interface())
	}
	return
}
