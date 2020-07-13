package manager

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/XHXHXHX/go-mysql/pool"
	"github.com/XHXHXHX/go-mysql/result"
	"github.com/XHXHXHX/go-mysql/sql_generators"
	"strconv"
)

const DEFAULT_CONFIG = "config/mysql.json"

/*
 * Manager
 */
type Manager struct {
	Generator *sql_generators.SqlGenerator
	ExecuteSql string
	client    *sql.DB
	tx        *sql.Tx
	*result.Result
}

func InitConfig(config string) error {
	if len(config) == 0 {
		config = DEFAULT_CONFIG
	}
	err := pool.InitPool(config)
	if err != nil {
		return err
	}
	return nil
}

func DB() *Manager {
	return &Manager{Generator: sql_generators.DB()}
}

func DbBegin() *Manager {
	client, err := pool.GetClient()
	if err != nil {
		panic(err)
	}
	tx, err := client.Begin()
	if err != nil {
		panic(err)
	}
	return &Manager{
		Generator: sql_generators.NewGenerator(),
		client:    client,
		tx:        tx,
	}
}

func (manage *Manager) DbCommit() error {
	if manage.tx == nil {
		return errors.New("Please begin transaction by DbBegin() first")
	}
	if err := manage.tx.Commit(); err != nil {
		return err
	}
	err := pool.CloseClient(manage.client)
	if err != nil {
		return err
	}
	manage.tx = nil
	manage.client = nil
	return nil
}

func (manage *Manager) DbRollBack() error {
	if manage.tx == nil {
		return errors.New("Please begin transaction by DbBegin() first")
	}
	err := manage.tx.Rollback()
	if err != sql.ErrTxDone && err != nil {
		return err
	}
	err = pool.CloseClient(manage.client)
	if err != nil {
		return err
	}
	manage.tx = nil
	manage.client = nil
	return nil
}

func (manage *Manager) GetLastInsertId(data map[string]interface{}) (*Manager, error) {
	manage.Generator.Insert(data)
	return manage.exec(true)
}

func (manage *Manager) Insert(data map[string]interface{}) (*Manager, error) {
	manage.Generator.Insert(data)
	return manage.exec(false)
}

func (manage *Manager) MultiInsert(data []map[string]interface{}) (*Manager, error) {
	manage.Generator.MultiInsert(data)
	return manage.exec(false)
}

func (manage *Manager) Update(data map[string]interface{}) (*Manager, error) {
	manage.Generator.Update(data)
	return manage.exec(false)
}

func (manage *Manager) Delete() (*Manager, error) {
	manage.Generator.Delete()
	return manage.exec(false)
}

func (manage *Manager) Get(args ...string) (*Manager, error) {
	manage.Generator.Get(args...)
	return manage.query()
}

func (manage *Manager) Value(field string) (interface{}, error) {
	manage.Generator.Value(field)
	res, err := manage.query()
	if err != nil {
		return nil, errors.New(err.Error() + ", sql: " + manage.Generator.ShowSql)
	}
	if len(res.Set) == 0 {
		return nil, nil
	}
	if val, ok := res.Set[0][field]; ok {
		return val, nil
	}
	return nil, errors.New("empty this field")
}

func (manage *Manager) First() (map[string]interface{}, error) {
	manage.Generator.First()
	res, err := manage.query()
	if err != nil {
		return nil, errors.New(err.Error() + ", sql: " + manage.Generator.ShowSql)
	}
	if len(res.Set) == 0 {
		return nil, nil
	}
	return res.Set[0], nil
}

func (manage *Manager) Model(model interface{}) (*Manager, error) {
	manage.Generator.First()
	res, err := manage.query()
	if err != nil {
		return manage, err
	}
	if len(res.Set) == 0 {
		return manage, nil
	}

	res.MapResult(&model)

	return manage, nil
}

func (manage *Manager) Models(model interface{}) (*Manager, error) {
	manage.Generator.Get()
	res, err := manage.query()
	if err != nil {
		return manage, err
	}
	if len(res.Set) == 0 {
		return manage, nil
	}

	res.MapResults(&model)

	return manage, nil
}

func (manage *Manager) PluckArray(field string) ([]interface{}, error) {
	manage.Generator.PluckArray(field)
	res, err := manage.query()
	if err != nil {
		return nil, errors.New(err.Error() + ", sql: " + manage.Generator.ShowSql)
	}
	if len(res.Set) == 0 {
		return []interface{}{}, nil
	}
	var data []interface{}
	for _, item := range res.Set {
		val, ok := item[field]
		if !ok {
			return nil, errors.New("no this field : " + field)
		}
		data = append(data, val)
	}
	return data, nil
}

func (manage *Manager) PluckMap(field, value string) (map[string]interface{}, error) {
	manage.Generator.PluckMap(field, value)
	res, err := manage.query()
	if err != nil {
		return nil, errors.New(err.Error() + ", sql: " + manage.Generator.ShowSql)
	}
	if len(res.Set) == 0 {
		return map[string]interface{}{}, nil
	}
	var data = make(map[string]interface{})
	for _, item := range res.Set {
		val, ok := item[value]
		if !ok {
			return nil, errors.New("no this field : " + value)
		}
		data[fmt.Sprintf("%v", item[field])] = val
	}
	return data, nil
}

// Todo Chunk
func (manage *Manager) Chunk(num int, callback func()) {
	manage.Generator.Chunk(num, callback)
}

