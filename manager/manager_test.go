package mysqlManager

import (
	"go-mysql/sql_generators"
	"testing"
)

type User struct {
	Id int64	`json:"id"`
	UserName string 	`json:"user_name"`
	Age int		`json:"age"`
	Sex uint8	`json:"sex"`
	LastUpdateTime string `json:"last_update_time"`
}

var userModel *User
var userModels []*User
var isInit bool

func TestInitConfig(t *testing.T) {
	SetPrefix("my_")
	err := InitConfig("../config.json")
	if err != nil {
		t.Error(err)
	}
	isInit = true
}

func TestManager_InsertToSql(t *testing.T) {
	// INSERT INTO `my_user` (user_name,age,sex) VALUES('小狼人',1,1)
	sql := DB().Table("user").InsertToSql(map[string]interface{}{
		"user_name": "小狼人",
		"age": 1,
		"sex": 1,
	})
	if sql != "INSERT INTO `my_user` (user_name,age,sex) VALUES('小狼人',1,1)" {
		//t.Error(sql)
	}
}
func TestManager_Insert(t *testing.T) {
	// 1
	res, err := DB().Table("user").Insert(map[string]interface{}{
		"user_name": "小狼人",
		"age": 1,
		"sex": 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.FuncResult == 1 {
		t.Error(res.FuncResult)
	}
}
func TestManager_MultiInsertToSql(t *testing.T) {
	// INSERT INTO `my_user` (user_name,age,sex) VALUES('小红帽',2,0),('小灰帽',3,1),('小黄帽',4,1)
	sql := DB().Table("user").MultiInsertToSql([]map[string]interface{}{
		{"user_name":"小红帽", "age": 2, "sex":0},
		{"user_name":"小灰帽", "age": 3, "sex":1},
		{"user_name":"小黄帽", "age": 4, "sex":1},
	})
	if sql != "INSERT INTO `my_user` (age,sex,user_name) VALUES(2,0,'小红帽'),(3,1,'小灰帽'),(4,1,'小黄帽')" {
		//t.Error(sql)
	}
}
func TestManager_MultiInsert(t *testing.T) {
	if isInit == false {
		TestInitConfig(t)
	}
	// 3
	res, err := DB().Table("user").MultiInsert([]map[string]interface{}{
		{"user_name":"小红帽", "age": 2, "sex":0},
		{"user_name":"小灰帽", "age": 3, "sex":1},
		{"user_name":"小黄帽", "age": 4, "sex":1},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.RowsAffected != 3 {
		t.Error(res.RowsAffected)
	}
}
func TestManager_LastInsertIdToSql(t *testing.T) {
	// INSERT INTO `my_user` (user_name,age,sex) VALUES('小红帽',1,0)
	sql := DB().Table("user").LastInsertIdToSql(map[string]interface{}{
		"user_name": "小红帽",
		"age": 1,
		"sex": 0,
	})
	if sql != "INSERT INTO `my_user` (user_name,age,sex) VALUES('小红帽',1,0)" {
		//t.Error(sql)
	}
}
func TestManager_GetLastInsertId(t *testing.T) {
	// 1
	res, err := DB().Table("user").GetLastInsertId(map[string]interface{}{
		"user_name": "小红帽",
		"age": 1,
		"sex": 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.LastInsertId == 0 {
		t.Error(res.LastInsertId)
	}
}
func TestManager_UpdateToSql(t *testing.T) {
	// UPDATE `my_user` SET user_name = '小绿帽',sex = 1 WHERE `my_user`.`user_name` = '小红帽'
	sql := DB().Table("user").Where("user_name", "小红帽").UpdateToSql(map[string]interface{}{
		"user_name": "小绿帽",
		"sex": 1,
	})
	if sql != "UPDATE `my_user` SET user_name = '小绿帽',sex = 1 WHERE `my_user`.`user_name` = '小红帽'" {
		//t.Error("UpdateToSql error")
	}
}
func TestManager_Update(t *testing.T) {
	// 2
	res, err := DB().Table("user").Where("user_name", "小红帽").Update(map[string]interface{}{
		"user_name": "小绿帽",
		"sex": 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.RowsAffected == 0 {
		t.Error(res)
	}
}
func TestManager_DeleteToSql(t *testing.T) {
	// DELETE FROM `my_user` WHERE `my_user`.`user_name` = '小狼人'
	sql := DB().Table("user").Where("user_name", "小狼人").DeleteToSql()
	if sql != "DELETE FROM `my_user` WHERE `my_user`.`user_name` = '小狼人'" {
		t.Error(sql)
	}
}
func TestManager_Delete(t *testing.T) {
	// 1
	res, err := DB().Table("user").Where("user_name", "小狼人").Delete()
	if err != nil {
		t.Fatal(err)
	}
	if res.RowsAffected != 1 {

	}
}
func TestManager_Get(t *testing.T) {
	// &{0xc0000ae180 [map[age:1 id:1 last_update_time:2020-06-14 10:12:52 sex:1 user_name:小明] map[age:1 id:4 last_update_time:2020-06-14 10:13:33 sex:2 user_name:小蓝]]}
	res, err := DB().Table("user").Where("age", ">", 0).OrderBy("age").GroupBy("id").Offset(0).Limit(10).Get()
	if err != nil {
		t.Fatal(err)
	}
	if res.Set == nil {
		t.Error("result MakeResult failed")
	}
}

func TestManager_First(t *testing.T) {
	// map[age:2 id:2 last_update_time:2020-06-14 10:13:13 sex:2 user_name:小红]
	var res interface{}
	res, err := DB().Table("user").Where("age", ">", 1).First()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.(map[string]interface{});!ok {
		t.Error("First error")
	}
}
func TestManager_Value(t *testing.T) {
	// 小明
	res, err := DB().Table("user").Where("id", 1).Value("user_name")
	if err != nil {
		t.Fatal(err)
	}
	if res != "小樱" {
		t.Error(res)
	}
}
func TestManager_Model(t *testing.T) {
	// &{2 小红 2 2 2020-06-14 10:13:13}
	err := DB().Table("user").Where("age", ">", 1).Model(&userModel)
	if err != nil {
		t.Fatal(err)
	}
	if userModel == nil {
		t.Error("Model error")
	}
}
func TestManager_Models(t *testing.T) {
	// [0xc000192840 0xc000192880]
	err := DB().Table("user").Where("age", ">", 1).Models(&userModels)
	if err != nil {
		t.Fatal(err)
	}
	if userModels == nil {
		t.Error("Models error")
	}
}
func TestManager_Count(t *testing.T) {
	// 3
	res, err := DB().Table("user").Where("age", ">", 1).Count()
	if err != nil {
		t.Fatal(err)
	}
	if res.FuncResult != 3 {
		//t.Error(res.FuncResult)
	}
}
func TestManager_Max(t *testing.T) {
	// 4
	res, err := DB().Table("user").Where("age", ">", 1).Max("age")
	if err != nil {
		t.Fatal(err)
	}
	if res.FuncResult != 4 {
		//t.Error(res.FuncResult)
	}
}
func TestManager_Sum(t *testing.T) {
	// 9
	res, err := DB().Table("user").Where("age", ">", 1).Sum("age")
	if err != nil {
		t.Fatal(err)
	}
	if res.FuncResult != 9 {
		//t.Error(res.FuncResult)
	}
}
func TestManager_PluckArray(t *testing.T) {
	// [2 3 4]
	res, err := DB().Table("user").Where("age", ">", 1).PluckArray("age")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 3 {
		//t.Error(res)
	}
}
func TestManager_PluckMap(t *testing.T) {
	// map[4:小绿帽 5:小灰帽 6:小黄帽]
	res, err := DB().Table("user").Where("age", ">", 1).PluckMap("id", "user_name")
	if err != nil {
		t.Fatal(err)
	}
	if res["4"] != "小绿帽" {
		t.Error(res)
	}
}

func TestDB(t *testing.T) {
	if isInit == false {
		TestInitConfig(t)
	}
	var res *Manager

	// SELECT user_name,age FROM `my_user` AS u WHERE `u`.`age` Between 1 AND 4 AND `u`.`user_name` IN ('小黄帽','小明')
	res, err := DB().Table("user", "u").WhereBetween("age", []interface{}{1, 4}).WhereIn("user_name", []interface{}{"小黄帽", "小明"}).Get("user_name", "age")
	if err != nil {
		t.Log(res.Generator.ShowSql)
		t.Error(err)
	}

	// SELECT * FROM `my_user` AS u WHERE `u`.`age` = 1 AND `u`.`sex` = 0 AND `u`.`user_name` = '小樱'
	res, err = DB().Table("user", "u").WhereArray([][]interface{}{
		{"age", 1},
		{"sex", 0},
		{"user_name", "小樱"},
	}).Get()
	if err != nil {
		t.Log(res.Generator.ShowSql)
		t.Error(err)
	}

	// SELECT * FROM `my_user` AS u WHERE `u`.`age` = 1 AND `u`.`sex` = 1 AND `u`.`user_name` = '小明'
	res, err = DB().Table("user", "u").WhereMap(map[string]interface{}{
		"age": 1, "sex": 1, "user_name": "小明",
	}).Get()
	if err != nil {
		t.Log(res.Generator.ShowSql)
		t.Error(err)
	}

	// SELECT * FROM `my_user` AS u WHERE `u`.`user_name` IN ('小明') AND DATE_FORMAT(last_update_time, '%m') = 6
	res, err = DB().Table("user", "u").WhereIn("user_name", []interface{}{"小明"}).WhereMonth("last_update_time", 6).Get()
	if err != nil {
		t.Log(res.Generator.ShowSql)
		t.Error(err)
	}

	// SELECT * FROM `my_class` AS c WHERE `c`.`teacher` IS Null
	res, err = DB().Table("class", "c").WhereNull("teacher").Get()
	if err != nil {
		t.Log(res.Generator.ShowSql)
		t.Error(err)
	}

	// SELECT u.user_name, c.grade, c.class FROM `my_user` AS u INNER JOIN `my_class` AS c ON `u`.`id` = `c`.`uid` WHERE `u`.`age` < 3
	res, err = DB().Table("user", "u").Join("class as c", "id", "=", "uid").Where("age", "<", 3).Get("u.user_name, c.grade, c.class")
	if err != nil {
		t.Log(res.Generator.ShowSql)
		t.Error(err)
	}

	// SELECT * FROM `my_user` AS u INNER JOIN `my_class` ON `u`.`id` = `my_class`.`uid` AND `my_class`.`teacher` IS Null
	res, err = DB().Table("user", "u").JoinFunc("class", func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator {
		return build.On("id", "=", "uid").WhereNull("teacher")
	}).Get()
	if err != nil {
		t.Log(res.Generator.ShowSql)
		t.Error(err)
	}

	// SELECT * FROM `my_user` AS u WHERE `u`.`sex` = 0
	res, err = DB().Table("user", "u").When(res != nil, func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator {
		return build.Where("sex", 0)
	}).Get()
	if err != nil {
		t.Log(res.Generator.ShowSql)
		t.Error(err)
	}

	// SELECT * FROM `my_user` AS u WHERE `u`.`sex` = 1
	res, err = DB().Table("user", "u").WhenElse(res == nil, func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator {
		return build.JoinFunc("class", func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator {
			return build.On("id", "=", "uid").WhereNull("teacher")
		})
	}, func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator {
		return build.Where("sex", 1).OrderBy("age", "asc")
	}).Get()
	if err != nil {
		t.Log(res.Generator.ShowSql)
		t.Error(err)
	}

	// SELECT count(id) as count, age FROM `my_user` GROUP BY age HAVING count > 1
	res, err = DB().Table("user").GroupBy("age").SelectRaw("count(id) as count, age").HavingRaw("count > 1").Get()
	if err != nil {
		t.Log(res.Generator.ShowSql)
		t.Error(err)
	}
}

// TODO
func TestDeBug(t *testing.T) {
	if isInit == false {
		TestInitConfig(t)
	}
	// SELECT count(id) as count, age FROM `my_user` GROUP BY age HAVING `my_user`.`count` > 1
	res, err := DB().Table("user").GroupBy("age").SelectRaw("count(id) as count, age").Having("count", ">", 1).Get()
	if err != nil {
		t.Log(res.Generator.ShowSql)
		t.Error(err)
	}
}