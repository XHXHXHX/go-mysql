package sql_generators

import (
	"strconv"
	"strings"
	"time"
)

type SqlGenerator struct {
	ExeSql string
	ExeParam []interface{}
	ShowSql string
	exeType string
	runtime time.Duration
	sqlPart *sqlInfo
	joinTmp *joinInfo
}

var whereFactory *whereInfo
var prefix string

func init () {
	whereFactory = &whereInfo{}
}

func SetPrefix(s string) {
	prefix = s
}

func DB() *SqlGenerator {
	return NewGenerator()
}

func NewGenerator() *SqlGenerator {
	return &SqlGenerator{
		sqlPart: &sqlInfo{
			selectData: "*",
			offset:0,
		},
	}
}

func (this *SqlGenerator) build() *SqlGenerator {
	var sql string
	var params []interface{}
	switch this.exeType {
	case "SELECT":
		sql, params = this.sqlPart.BuildQuery()
	case "INSERT":
		sql, params = this.sqlPart.BuildInsert()
	case "UPDATE":
		sql, params = this.sqlPart.BuildUpdate()
	case "DELETE":
		sql, params = this.sqlPart.BuildDelete()
	case "ALTER" :
	default:

	}

	this.ExeSql = sql
	this.setShowSql(params)

	return this
}

func (this *SqlGenerator) NewGenerator() *SqlGenerator {
	build := NewGenerator()

	return build.Table(this.sqlPart.table, this.sqlPart.alias)
}

func (this *SqlGenerator) setShowSql(params []interface{}) {
	this.ShowSql = this.ExeSql

	for _, value := range params {
		switch val := value.(type) {
			case string:
				this.ShowSql = strings.Replace(this.ShowSql, "?", AddSingleSymbol(val), 1)
			case int:
				this.ShowSql = strings.Replace(this.ShowSql, "?", strconv.Itoa(val), 1)
			default:
				panic("param error")
		}
		this.ExeParam = append(this.ExeParam, value)
	}
}

func (this *SqlGenerator) whereResult(result *whereInfo, err error) *SqlGenerator {
	if err != nil {
		panic(err)
	}

	this.sqlPart.setWhere(result)
	return this
}

func (this *SqlGenerator) havingResult(result *whereInfo, err error) *SqlGenerator {
	if err != nil {
		panic(err)
	}

	this.sqlPart.setHaving(result)
	return this
}

func (this *SqlGenerator) Table(args... string) *SqlGenerator {
	if len(args) == 0 {
		panic("Table param error")
	}
	table, alias := args[0], ""
	if len(args) > 1 && len(args[1]) > 0 {
		alias = args[1]
	}

	this.sqlPart.setTable(prefix + table, alias)
	return this
}

func (this *SqlGenerator) JoinTable(table string) {
	if len(table) == 0 {
		panic("JoinTable param error")
	}
	var alias string = ""
	if strings.Count(table, " as ") == 1 {
		tmp := strings.Split(table, " as ")
		table = tmp[0]
		alias = tmp[1]
	}
	_ = this.Table(table, alias)
}

func (this *SqlGenerator) Select(args... string) *SqlGenerator {
	if len(args) > 0 {
		this.sqlPart.selectData = strings.Join(args, ",")
	} else {
		this.sqlPart.selectData = "*"
	}

	return this
}

func (this *SqlGenerator) SelectRaw(sql string) *SqlGenerator {
	if len(sql) == 0 {
		panic("SelectRaw param error")
	}
	this.sqlPart.selectData = sql
	return this
}

/***************************************  JOIN  **********************************************************/

func (this *SqlGenerator) JoinFactory(joinType, table, thatRelationField, relationCondition, thisRelationField string) *SqlGenerator {
	build := NewGenerator()
	build.JoinTable(table)
	this.sqlPart.setJoin(build, joinType, thatRelationField, relationCondition, thisRelationField)

	return this
}

func (this *SqlGenerator) Join(table, thatRelationField, relationCondition, thisRelationField string) *SqlGenerator {
	return this.JoinFactory("Inner Join", table, thatRelationField, relationCondition, thisRelationField)
}

