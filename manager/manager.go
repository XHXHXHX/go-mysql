package mysqlManager

import (
	"database/sql"
	"errors"
	"fmt"
	mysqlPool "go-mysql/pool"
	mysqlResult "go-mysql/result"
	generators "go-mysql/sql-generators"
	"strconv"
)

const DEFAULT_CONFIG = "config/mysql.json"

/*
 * Manager
 */
type Manager struct {
	generator *generators.SqlGenerator
	client    *sql.DB
	tx        *sql.Tx
	model     interface{}
	models    []interface{}
}

func InitConfig(config string) error {
	if len(config) == 0 {
		config = DEFAULT_CONFIG
	}
	err := mysqlPool.InitPool(config)
	if err != nil {
		return err
	}
	return nil
}

func DB() *Manager {
	return &Manager{}
}

func (manage *Manager) BindModel(ptr interface{}) *Manager {
	manage.model = ptr
	return manage
}

func (manage *Manager) BindModels(ptr []interface{}) *Manager {
	manage.models = ptr
	return manage
}

func DbBegin() *Manager {
	client, err := mysqlPool.GetClient()
	if err != nil {
		panic(err)
	}
	tx, err := client.Begin()
	if err != nil {
		panic(err)
	}
	return &Manager{
		generator: generators.NewGenerator(),
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
	err := mysqlPool.CloseClient(manage.client)
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
	err = mysqlPool.CloseClient(manage.client)
	if err != nil {
		return err
	}
	manage.tx = nil
	manage.client = nil
	return nil
}

func (manage *Manager) LastInsertId(data map[string]interface{}) (int, error) {
	manage.generator.Insert(data)
	return manage.exec(true)
}

func (manage *Manager) Insert(data map[string]interface{}) (int, error) {
	manage.generator.Insert(data)
	return manage.exec(false)
}

func (manage *Manager) MultiInsert(data []map[string]interface{}) (int, error) {
	manage.generator.MultiInsert(data)
	return manage.exec(false)
}

func (manage *Manager) Update(data map[string]interface{}) (int, error) {
	manage.generator.Update(data)
	return manage.exec(false)
}

func (manage *Manager) Delete() (int, error) {
	manage.generator.Delete()
	return manage.exec(false)
}

func (manage *Manager) Get(args ...string) (*mysqlResult.Result, error) {
	manage.generator.Get(args...)
	return manage.query()
}

func (manage *Manager) Value(field string) (interface{}, error) {
	manage.generator.Value(field)
	res, err := manage.query()
	if err != nil || len(res.Set) == 0 {
		return nil, err
	}
	if val, ok := res.Set[0][field]; ok {
		return val, nil
	}
	return nil, errors.New("empty this field")
}

func (manage *Manager) First() (map[string]interface{}, error) {
	manage.generator.First()
	res, err := manage.query()
	if err != nil || len(res.Set) == 0 {
		return nil, err
	}
	return res.Set[0], nil
}

func (manage *Manager) FirstModel(model interface{}) error {
	manage.generator.First()
	res, err := manage.query()
	if err != nil || len(res.Set) == 0 {
		return err
	}

	res.MapResult(&model)

	return nil
}

func (manage *Manager) PluckArray(field string) ([]interface{}, error) {
	manage.generator.PluckArray(field)
	res, err := manage.query()
	if err != nil || len(res.Set) == 0 {
		return nil, err
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
	manage.generator.PluckMap(field, value)
	res, err := manage.query()
	if err != nil || len(res.Set) == 0 {
		return nil, err
	}
	var data = make(map[string]interface{})
	for _, item := range res.Set {
		val, ok := item[field]
		if !ok {
			return nil, errors.New("no this field : " + field)
		}
		data[fmt.Sprintf("%v", item[field])] = val
	}
	return data, nil
}

// Todo Chunk
func (manage *Manager) Chunk(num int, callback func()) {
	manage.generator.Chunk(num, callback)
}

func (manage *Manager) Count() (int, error) {
	manage.generator.Count()
	res, err := manage.query()
	if err != nil || len(res.Set) == 0 {
		return 0, err
	}
	return manage.singleFuncField(res, "count")
}

func (manage *Manager) Max(field string) (int, error) {
	manage.generator.Max(field)
	res, err := manage.query()
	if err != nil || len(res.Set) == 0 {
		return 0, err
	}
	return manage.singleFuncField(res, "max")
}

func (manage *Manager) Sum(field string) (int, error) {
	manage.generator.Sum(field)
	res, err := manage.query()
	if err != nil || len(res.Set) == 0 {
		return 0, err
	}
	return manage.singleFuncField(res, "max")
}

func (manage *Manager) singleFuncField(res *mysqlResult.Result, field string) (int, error) {
	if val, ok := res.Set[0]["count"].(string); ok {
		v, _ := strconv.Atoi(val)
		return v, nil
	}

	val := fmt.Sprintf("%v", res.Set[0]["count"])
	v, _ := strconv.Atoi(val)
	return v, nil
}

// Todo Exists
func (manage *Manager) Exists() {
	manage.generator.Exists()
}

// Todo DoesntExists
func (manage *Manager) DoesntExists() {
	manage.generator.DoesntExists()
}

func (manage *Manager) query() (*mysqlResult.Result, error) {
	var rows *sql.Rows
	var err error
	if manage.tx != nil {
		rows, err = manage.tx.Query(manage.generator.ExeSql, manage.generator.ExeParam...)
	} else {
		client, err := mysqlPool.GetClient()
		if err != nil {
			return nil, err
		}
		stmt, err := client.Prepare(manage.generator.ExeSql)
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		rows, err = stmt.Query(manage.generator.ExeParam...)
		defer mysqlPool.CloseClient(client)
	}

	if err != nil {
		return nil, err
	}

	m := new(mysqlResult.Result)
	m.Rows = rows
	m.MakeResult()

	return m, nil
}

func (manage *Manager) exec(InsertId bool) (int, error) {
	var ret sql.Result
	var err error
	if manage.tx != nil {
		ret, err = manage.tx.Exec(manage.generator.ExeSql, manage.generator.ExeParam...)
	} else {
		client, err := mysqlPool.GetClient()
		if err != nil {
			return 0, err
		}
		stmt, err := client.Prepare(manage.generator.ExeSql)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()
		ret, err = stmt.Exec(manage.generator.ExeParam...)
		defer mysqlPool.CloseClient(client)
	}

	if err != nil {
		return 0, err
	}

	var num int64
	if InsertId {
		num, err = ret.LastInsertId()
	} else {
		num, err = ret.RowsAffected()
	}
	if err != nil {
		return 0, err
	}

	return int(num), nil
}

func (manage *Manager) LastInsertIdToSql(data map[string]interface{}) string {
	manage.generator.Insert(data)
	return manage.generator.ShowSql
}

func (manage *Manager) InsertToSql(data map[string]interface{}) string {
	manage.generator.Insert(data)
	return manage.generator.ShowSql
}

func (manage *Manager) MultiInsertToSql(data []map[string]interface{}) string {
	manage.generator.MultiInsert(data)
	return manage.generator.ShowSql
}

func (manage *Manager) UpdateToSql(data map[string]interface{}) string {
	manage.generator.Update(data)
	return manage.generator.ShowSql
}

func (manage *Manager) DeleteToSql() string {
	manage.generator.Delete()
	return manage.generator.ShowSql
}
func (manage *Manager) GetToSql(args ...string) string {
	manage.generator.Get(args...)
	return manage.generator.ShowSql
}
func (manage *Manager) ValueToSql(field string) string {
	manage.generator.Value(field)
	return manage.generator.ShowSql
}
func (manage *Manager) FirstToSql(args ...string) string {
	manage.generator.First()
	return manage.generator.ShowSql
}
func (manage *Manager) PluckArrayToSql(field string) string {
	manage.generator.PluckArray(field)
	return manage.generator.ShowSql
}
func (manage *Manager) PluckMapToSql(field, value string) string {
	manage.generator.PluckMap(field, value)
	return manage.generator.ShowSql
}
func (manage *Manager) CountToSql() string {
	manage.generator.Count()
	return manage.generator.ShowSql
}
func (manage *Manager) MaxToSql(field string) string {
	manage.generator.Max(field)
	return manage.generator.ShowSql
}
func (manage *Manager) SumToSql(field string) string {
	manage.generator.Sum(field)
	return manage.generator.ShowSql
}
func (manage *Manager) ChunkToSql(num int) string {
	manage.generator.Limit(num)
	manage.generator.Get()
	return manage.generator.ShowSql
}