func (manage *Manager) Count() (*Manager, error) {
	manage.Generator.Count()
	res, err := manage.query()
	if err != nil {
		return nil, errors.New(err.Error() + ", sql: " + manage.Generator.ShowSql)
	}
	if len(res.Set) == 0 {
		return nil, nil
	}
	err = manage.singleFuncField(res.Result, "count")
	if err != nil {
		return nil, err
	}
	return manage, nil
}

func (manage *Manager) Max(field string) (*Manager, error) {
	manage.Generator.Max(field)
	res, err := manage.query()
	if err != nil {
		return nil, errors.New(err.Error() + ", sql: " + manage.Generator.ShowSql)
	}
	if len(res.Set) == 0 {
		return nil, nil
	}
	err = manage.singleFuncField(res.Result, "max")
	if err != nil {
		return nil, err
	}
	return manage, nil
}

func (manage *Manager) Sum(field string) (*Manager, error) {
	manage.Generator.Sum(field)
	res, err := manage.query()
	if err != nil {
		return nil, errors.New(err.Error() + ", sql: " + manage.Generator.ShowSql)
	}
	if len(res.Set) == 0 {
		return nil, nil
	}
	err = manage.singleFuncField(res.Result, "sum")
	if err != nil {
		return nil, err
	}
	return manage, nil
}

func (manage *Manager) singleFuncField(res *result.Result, field string) (error) {
	if val, ok := res.Set[0][field].(string); ok {
		v, _ := strconv.Atoi(val)
		res.FuncResult = v
		return nil
	}

	val := fmt.Sprintf("%v", res.Set[0][field])
	v, _ := strconv.Atoi(val)
	res.FuncResult = v
	return nil
}

// Todo Exists
func (manage *Manager) Exists() {
	manage.Generator.Exists()
}

// Todo DoesntExists
func (manage *Manager) DoesntExists() {
	manage.Generator.DoesntExists()
}

func (manage *Manager) query() (*Manager, error) {
	var rows *sql.Rows
	var err error
	if manage.tx != nil {
		rows, err = manage.tx.Query(manage.Generator.ExeSql, manage.Generator.ExeParam...)
	} else {
		client, err := pool.GetClient()
		if err != nil {
			return manage, err
		}
		stmt, err := client.Prepare(manage.Generator.ExeSql)
		if err != nil {
			return manage, err
		}
		defer stmt.Close()
		rows, err = stmt.Query(manage.Generator.ExeParam...)
		defer pool.CloseClient(client)
	}

	if err != nil {
		return manage, errors.New(err.Error() + ", sql: " + manage.Generator.ShowSql)
	}

	manage.Result = new(result.Result)
	manage.Result.Rows = rows
	if rows != nil {
		manage.Result.MakeResult()
	}

	return manage, nil
}

func (manage *Manager) exec(InsertId bool) (*Manager, error) {
	var ret sql.Result
	var err error
	if manage.tx != nil {
		ret, err = manage.tx.Exec(manage.Generator.ExeSql, manage.Generator.ExeParam...)
	} else {
		client, err := pool.GetClient()
		if err != nil {
			return manage, err
		}
		stmt, err := client.Prepare(manage.Generator.ExeSql)
		if err != nil {
			return manage, err
		}
		defer stmt.Close()
		ret, err = stmt.Exec(manage.Generator.ExeParam...)
		defer pool.CloseClient(client)
	}

	if err != nil {
		return manage, errors.New(err.Error() + ", sql: " + manage.Generator.ShowSql)
	}

	manage.Result = new(result.Result)

	if InsertId {
		manage.Result.LastInsertId, err = ret.LastInsertId()
	} else {
		manage.Result.RowsAffected, err = ret.RowsAffected()
	}
	if err != nil {
		return manage, err
	}

	return manage, nil
}

func (manage *Manager) LastInsertIdToSql(data map[string]interface{}) string {
	manage.Generator.Insert(data)
	return manage.Generator.ShowSql
}

func (manage *Manager) InsertToSql(data map[string]interface{}) string {
	manage.Generator.Insert(data)
	return manage.Generator.ShowSql
}

func (manage *Manager) MultiInsertToSql(data []map[string]interface{}) string {
	manage.Generator.MultiInsert(data)
	return manage.Generator.ShowSql
}

func (manage *Manager) UpdateToSql(data map[string]interface{}) string {
	manage.Generator.Update(data)
	return manage.Generator.ShowSql
}

func (manage *Manager) DeleteToSql() string {
	manage.Generator.Delete()
	return manage.Generator.ShowSql
}
func (manage *Manager) GetToSql(args ...string) string {
	manage.Generator.Get(args...)
	return manage.Generator.ShowSql
}
func (manage *Manager) ValueToSql(field string) string {
	manage.Generator.Value(field)
	return manage.Generator.ShowSql
}
func (manage *Manager) FirstToSql(args ...string) string {
	manage.Generator.First()
	return manage.Generator.ShowSql
}
func (manage *Manager) PluckArrayToSql(field string) string {
	manage.Generator.PluckArray(field)
	return manage.Generator.ShowSql
}
func (manage *Manager) PluckMapToSql(field, value string) string {
	manage.Generator.PluckMap(field, value)
	return manage.Generator.ShowSql
}
func (manage *Manager) CountToSql() string {
	manage.Generator.Count()
	return manage.Generator.ShowSql
}
func (manage *Manager) MaxToSql(field string) string {
	manage.Generator.Max(field)
	return manage.Generator.ShowSql
}
func (manage *Manager) SumToSql(field string) string {
	manage.Generator.Sum(field)
	return manage.Generator.ShowSql
}
func (manage *Manager) ChunkToSql(num int) string {
	manage.Generator.Limit(num)
	manage.Generator.Get()
	return manage.Generator.ShowSql
}