func (this *SqlGenerator) LeftJoin(table, thatRelationField, relationCondition, thisRelationField string) *SqlGenerator {
	return this.JoinFactory("Left Join", table, thatRelationField, relationCondition, thisRelationField)
}

func (this *SqlGenerator) RightJoin(table, thatRelationField, relationCondition, thisRelationField string) *SqlGenerator {
	return this.JoinFactory("Right Join", table, thatRelationField, relationCondition, thisRelationField)
}

func (this *SqlGenerator) InnerJoin(table, thatRelationField, relationCondition, thisRelationField string) *SqlGenerator {
	return this.JoinFactory("Inner Join", table, thatRelationField, relationCondition, thisRelationField)
}

func (this *SqlGenerator) On(thatRelationField, relationCondition, thisRelationField string) *SqlGenerator {
	this.joinTmp = &joinInfo{}
	this.joinTmp.JoinOn(thatRelationField, relationCondition, thisRelationField)
	return this
}

func (this *SqlGenerator) JoinFuncFactory(joinType, table string, callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	build :=  NewGenerator()
	build.JoinTable(table)
	join_build := callback(build)

	tmpJoin := &joinInfo{
		buildInfo: join_build,
		joinType: joinType,
		thisRelationField: join_build.joinTmp.thisRelationField,
		thatRelationField: join_build.joinTmp.thatRelationField,
		relationCondition: join_build.joinTmp.relationCondition,
	}
	this.sqlPart.joinData = append(this.sqlPart.joinData, tmpJoin)
	return this
}

func (this *SqlGenerator) LeftJoinFunc(table string, callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	return this.JoinFuncFactory("Left Join", table, callback)
}

func (this *SqlGenerator) RightJoinFunc(table string, callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	return this.JoinFuncFactory("Right Join", table, callback)
}

func (this *SqlGenerator) InnerJoinFunc(table string, callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	return this.JoinFuncFactory("Inner Join", table, callback)
}

func (this *SqlGenerator) JoinFunc(table string, callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	return this.JoinFuncFactory("Inner Join", table, callback)
}
/*************************************  JOIN  END  *******************************************************/

/***************************************  WHERE  **********************************************************/

func (this *SqlGenerator) Where(args... interface{}) *SqlGenerator {
	if len(args) == 0 {
		return this
	}
	args = append(args, 1)
	return this.whereResult( whereFactory.Where(FormatWhereParam(args...)...))
}

func (this *SqlGenerator) OrWhere(args... interface{}) *SqlGenerator {
	if len(args) == 0 {
		panic("Where param error")
	}
	args = append(args, 0)
	return this.whereResult( whereFactory.Where(FormatWhereParam(args...)...))
}

func (this *SqlGenerator) WhereArray(arrayWhere [][]interface{}) *SqlGenerator {
	if len(arrayWhere) == 0 {
		panic("WhereArray param error")
	}
	return this.whereResult( whereFactory.WhereArray(arrayWhere, this.NewGenerator(), false))
}

func (this *SqlGenerator) OrWhereArray(arrayWhere [][]interface{}) *SqlGenerator {
	if len(arrayWhere) == 0 {
		panic("WhereArray param error")
	}
	return this.whereResult( whereFactory.WhereArray(arrayWhere, this.NewGenerator(), true))
}

func (this *SqlGenerator) WhereMap(mapWhere map[string] interface{}) *SqlGenerator {
	if len(mapWhere) == 0 {
		panic("WhereMap param error")
	}
	return this.whereResult( whereFactory.WhereMap(mapWhere, this.NewGenerator(), false))
}

func (this *SqlGenerator) OrWhereMap(mapWhere map[string] interface{}) *SqlGenerator {
	if len(mapWhere) == 0 {
		panic("WhereMap param error")
	}
	return this.whereResult( whereFactory.WhereMap(mapWhere, this.NewGenerator(), true))
}

func (this *SqlGenerator) WhereIn(field string, listValue [] interface{}) *SqlGenerator {
	if len(field) == 0 || len(listValue) == 0 {
		panic("WhereIn param error")
	}
	return this.whereResult( whereFactory.WhereIn(field, listValue, false))
}

