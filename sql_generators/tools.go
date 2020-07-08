package sql_generators

import (
	"strconv"
	"strings"
)

/*
 * 格式化 where 参数
 */
func FormatWhereParam(args... interface{}) []interface{} {
	var new_args []interface{}
	l := len(args)
	for index, value := range args {
		switch val := value.(type) {
		case string:
			if l - 1 == index {
				new_args = append(new_args, "'" + val + "'")
			} else {
				new_args = append(new_args, val)
			}
		case int:
			new_args = append(new_args, val)
		default:
			new_args = append(new_args, val)
		}
	}

	return new_args
}

/*
 * 找到日期格式的连接符
 */
func GetDateStringJoiner(s string) string {
	for _, char := range s {
		if char < '0' || char > '9' {
			return string(char)
		}
	}

	return ""
}

/*
 * 返回 num 个占位符字符串
 */
func MakeInPlaceholderString(num int) string {
	return "(" + strings.Trim(strings.Repeat("?,", num), ",") + ")"
}
/*
 * 给字段或表名添加特殊引号
 */
func AddCharForString(s string) string {
	return "`" + s + "`"
}

/*
 * 给数组中的字段或表名添加特殊引号
 */
func AddCharForArray(arr []string) string {
	for index, val := range arr {
		arr[index] = AddCharForString(val)
	}
	return strings.Join(arr, ".")
}

/*
 * 用 char 连接符生成 Mysql 日期格式化字符串
 * %Y%m%d
 */
func SetMysqlDateFormatByChar(char string) string {
	return "%Y" + char + "%m" + char + "%d"
}

/*
 * 给字符串加引号
 */
func AddSingleSymbol(s string) string {
	return "'" + string(s) + "'"
}

/*
 * 转化成字符串
 */
func TransferString(s interface{}) string {
	switch r := s.(type) {
	case string:
		return r
	case int:
		return strconv.Itoa(r)
	default:
		panic("transfer string error")
	}
}


func InvalidOperator(option string) bool {
	for _, char := range options {
		if char == strings.ToLower(option) {
			return true
		}
	}
	return false
}