package generators

import (
	"strconv"
	"strings"
)

type sqlInfo struct {
	table string
	alias string
	joinData [] *joinInfo
	whereData [] *whereInfo
	groupData []string
	orderData []string
	havingData [] *whereInfo
	offset	int
	limit 	int
	selectData string
	insertData []map[string]interface{}
	insertField []string
	updateData map[string]interface{}
	sql string
	preproParam []interface{}
}

func (this *sqlInfo) setTable(table, alias string) {
	this.table = table
	if len(alias) == 0 {
		this.alias = ""
	} else {
		this.alias = alias
	}
}

func (this *sqlInfo) setJoin(buildInfo *SqlGenerator, joinType, thisRelationField, relationCondition, thatRelationField string) {
	joinData := &joinInfo{
		buildInfo: buildInfo,
		joinType: joinType,
		thisRelationField: thisRelationField,
		relationCondition: relationCondition,
		thatRelationField: thatRelationField,
	}
	this.joinData = append(this.joinData, joinData)
}

func (this *sqlInfo) setWhere(whereData *whereInfo) {
	this.whereData = append(this.whereData, whereData)
}

func (this *sqlInfo) BuildInsert() (string, []interface{}) {
	var sql_array [] string
	sql_array = append(sql_array, "INSERT INTO")
	sql_array = append(sql_array, this.BuildTable(true))
	sql_array = append(sql_array, this.BuildInsertField())
	sql_array = append(sql_array, this.BuildInsertValue())

	return strings.Join(sql_array, " "), this.preproParam
}

func (this *sqlInfo) BuildUpdate() (string, []interface{}) {
	var sql_array [] string
	sql_array = append(sql_array, "UPDATE")
	sql_array = append(sql_array, this.BuildTable(true))
	sql_array = append(sql_array, "SET")
	sql_array = append(sql_array, this.BuildUpdateData())
	sql_array = append(sql_array, this.BuildWhere(false))

	return strings.Join(sql_array, " "), this.preproParam
}

func (this *sqlInfo) BuildDelete() (string, []interface{}) {
	var sql_array [] string
	sql_array = append(sql_array, "DELETE")
	sql_array = append(sql_array, this.BuildTable(false))
	sql_array = append(sql_array, this.BuildWhere(false))

	return strings.Join(sql_array, " "), this.preproParam
}

func (this *sqlInfo) BuildQuery() (string, []interface{}) {
	var sql_array [] string
	sql_array = append(sql_array, "SELECT")
	sql_array = append(sql_array, this.selectData)
	sql_array = append(sql_array, this.BuildTable(false))
	sql_array = append(sql_array, this.BuildJoin())
	sql_array = append(sql_array, this.BuildWhere(false))
	sql_array = append(sql_array, this.BuildGroup())
	sql_array = append(sql_array, this.BuildOrder())
	sql_array = append(sql_array, this.BuildLimit())
	sql_array = append(sql_array, this.BuildWhere(true))

	var new_array []string
	for _, item := range sql_array {
		if item != "" {
			new_array = append(new_array, item)
		}
	}

	return strings.Join(new_array, " "), this.preproParam
}

func (this *sqlInfo) BuildUpdateData() string {
	var s []string
	for field, value := range this.updateData {
		if val, ok := value.(string);ok && strings.Index(val, "&/") == 0 {
			s = append(s, field + " = " + val)
		} else {
			s = append(s, field + " = ?")
			this.preproParam = append(this.preproParam, value)
		}
	}

	return strings.Join(s, ",")
}

func (this *sqlInfo) BuildInsertField() string {
	for field, _ := range this.insertData[0] {
		this.insertField = append(this.insertField, field)
	}

	return "(" + strings.Join(this.insertField, ",") + ")"
}

func (this *sqlInfo) BuildInsertValue() string {
	var s []string
	l := len(this.insertField)
	for _, value := range this.insertData {
		for _, field := range this.insertField {
			this.preproParam = append(this.preproParam, value[field])
		}
		s = append(s, "(" + strings.Trim(strings.Repeat("?,", l), ",") + ")")
	}

	return "VALUES" + strings.Join(s, ",")
}

func (this *sqlInfo) BuildLimit() string {
	if this.limit == 0 {
		return ""
	}

	return "Limit " + strconv.Itoa(this.offset) + ", " + strconv.Itoa(this.limit)
}

func (this *sqlInfo) BuildGroup() string {
	if len(this.groupData) == 0 {
		return ""
	}
	return "GROUP BY " + strings.Join(this.groupData, ",")
}

func (this *sqlInfo) BuildOrder() string {
	if len(this.orderData) == 0 {
		return ""
	}
	return "ORDER BY " + strings.Join(this.orderData, ",")
}

func (this *sqlInfo) BuildTable(notFrom bool) string {
	var sql []string
	if notFrom {
		sql = []string {AddCharForString(this.table)}
	} else {
		sql = []string {"FROM", AddCharForString(this.table)}
	}

	if len(this.alias) > 0 {
		sql = append(sql, "AS")
		sql = append(sql, this.alias)
	}

	return strings.Join(sql, " ")
}