func (this *SqlGenerator) OrWhereIn(field string, listValue [] interface{}) *SqlGenerator {
	if len(field) == 0 || len(listValue) == 0 {
		panic("WhereIn param error")
	}
	return this.whereResult( whereFactory.WhereIn(field, listValue, true))
}

func (this *SqlGenerator) WhereNotIn(field string, listValue [] interface{}) *SqlGenerator {
	if len(field) == 0 || len(listValue) == 0 {
		panic("WhereNotIn param error")
	}
	return this.whereResult( whereFactory.WhereNotIn(field, listValue, false))
}

func (this *SqlGenerator) OrWhereNotIn(field string, listValue [] interface{}) *SqlGenerator {
	if len(field) == 0 || len(listValue) == 0 {
		panic("WhereNotIn param error")
	}
	return this.whereResult( whereFactory.WhereNotIn(field, listValue, true))
}

func (this *SqlGenerator) WhereBetween(field string, interval [] interface{}) *SqlGenerator {
	if len(field) == 0 || len(interval) != 2 {
		panic("WhereBetween param error")
	}
	return this.whereResult( whereFactory.WhereBetween(field, interval, false))
}

func (this *SqlGenerator) OrWhereBetween(field string, interval [] interface{}) *SqlGenerator {
	if len(field) == 0 || len(interval) != 2 {
		panic("WhereBetween param error")
	}
	return this.whereResult( whereFactory.WhereBetween(field, interval, true))
}

func (this *SqlGenerator) WhereNotBetween(field string, interval [] interface{}) *SqlGenerator {
	if len(field) == 0 || len(interval) != 2 {
		panic("WhereNotBetween param error")
	}
	return this.whereResult( whereFactory.WhereNotBetween(field, interval, false))
}

func (this *SqlGenerator) OrWhereNotBetween(field string, interval [] interface{}) *SqlGenerator {
	if len(field) == 0 || len(interval) != 2 {
		panic("WhereNotBetween param error")
	}
	return this.whereResult( whereFactory.WhereNotBetween(field, interval, true))
}

func (this *SqlGenerator) WhereNull(field string) *SqlGenerator {
	if len(field) == 0 {
		panic("WhereNull param error")
	}
	return this.whereResult( whereFactory.WhereNull(field, false))
}

func (this *SqlGenerator) OrWhereNull(field string) *SqlGenerator {
	if len(field) == 0 {
		panic("WhereNull param error")
	}
	return this.whereResult( whereFactory.WhereNull(field, true))
}

func (this *SqlGenerator) WhereNotNull(field string) *SqlGenerator {
	if len(field) == 0 {
		panic("WhereNotNull param error")
	}
	return this.whereResult( whereFactory.WhereNotNull(field, false))
}

func (this *SqlGenerator) OrWhereNotNull(field string) *SqlGenerator {
	if len(field) == 0 {
		panic("WhereNotNull param error")
	}
	return this.whereResult( whereFactory.WhereNotNull(field, true))
}

func (this *SqlGenerator) WhereDate(field, date string) *SqlGenerator {
	if len(field) == 0 || len(date) == 0 {
		panic("WhereDate param error")
	}
	return this.whereResult( whereFactory.WhereDate(field, date, false))
}

func (this *SqlGenerator) OrWhereDate(field, date string) *SqlGenerator {
	if len(field) == 0 || len(date) == 0 {
		panic("WhereDate param error")
	}
	return this.whereResult( whereFactory.WhereDate(field, date, true))
}

func (this *SqlGenerator) WhereMonth(field string, month int) *SqlGenerator {
	if len(field) == 0 || month == 0 {
		panic("WhereMonth param error")
	}
	return this.whereResult( whereFactory.WhereMonth(field, month, false))
}

func (this *SqlGenerator) OrWhereMonth(field string, month int) *SqlGenerator {
	if len(field) == 0 || month == 0 {
		panic("WhereMonth param error")
	}
	return this.whereResult( whereFactory.WhereMonth(field, month, true))
}

func (this *SqlGenerator) WhereDay(field string, day int) *SqlGenerator {
	if len(field) == 0 || day == 0 {
		panic("WhereDay param error")
	}
	return this.whereResult( whereFactory.WhereDay(field, day, false))
}

