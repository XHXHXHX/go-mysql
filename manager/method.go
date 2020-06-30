package manager

import "github.com/XHXHXHX/go-mysql/sql_generators"

func SetPrefix(s string) {
	sql_generators.SetPrefix(s)
}

func (manage *Manager) Table(args... string) *Manager {
	manage.Generator.Table(args...)
	return manage
}
func (manage *Manager) JoinTable(table string) *Manager {
	manage.Generator.JoinTable(table)
	return manage
}
func (manage *Manager) Select(args... string) *Manager {
	manage.Generator.Select(args...)
	return manage
}
func (manage *Manager) SelectRaw(sql string) *Manager {
	manage.Generator.SelectRaw(sql)
	return manage
}
func (manage *Manager) JoinFactory(joinType, table, thatRelationField, relationCondition, thisRelationField string) *Manager {
	manage.Generator.JoinFactory(joinType, table, thatRelationField, relationCondition, thisRelationField)
	return manage
}
func (manage *Manager) Join(table, thatRelationField, relationCondition, thisRelationField string) *Manager {
	manage.Generator.Join(table, thatRelationField, relationCondition, thisRelationField)
	return manage
}
func (manage *Manager) LeftJoin(table, thatRelationField, relationCondition, thisRelationField string) *Manager {
	manage.Generator.LeftJoin(table, thatRelationField, relationCondition, thisRelationField)
	return manage
}
func (manage *Manager) RightJoin(table, thatRelationField, relationCondition, thisRelationField string) *Manager {
	manage.Generator.RightJoin(table, thatRelationField, relationCondition, thisRelationField)
	return manage
}
func (manage *Manager) InnerJoin(table, thatRelationField, relationCondition, thisRelationField string) *Manager {
	manage.Generator.InnerJoin(table, thatRelationField, relationCondition, thisRelationField)
	return manage
}
func (manage *Manager) JoinFuncFactory(joinType, table string, callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.JoinFuncFactory(joinType, table, callback)
	return manage
}
func (manage *Manager) LeftJoinFunc(table string, callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.LeftJoinFunc(table, callback)
	return manage
}
func (manage *Manager) RightJoinFunc(table string, callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.RightJoinFunc(table, callback)
	return manage
}
func (manage *Manager) InnerJoinFunc(table string, callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.InnerJoinFunc(table, callback)
	return manage
}
func (manage *Manager) JoinFunc(table string, callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.JoinFunc(table, callback)
	return manage
}
func (manage *Manager) Where(args... interface{}) *Manager {
	manage.Generator.Where(args...)
	return manage
}
func (manage *Manager) OrWhere(args... interface{}) *Manager {
	manage.Generator.OrWhere(args...)
	return manage
}
func (manage *Manager) WhereArray(arrayWhere [][]interface{}) *Manager {
	manage.Generator.WhereArray(arrayWhere)
	return manage
}
func (manage *Manager) OrWhereArray(arrayWhere [][]interface{}) *Manager {
	manage.Generator.OrWhereArray(arrayWhere)
	return manage
}
func (manage *Manager) WhereMap(mapWhere map[string] interface{}) *Manager {
	manage.Generator.WhereMap(mapWhere)
	return manage
}
func (manage *Manager) OrWhereMap(mapWhere map[string] interface{}) *Manager {
	manage.Generator.OrWhereMap(mapWhere)
	return manage
}
func (manage *Manager) WhereIn(field string, listValue [] interface{}) *Manager {
	manage.Generator.WhereIn(field, listValue)
	return manage
}
func (manage *Manager) OrWhereIn(field string, listValue [] interface{}) *Manager {
	manage.Generator.OrWhereIn(field, listValue)
	return manage
}
func (manage *Manager) WhereNotIn(field string, listValue [] interface{})  *Manager {
	manage.Generator.WhereNotIn(field, listValue)
	return manage
}
func (manage *Manager) OrWhereNotIn(field string, listValue [] interface{})  *Manager {
	manage.Generator.OrWhereNotIn(field, listValue)
	return manage
}
func (manage *Manager) WhereBetween(field string, listValue [] interface{})  *Manager {
	manage.Generator.WhereBetween(field, listValue)
	return manage
}
func (manage *Manager) OrWhereBetween(field string, listValue [] interface{})  *Manager {
	manage.Generator.OrWhereBetween(field, listValue)
	return manage
}
func (manage *Manager) WhereNotBetween(field string, listValue [] interface{})  *Manager {
	manage.Generator.WhereNotBetween(field, listValue)
	return manage
}
func (manage *Manager) OrWhereNotBetween(field string, listValue [] interface{})  *Manager {
	manage.Generator.OrWhereNotBetween(field, listValue)
	return manage
}
func (manage *Manager) WhereNull(field string) *Manager {
	manage.Generator.WhereNull(field)
	return manage
}
func (manage *Manager) OrWhereNull(field string) *Manager {
	manage.Generator.OrWhereNull(field)
	return manage
}
func (manage *Manager) WhereNotNull(field string) *Manager {
	manage.Generator.WhereNotNull(field)
	return manage
}
func (manage *Manager) OrWhereNotNull(field string) *Manager {
	manage.Generator.OrWhereNotNull(field)
	return manage
}
func (manage *Manager) WhereDate(field, date string) *Manager {
	manage.Generator.WhereDate(field, date)
	return manage
}
func (manage *Manager) OrWhereDate(field, date string) *Manager {
	manage.Generator.OrWhereDate(field, date)
	return manage
}
func (manage *Manager) WhereMonth(field string, month int) *Manager {
	manage.Generator.WhereMonth(field, month)
	return manage
}
func (manage *Manager) OrWhereMonth(field string, month int) *Manager {
	manage.Generator.OrWhereMonth(field, month)
	return manage
}
func (manage *Manager) WhereDay(field string, day int) *Manager {
	manage.Generator.WhereDay(field, day)
	return manage
}
func (manage *Manager) OrWhereDay(field string, day int) *Manager {
	manage.Generator.OrWhereDay(field, day)
	return manage
}
func (manage *Manager) WhereYear(field string, year int) *Manager {
	manage.Generator.WhereYear(field, year)
	return manage
}
func (manage *Manager) OrWhereYear(field string, year int) *Manager {
	manage.Generator.OrWhereYear(field, year)
	return manage
}
func (manage *Manager) WhereTime(field, condition, timestamp string) *Manager {
	manage.Generator.WhereTime(field, condition, timestamp)
	return manage
}
func (manage *Manager) OrWhereTime(field, condition, timestamp string) *Manager {
	manage.Generator.OrWhereTime(field, condition, timestamp)
	return manage
}
func (manage *Manager) WhereFunc(callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.WhereFunc(callback)
	return manage
}
func (manage *Manager) OrWhereFunc(callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.OrWhereFunc(callback)
	return manage
}
func (manage *Manager) WhereExists(callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.WhereExists(callback)
	return manage
}
func (manage *Manager) OrWhereExists(callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.OrWhereExists(callback)
	return manage
}
func (manage *Manager) WhereNotExists(field string, callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.WhereNotExists(field, callback)
	return manage
}
func (manage *Manager) OrWhereNotExists(field string, callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.OrWhereNotExists(field, callback)
	return manage
}
func (manage *Manager) When(boolean bool, callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.When(boolean, callback)
	return manage
}
func (manage *Manager) OrWhen(boolean bool, callback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.OrWhen(boolean, callback)
	return manage
}
func (manage *Manager) WhenElse(boolean bool, trueCallback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator, falseCallback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.WhenElse(boolean, trueCallback, falseCallback)
	return manage
}
func (manage *Manager) OrWhenElse(boolean bool, trueCallback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator, falseCallback func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator) *Manager {
	manage.Generator.OrWhenElse(boolean, trueCallback, falseCallback)
	return manage
}
func (manage *Manager) OrderBy(args... string) *Manager {
	manage.Generator.OrderBy(args...)
	return manage
}
func (manage *Manager) OrderByRaw(sql string) *Manager {
	manage.Generator.OrderByRaw(sql)
	return manage
}
func (manage *Manager) GroupBy(args... string) *Manager {
	manage.Generator.GroupBy(args...)
	return manage
}
func (manage *Manager) GroupByRaw(sql string) *Manager {
	manage.Generator.GroupByRaw(sql)
	return manage
}
func (manage *Manager) Having(args... interface{}) *Manager {
	manage.Generator.Having(args...)
	return manage
}
func (manage *Manager) HavingRaw(sql string) *Manager {
	manage.Generator.HavingRaw(sql)
	return manage
}
func (manage *Manager) Offset(num int) *Manager {
	manage.Generator.Offset(num)
	return manage
}
func (manage *Manager) Limit(num int) *Manager {
	manage.Generator.Limit(num)
	return manage
}
