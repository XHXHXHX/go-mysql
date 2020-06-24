### Go-Mysql
---
由链接池/结果集/sql生成器组成 `Laravel` 风格的ORM，每个部分都可以独立使用。
    
    
### 安装
---

    $ go get github.com/XHXHXHX/go-mysql
    

### 使用
---

        
- 简易查询

        import _ "github.com/XHXHXHX/go-mysql/mysqlManager"
    
        InitConfig()
        SetPrefix("my_")
        
        res, err := DB().Table("user").Where("age", ">", 0).OrderBy("age").GroupBy("id").Offset(0).Limit(10).Get()
        ......
        
- 稍微复杂的查询

        res, err = DB().Table("user", "u").WhenElse(res == nil, func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator {
            return build.JoinFunc("class", func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator {
                return build.On("id", "=", "uid").WhereNull("teacher")
            })
        }, func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator {
            return build.Where("sex", 1).OrderBy("age", "asc")
        }).Get()
        
        
- 执行后查看执行sql

        res, err := DB().Table("user").Where("age", ">", 0).OrderBy("age").GroupBy("id").Offset(0).Limit(10).Get()
        fmt.Println(res.Generator.ShowSql)
        
- 也可以将查询结果绑定到 `struct`中

        DB().Table("user").Model(model interface{}) error
        DB().Table("user").Models(model interface{}) error
        
        
- 结果集结构体

        type Result struct {
        	*sql.Rows
        	Set []map[string]interface{}
        	LastInsertId int64
        	RowsAffected int64
        	FuncResult int
        }
