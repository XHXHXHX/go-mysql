package result

/*
 * 韦钊大佬写的
 */

import (
	"database/sql"
	"encoding/json"
)

type Result struct {
	*sql.Rows
	Set []map[string]interface{}
	LastInsertId int64
	RowsAffected int64
	FuncResult int			// 功能函数如 sum, count, max 的值
}

func (r *Result) MapResult(ptr interface{}) {
	if r == nil {
		return
	}

	if r.Set == nil {
		r.MakeResult()
	}

	//如果map后仍然是nil，说明没有值
	if r.Set == nil {
		return
	}

	setbytes, err := json.Marshal(r.Set[0])
	if err != nil {
		return
	}

	err = json.Unmarshal(setbytes, ptr)
	if err != nil {
		return
	}
}

func (r *Result) MapResults(ptr interface{}) {
	if r == nil {
		return
	}

	if r.Set == nil {
		r.MakeResult()
	}

	//如果map后仍然是nil，说明没有值
	if r.Set == nil {
		return
	}

	setbytes, err := json.Marshal(r.Set)
	if err != nil {
		return
	}

	err = json.Unmarshal(setbytes, ptr)
	if err != nil {
		return
	}
}

func (r *Result) GetColumnValue(fieldName string, value interface{}, MessageRowMap map[string]interface{}) map[string]interface{} {
	//fmt.Printf("%v %s %s, ", fieldName, value, reflect.TypeOf(value))
	switch v := value.(type) {
	case int:
		MessageRowMap[fieldName] = int64(v)
	case int8:
		MessageRowMap[fieldName] = int64(v)
	case int16:
		MessageRowMap[fieldName] = int64(v)
	case int32:
		MessageRowMap[fieldName] = int64(v)
	case int64:
		MessageRowMap[fieldName] = int64(v)
	case uint:
		MessageRowMap[fieldName] = uint64(v)
	case uint8:
		MessageRowMap[fieldName] = uint64(v)
	case uint16:
		MessageRowMap[fieldName] = uint64(v)
	case uint32:
		MessageRowMap[fieldName] = uint64(v)
	case uint64:
		MessageRowMap[fieldName] = uint64(v)
	case float32:
		MessageRowMap[fieldName] = float64(v)
	case float64:
		MessageRowMap[fieldName] = float64(v)
	case string:
		MessageRowMap[fieldName] = string(v)
	case []byte:
		MessageRowMap[fieldName] = string(v)
	case nil:
		MessageRowMap[fieldName] = nil
	case bool:
		MessageRowMap[fieldName] = bool(v)
	default:
		MessageRowMap[fieldName] = nil
	}

	//fmt.Printf("%T", MessageRowMap[fieldName])
	//fmt.Println()
	return MessageRowMap
}


func (r *Result) MakeResult() {
	columns, err := r.Rows.Columns()
	if err != nil {
		return
	}

	l := len(columns)
	values := make([]interface{}, l)
	valueCollect := make([]interface{}, l)

	for i := 0; i < l; i++ {
		valueCollect[i] = &values[i]
	}

	for r.Rows.Next() {
		tmp := make(map[string]interface{})
		_ = r.Rows.Scan(valueCollect...)
		for i, name := range columns {
			tmp = r.GetColumnValue(name, values[i], tmp)
			//tmp[name] = values[i]
		}
		//fmt.Println(fmt.Sprintf("%T", tmp["user_name"]))

		r.Set = append(r.Set, tmp)
	}
	return
}