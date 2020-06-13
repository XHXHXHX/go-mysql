package sql_generators

import (
	"fmt"
	"strings"
	"testing"
)

const (
	TextHeaderNum = 20
)

func TestDB(t *testing.T) {
	result := DB().Table("class", "c").Where("a", 1).OrWhere("b", "2").Get()
	myPrintln("Where", result.ShowSql)
}

func TestSqlGenerator_WhereIn(t *testing.T) {
	var arr []interface{}
	arr = append(arr, 1)
	arr = append(arr, "2")
	arr = append(arr, 4)
	arr = append(arr, 5)

	result := DB().Table("class").Where("a", "1").WhereIn("type", arr).Get()
	myPrintln("WhereIn", result.ShowSql)
}

func TestSqlGenerator_WhereNull(t *testing.T) {
	result := DB().Table("goods").WhereNull("exrat_data").Get()
	myPrintln("WhereNull", result.ShowSql)
}

func TestSqlGenerator_WhereBetween(t *testing.T) {
	var arr []interface{}
	arr = append(arr, 123)
	arr = append(arr, "456")

	result := DB().Table("class").WhereBetween("type", arr).Get()
	myPrintln("WhereBetween", result.ShowSql)
}

func TestSqlGenerator_WhereArray(t *testing.T) {
	var arr = [][]interface{} {
		{
			"id", 2,
		},{
			"is_del", "!=", 0,
		},
	}

	result := DB().Table("class").WhereArray(arr).Get()
	myPrintln("WhereArray", result.ShowSql)
}

func TestSqlGenerator_WhereMap(t *testing.T) {
	myMap := make(map[string]interface{})
	myMap["id"] = 2
	myMap["name"] = 1
	myMap["age"] = "23"

	result := DB().Table("class").WhereMap(myMap).Get()
	myPrintln("WhereMap", result.ShowSql)
}

func TestSqlGenerator_WhereDate(t *testing.T) {
	result := DB().Table("goods").WhereDate("add_time", "2019-09-21").Get()
	myPrintln("WhereDate", result.ShowSql)
}

func TestSqlGenerator_WhereMonth(t *testing.T) {
	result := DB().Table("goods").WhereMonth("add_time", "10").Get()
	myPrintln("WhereMonth", result.ShowSql)
}

func TestSqlGenerator_WhereDay(t *testing.T) {
	result := DB().Table("goods").WhereDay("add_time", "31").Get()
	myPrintln("WhereDay", result.ShowSql)
}

func TestSqlGenerator_WhereYear(t *testing.T) {
	result := DB().Table("goods").WhereYear("add_time", "2020").Get()
	myPrintln("WhereYear", result.ShowSql)
}

func TestSqlGenerator_WhereTime(t *testing.T) {
	result := DB().Table("goods").WhereTime("add_time", "<", "13:20:11").Get()
	myPrintln("WhereTime", result.ShowSql)
}

func TestSqlGenerator_WhereFunc(t *testing.T) {
	var aaa = 1
	result := DB().Table("goods", "g").WhereFunc(func(query *SqlGenerator) *SqlGenerator {
		aaa += 2
		var arr []interface{}
		arr = append(arr, aaa)
		arr = append(arr, "99")
		return query.Where("goods_stock", 1).WhereBetween("goods_stock", arr)
	}).Where("id", 23).Get()

	myPrintln("WhereFunc", result.ShowSql)
}

func TestSqlGenerator_WhereRaw(t *testing.T) {
	result := DB().Table("goods").WhereRaw("goods_name = '花卷'").Get()
	myPrintln("WhereRaw", result.ShowSql)
}

func TestSqlGenerator_GroupBy(t *testing.T) {
	result := DB().Table("goods").GroupBy("goods_type").Get()
	myPrintln("GroupBy", result.ShowSql)
}

func TestSqlGenerator_OrderBy(t *testing.T) {
	result := DB().Table("goods").OrderBy("id desc", "add_time asc").Get()
	myPrintln("OrderBy", result.ShowSql)
}

func TestSqlGenerator_Join(t *testing.T) {
	result := DB().Table("goods", "g").Join("goods_score", "goods_id", "=", "goods_id").Get()
	myPrintln("Join", result.ShowSql)
}

func TestSqlGenerator_LeftJoin(t *testing.T) {
	result := DB().Table("goods", "g").LeftJoin("goods_score as s", "goods_id", "=", "goods_id").Get()
	myPrintln("LeftJoin", result.ShowSql)
}

func TestSqlGenerator_LeftJoinFunc(t *testing.T) {
	result := DB().Table("hj_goods", "g").LeftJoinFunc("hj_goods_score as s", func(build *SqlGenerator) *SqlGenerator {
		return build.On("goods_id", "=", "goods_id").Where("is_del", 0)
	}).Where("g.goods_stock", ">", 0).Get()
	myPrintln("LeftJoinFunc", result.ShowSql)
}

func TestSqlGenerator_When(t *testing.T) {
	a := 3
	b := 2
	result := DB().Table("hj_goods").When(a < b, func(build *SqlGenerator) *SqlGenerator {
		return build.Where("goods_state", 1)
	}).WhenElse(6 < 5, func(build *SqlGenerator) *SqlGenerator {
		return build.Where("is_show", 1)
	}, func(build *SqlGenerator) *SqlGenerator {
		return build.Where("goods_state", 1)
	}).Get()
	myPrintln("When", result.ShowSql)
}

func TestSqlGenerator_WhenElse(t *testing.T) {
	a := 3
	b := 2
	result := DB().Table("hj_goods").WhenElse(a < b, func(build *SqlGenerator) *SqlGenerator {
		return build.Where("is_show", 1)
	}, func(build *SqlGenerator) *SqlGenerator {
		return build.Where("goods_state", 1)
	}).Get()
	myPrintln("WhenElse", result.ShowSql)
}

func TestSqlGenerator_Insert(t *testing.T) {
	result := DB().Table("goods").Insert(map[string]interface{}{
		"name": "xiaoming",
		"age": 1,
		"hobby": "basketball",
	})
	myPrintln("Insert", result.ShowSql)
}

func TestSqlGenerator_MultiInsert(t *testing.T) {
	result := DB().Table("goods").MultiInsert([]map[string]interface{}{
		map[string]interface{}{
			"name":  "xiaoming",
			"age":   1,
			"hobby": "basketball",
		},map[string]interface{}{
			"name":  "小黄",
			"age":   2,
			"hobby": "football",
		}})
	myPrintln("MultiInsert", result.ShowSql)
}

func TestSqlGenerator_Update(t *testing.T) {
	result := DB().Table("goods").Where("name", "Uzi").Update(map[string]interface{}{
		"name": "legend",
		"age": 1,
		"hobby": "basketball",
	})
	myPrintln("Update", result.ShowSql)
}

func TestSqlGenerator_Delete(t *testing.T) {
	result := DB().Table("goods").Where("name", "Uzi").Delete()
	myPrintln("Delete", result.ShowSql)
}


func myPrintln(s, sql string) {
	l := len(s)
	fmt.Println(s, strings.Repeat(" ", TextHeaderNum - l), sql)
}