func (this *SqlGenerator) OrWhereDay(field string, day int) *SqlGenerator {
	if len(field) == 0 || day == 0 {
		panic("WhereDay param error")
	}
	return this.whereResult( whereFactory.WhereDay(field, day, true))
}

func (this *SqlGenerator) WhereYear(field string, year int) *SqlGenerator {
	if len(field) == 0 || year == 0 {
		panic("whereYear param error")
	}
	return this.whereResult( whereFactory.whereYear(field, year, false))
}

func (this *SqlGenerator) OrWhereYear(field string, year int) *SqlGenerator {
	if len(field) == 0 || year == 0 {
		panic("whereYear param error")
	}
	return this.whereResult( whereFactory.whereYear(field, year, true))
}

func (this *SqlGenerator) WhereTime(field, condition, timestamp string) *SqlGenerator {
	if len(field) == 0 || len(condition) == 0 || len(timestamp) == 0 {
		panic("WhereTime param error")
	}
	return this.whereResult( whereFactory.WhereTime(field, condition, timestamp, false))
}

func (this *SqlGenerator) OrWhereTime(field, condition, timestamp string) *SqlGenerator {
	if len(field) == 0 || len(condition) == 0 || len(timestamp) == 0 {
		panic("WhereTime param error")
	}
	return this.whereResult( whereFactory.WhereTime(field, condition, timestamp, true))
}

func (this *SqlGenerator) WhereFunc(callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	return this.whereResult( whereFactory.WhereFunc(callback, this.NewGenerator(), false))
}

func (this *SqlGenerator) OrWhereFunc(callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	return this.whereResult( whereFactory.WhereFunc(callback, this.NewGenerator(), true))
}

func (this *SqlGenerator) WhereExists(callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	return this.whereResult( whereFactory.WhereExists("Exists", callback, false))
}

func (this *SqlGenerator) OrWhereExists(callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	return this.whereResult( whereFactory.WhereExists("Exists", callback, true))
}

func (this *SqlGenerator) WhereNotExists(field string, callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	return this.whereResult( whereFactory.WhereExists("Not Exists", callback, false))
}

func (this *SqlGenerator) OrWhereNotExists(field string, callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	return this.whereResult( whereFactory.WhereExists("Not Exists", callback, true))
}

func (this *SqlGenerator) WhereRaw(sql string) *SqlGenerator {
	return this.whereResult( whereFactory.WhereRaw(sql, false))
}

func (this *SqlGenerator) OrWhereRaw(sql string) *SqlGenerator {
	return this.whereResult( whereFactory.WhereRaw(sql, true))
}

func (this *SqlGenerator) When(boolean bool, callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	if boolean {
		return this.whereResult( whereFactory.WhereFunc(callback, this.NewGenerator(), false))
	}

	return this
}

func (this *SqlGenerator) OrWhen(boolean bool, callback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	if boolean {
		return this.whereResult( whereFactory.WhereFunc(callback, this.NewGenerator(), true))
	}

	return this
}

func (this *SqlGenerator) WhenElse(boolean bool, trueCallback func(build *SqlGenerator) *SqlGenerator, falseCallback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	if boolean {
		return this.whereResult( whereFactory.WhereFunc(trueCallback, this.NewGenerator(), false))
	} else {
		return this.whereResult( whereFactory.WhereFunc(falseCallback, this.NewGenerator(), false))
	}
}

func (this *SqlGenerator) OrWhenElse(boolean bool, trueCallback func(build *SqlGenerator) *SqlGenerator, falseCallback func(build *SqlGenerator) *SqlGenerator) *SqlGenerator {
	if boolean {
		return this.whereResult( whereFactory.WhereFunc(trueCallback, this.NewGenerator(), true))
	} else {
		return this.whereResult( whereFactory.WhereFunc(falseCallback, this.NewGenerator(), true))
	}
}

/*************************************  WHERE  END  *******************************************************/

/*************************************  Other  *******************************************************/
func (this *SqlGenerator) OrderBy(args... string) *SqlGenerator {
	if len(args) == 0 {
		panic("OrderBy param error")
	}
	for _, item := range args{
		this.sqlPart.orderData = append(this.sqlPart.orderData, item)
	}
	return this
}