func (this *sqlInfo) BuildJoin() string {
	if len(this.joinData) == 0 {
		return ""
	}

	var join []string
	for _, value := range this.joinData {
		join = append(join, this.makeJoinString(value))
	}

	return strings.Join(join, " ")
}

func (this *sqlInfo) makeJoinString(data *joinInfo) string {
	var s []string
	s = append(s, strings.ToUpper(data.joinType))
	s = append(s, data.buildInfo.sqlPart.BuildTable(true))
	if len(data.thisRelationField) > 0 {
		thatRelationField := data.buildInfo.sqlPart.addTablePrefixForField(data.thatRelationField)
		thisRelationField := this.addTablePrefixForField(data.thisRelationField)
		s = append(s, strings.Join([]string{"ON", thatRelationField, data.relationCondition, thisRelationField}, " "))
	}
	if len(data.buildInfo.sqlPart.whereData) > 0 {
		join_where := data.buildInfo.sqlPart.BuildWhere(false)
		s = append(s, strings.Replace(join_where, "WHERE", "AND", 1))
		for _, value := range data.buildInfo.sqlPart.preproParam {
			this.preproParam = append(this.preproParam, value)
		}
	}

	return strings.Join(s, " ")
}

func (this *sqlInfo) BuildWhere(isHaving bool) string {
	var data [] *whereInfo
	var where [] string
	if isHaving {
		data = this.havingData
		where = append(where, "HAVING")
	} else {
		data = this.whereData
		where = append(where, "WHERE")
	}
	if len(data) == 0 {
		return ""
	}
	for _, item := range data {
		var s string
		 switch item.whereType{
			 case "Basic":
				 s = this.BasicWhere(item)
			 case "In":
				 s = this.InWhere(item)
			 case "Between":
				 s = this.BetweenWhere(item)
			 case "Null":
				 s = this.NullWhere(item)
			 case "Func":
				 s = this.FuncWhere(item)
			 case "Raw":
			 	value, _ := item.value.(string)
				 s = value
			 case "Nested":
				 s = item.build.sqlPart.BuildWhere(false)
				 s = strings.TrimPrefix(s, "WHERE")
				 for _, val := range item.build.sqlPart.preproParam {
				 	this.preproParam = append(this.preproParam, val)
				 }
			 case "Exists":
				 s = item.build.sqlPart.BuildTable(false)
			 default:
		 }

		where = append(where, this.LittleChangeForWhere(item.boolean, item.build, s))
	}

	return this.dealHeaderWhere(where)
}

func (this *sqlInfo) BasicWhere(data *whereInfo) string {
	sql := this.SetWhereCommonPart(data.field, data.conditionType)
	sql = append(sql, "?")
	this.preproParam = append(this.preproParam, data.value)

	return strings.Join(sql, " ")
}

func (this *sqlInfo) InWhere(data *whereInfo) string {
	sql := this.SetWhereCommonPart(data.field, data.conditionType)

	value, _ := data.value.([]interface{})
	sql = append(sql, MakeInPlaceholderString(len(value)))

	for _, item := range value {
		this.preproParam = append(this.preproParam, item)
	}

	return strings.Join(sql, " ")
}

func (this *sqlInfo) BetweenWhere(data *whereInfo) string {
	sql := this.SetWhereCommonPart(data.field, data.conditionType)
	sql = append(sql, "? AND ?")

	value, _ := data.value.([]interface{})
	for _, item := range value {
		this.preproParam = append(this.preproParam,item)
	}

	return strings.Join(sql, " ")
}

func (this *sqlInfo) NullWhere(data *whereInfo) string {
	sql := this.SetWhereCommonPart(data.field, data.conditionType)
	sql = append(sql, "Null")

	return strings.Join(sql, " ")
}

func (this *sqlInfo) FuncWhere(data *whereInfo) string {
	var sql []string
	sql = append(sql, data.field)
	sql = append(sql, data.conditionType)
	sql = append(sql, "?")
	value, _ := data.value.(string)
	this.preproParam = append(this.preproParam, value)

	return strings.Join(sql, " ")
}

func (this *sqlInfo) dealHeaderWhere(where []string) string {
	where[1] = strings.TrimPrefix(where[1], "AND ")
	where[1] = strings.TrimPrefix(where[1], "OR ")

	if len(where) == 2 {
		where[1] = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(where[1], "("), ")"))
	}

	return strings.Join(where, " ")
}

func (this *sqlInfo) LittleChangeForWhere(boolean bool, subQuery *SqlGenerator, s string) string {
	c := "AND "
	if boolean {
		c = "OR "
	}

	if subQuery == nil {
		return c + s
	}

	return c + "(" + s + ")"
}

func (this *sqlInfo) addTablePrefixForField(field string) string {
	if strings.Count(field, ".") > 0 {
		arr := strings.Split(field, ".")
		return AddCharForArray(arr)
	}

	if this.alias == "" {
		return AddCharForString(this.table) + "." + AddCharForString(field)
	} else {
		return AddCharForString(this.alias) + "." + AddCharForString(field)
	}

}

func (this *sqlInfo) SetWhereCommonPart(field, conditionType string) []string {
	var sql []string
	sql = append(sql, this.addTablePrefixForField(field))
	if len(conditionType) > 0 {
		sql = append(sql, conditionType)
	}

	return sql
}