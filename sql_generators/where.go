package sql_generators

import (
	"errors"
)

var (
	options = [] string {
		"=", "<", ">", "<=", ">=", "<>", "!=", "<=>",
		"ike", "ike binary", "ot like", "like",
		"&", "|", "^", "<<", ">>",
		"like", "ot rlike", "egexp", "ot regexp",
		"~", "~*", "!~", "!~*", "imilar to",
		"ot similar to", "ot ilike", "~~*", "!~~*",}
)

type whereInfo struct {
	boolean bool			// false AND     true OR
	whereType string
	conditionType string
	field string
	value interface{}
	build *SqlGenerator
}

func newWhere(whereType, conditionType, field string, value interface{}, boolean bool) (*whereInfo, error) {
	return &whereInfo{
		boolean: boolean,
		whereType: whereType,
		conditionType: conditionType,
		field: field,
		value: value,
	}, nil
}

func newBuild(build *SqlGenerator, whereType, conditionType string, boolean bool) (*whereInfo, error) {
	return &whereInfo{
		boolean: boolean,
		whereType: whereType,
		conditionType: conditionType,
		build: build,
	}, nil
}

func (this *whereInfo) Where(args... interface{}) (*whereInfo, error) {
	boolean := args[len(args)-1]
	args = args[:len(args) - 1]
	if len(args) == 0 {
		panic("Where param error")
	}

	var field, conditionType string
	var value interface{}
	switch len(args) {
		case 1:
			value = args[0]
		case 2:
			field = args[0].(string)
			value = args[1]
			conditionType = "="
		case 3:
			field = args[0].(string)
			value = args[2]
			conditionType = args[1].(string)
		default:
			panic("Where func params num up to 3")
	}

	if !InvalidOperator(conditionType) {
		return nil, errors.New("condition option error")
	}

	return newWhere("Basic", conditionType, field, value, boolean == 0)
}

func (this *whereInfo) WhereFunc(callback func(build *SqlGenerator) *SqlGenerator, build *SqlGenerator, boolean bool) (*whereInfo, error) {
	return newBuild(callback(build), "Nested", "", boolean)
}

func (this *whereInfo) WhereExists(conditionType string, callback func(build *SqlGenerator) *SqlGenerator, boolean bool) (*whereInfo, error) {
	return newBuild(callback(NewGenerator()), "Exists", conditionType, boolean)
}

func (this *whereInfo) WhereArray(arrayWhere [][]interface{}, newBuild *SqlGenerator, boolean bool) (*whereInfo, error) {
	return this.WhereFunc(func(build *SqlGenerator) *SqlGenerator {
		for _, value := range arrayWhere {
			if boolean {
				_ = build.OrWhere(value...)
			} else {
				_ = build.Where(value...)
			}
		}
		return build
	}, newBuild, boolean)
}

func (this *whereInfo) WhereMap(mapWhere map[string] interface{}, newBuild *SqlGenerator, boolean bool) (*whereInfo, error) {
	return this.WhereFunc(func(build *SqlGenerator) *SqlGenerator {
		for key, value := range mapWhere {
			_ = build.Where(key, value)
		}
		return build
	}, newBuild, boolean)
}

func (this *whereInfo) whereInFactory(isNot bool, field string, listValue [] interface{}, boolean bool) (*whereInfo, error) {

	conditionType := "IN"
	if isNot {
		conditionType = "NOT IN"
	}
	return newWhere("In", conditionType, field, listValue, boolean)
}


func (this *whereInfo) WhereIn(field string, listValue [] interface{}, boolean bool) (*whereInfo, error) {
	return this.whereInFactory(false ,field, listValue, boolean)
}

func (this *whereInfo) WhereNotIn(field string, listValue [] interface{}, boolean bool) (*whereInfo, error) {
	return this.whereInFactory(true ,field, listValue, boolean)
}

func (this *whereInfo) whereBetweenFactory(isNot bool, field string, interval [] interface{}, boolean bool) (*whereInfo, error) {
	conditionType := "Between"
	if isNot {
		conditionType = "NOT Between"
	}

	return newWhere("Between", conditionType, field, interval, boolean)
}

func (this *whereInfo) WhereBetween(field string, interval [] interface{}, boolean bool) (*whereInfo, error) {
	if len(field) == 0 || len(interval) != 2 {
		panic("WhereBetween param error")
	}
	return this.whereBetweenFactory(false, field, interval, boolean)
}

func (this *whereInfo) WhereNotBetween(field string, interval [] interface{}, boolean bool) (*whereInfo, error) {
	return this.whereBetweenFactory(true, field, interval, boolean)
}

func (this *whereInfo) WhereNull(field string, boolean bool) (*whereInfo, error) {
	return newWhere("Null","IS", field, "NULL", boolean)
}

func (this *whereInfo) WhereNotNull(field string, boolean bool) (*whereInfo, error) {
	return newWhere("Null","GIS NOT", field, "NULL", boolean)
}

func (this *whereInfo) WhereDate(field string, date string, boolean bool) (*whereInfo, error) {
	char := GetDateStringJoiner(date)
	format := SetMysqlDateFormatByChar(char)
	func_field := "DATE_FORMAT(" + field + ", '" + format + "')"
	return newWhere("Func", "=", func_field, date, boolean)
}

func (this *whereInfo) WhereMonth(field string, month int, boolean bool) (*whereInfo, error) {
	func_field := "DATE_FORMAT(" + field + ", '%m')"
	return newWhere("Func","=", func_field, month, boolean)
}

func (this *whereInfo) WhereDay(field string, day int, boolean bool) (*whereInfo, error) {
	func_field := "DATE_FORMAT(" + field + ", '%d')"
	return newWhere("Func","=", func_field, day, boolean)
}

func (this *whereInfo) whereYear(field string, year int, boolean bool) (*whereInfo, error) {
	func_field := "DATE_FORMAT(" + field + ", '%Y')"
	return newWhere("Func","=", func_field, year, boolean)
}

func (this *whereInfo) WhereTime(field, condition, timestamp string, boolean bool) (*whereInfo, error) {
	func_field := "DATE_FORMAT(" + field + ", '%H:%i:%s')"
	return newWhere("Func", condition, func_field, timestamp, boolean)
}

func (this *whereInfo) WhereRaw(sql string, boolean bool) (*whereInfo, error) {
	return newWhere("Raw", "", "", sql, boolean)
}