func (this *SqlGenerator) OrderByRaw(sql string) *SqlGenerator {
	if len(sql) == 0 {
		panic("OrderByRaw param error")
	}
	this.sqlPart.orderData = append(this.sqlPart.orderData, sql)
	return this
}

func (this *SqlGenerator) GroupBy(args... string) *SqlGenerator {
	if len(args) == 0 {
		panic("GroupBy param error")
	}
	for _, item := range args{
		this.sqlPart.groupData = append(this.sqlPart.groupData, item)
	}
	return this
}

func (this *SqlGenerator) GroupByRaw(sql string) *SqlGenerator {
	if len(sql) == 0 {
		panic("GroupByRaw param error")
	}
	this.sqlPart.groupData = append(this.sqlPart.groupData, sql)
	return this
}

// TODO BUG SELECT中存在别名字段时无法识别表名加字段名
func (this *SqlGenerator) Having(args... interface{}) *SqlGenerator {
	if len(args) == 0 {
		panic("Having param error")
	}
	args = append(args, 1)
	return this.havingResult( whereFactory.Where(FormatWhereParam(args...)...))
}

func (this *SqlGenerator) HavingRaw(sql string) *SqlGenerator {
	if len(sql) == 0 {
		panic("HavingRaw param error")
	}
	return this.havingResult( whereFactory.WhereRaw(sql, false))
}

func (this *SqlGenerator) Offset(num int) *SqlGenerator {
	this.sqlPart.offset = num
	return this
}

func (this *SqlGenerator) Limit(num int) *SqlGenerator {
	if num == 0 {
		panic("Limit param error")
	}
	this.sqlPart.limit = num
	return this
}

/*************************************  Other END *******************************************************/

/*************************************  SELECT *******************************************************/

func (this *SqlGenerator) Get(args... string) *SqlGenerator {
	if len(args) > 0 {
		_ = this.Select(args...)
	}

	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlGenerator) Value(field string) *SqlGenerator {
	_ = this.Select(field)
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlGenerator) First() *SqlGenerator {
	_ = this.Limit(1)
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlGenerator) PluckArray(field string) *SqlGenerator {
	_ = this.Select(field)
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlGenerator) PluckMap(field, value string) *SqlGenerator {
	_ = this.Select(field, value)
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlGenerator) Chunk(num int, callback func()) *SqlGenerator {
	if num == 0 {
		panic("Chunk param error")
	}
	_ = this.Limit(num)
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlGenerator) Count() *SqlGenerator {
	_ = this.SelectRaw("COUNT(1) AS count")
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlGenerator) Max(field string) *SqlGenerator {
	if len(field) == 0 {
		panic("Max param error")
	}
	_ = this.SelectRaw("Max("+field+") AS max")
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlGenerator) Sum(field string) *SqlGenerator {
	if len(field) == 0 {
		panic("Sum param error")
	}
	_ = this.SelectRaw("Sum("+field+") AS sum")
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlGenerator) Exists() *SqlGenerator {
	_ = this.SelectRaw("COUNT(1) AS count")
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlGenerator) DoesntExists() *SqlGenerator {
	_ = this.SelectRaw("COUNT(1) AS count")
	this.exeType = "SELECT"
	return this.build()
}

func (this *SqlGenerator) ToSql() string {
	return this.ShowSql
}

/*************************************  SELECT END *******************************************************/

func (this *SqlGenerator) Insert(data map[string]interface{}) *SqlGenerator {
	if len(data) == 0 {
		panic("Insert param error")
	}

	this.sqlPart.insertData = append(this.sqlPart.insertData, data)
	this.exeType = "INSERT"
	return this.build()
}

func (this *SqlGenerator) MultiInsert(data []map[string]interface{}) *SqlGenerator {
	if len(data) == 0 {
		panic("MultiInsert param error")
	}
	this.sqlPart.insertData = data
	this.exeType = "INSERT"
	return this.build()
}

func (this *SqlGenerator) Update(data map[string]interface{}) *SqlGenerator {
	if len(data) == 0 {
		panic("Update param error")
	}
	this.sqlPart.updateData = data
	this.exeType = "UPDATE"
	return this.build()
}

func (this *SqlGenerator) Delete() *SqlGenerator {
	this.exeType = "DELETE"
	return this.build()
}

