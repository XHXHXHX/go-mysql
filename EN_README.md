### Go-Mysql
---
one of the Mysql ORM for Golang, it consists of simple client pool, simple Mysql results and simple generators of sql as similar as `Laravel` DB. 

The three parts could work independently, it's fixible.
    
    
    
### Installation
---
Simple install the package to you `$GOPATH` with `go tool` from shell:
    
    $ go get github.com/XHXHXHX/go-mysql
    
Make sure Git is installed on your machine and in your system's `PATH`

### Usage
---

        
- simple select

        import _ "github.com/XHXHXHX/go-mysql"
    
        InitConfig()
        SetPrefix("my_")
        
        res, err := DB().Table("user").Where("age", ">", 0).OrderBy("age").GroupBy("id").Offset(0).Limit(10).Get()
        ......
        
- complex query

        res, err = DB().Table("user", "u").WhenElse(res == nil, func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator {
            return build.JoinFunc("class", func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator {
                return build.On("id", "=", "uid").WhereNull("teacher")
            })
        }, func(build *sql_generators.SqlGenerator) *sql_generators.SqlGenerator {
            return build.Where("sex", 1).OrderBy("age", "asc")
        }).Get()
        
        
- show sql for execute

        res, err := DB().Table("user").Where("age", ">", 0).OrderBy("age").GroupBy("id").Offset(0).Limit(10).Get()
        fmt.Println(res.Generator.ShowSql)
        
- query result for struct

        DB().Table("user").Model(model interface{}) error
        DB().Table("user").Models(model interface{}) error
        
        
- Result struct. I think you need

        type Result struct {
        	*sql.Rows
        	Set []map[string]interface{}
        	LastInsertId int64
        	RowsAffected int64
        	FuncResult int
        }