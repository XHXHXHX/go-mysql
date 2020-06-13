package pool

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"errors"
	"sync"
	"time"
)

type Client struct {
	MysqlClient *sql.DB			// Mysql 连接
	expire time.Time			// 超时时间 	0：不超时
}

type MysqlConfig struct {
	InitCap int
	MaxCap int
	DBName string
	Dsn string
	ClientTimeOut int
	CheckClientAliveInterval int
	KeepClientTime int
	Host string
	Port int
	Username string
	Password string
}

type Pool struct {
	wait sync.RWMutex
	useMap map[*sql.DB] *Client		// 使用中链接
	Config *MysqlConfig				// 配置信息
	Clients chan *Client			// 空闲链接池
	ClientNum int					// 已生成链接数
}

var myPool *Pool
var waitGroup sync.WaitGroup


/*
 * 初始化连接池
 */
func (this *Pool) InitClient() {
	waitGroup.Add(this.Config.InitCap)
	for i := 0; i < this.Config.InitCap; i++ {
		go func() {
			client, err := this.CreateClient()
			if err != nil {
				panic(err)
			}
			this.Clients <- client
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
}

/*
 * 连接Mysql
 */
func (this *Pool) clientMysql() (*sql.DB, error) {
	this.wait.RLock()
	defer this.wait.RUnlock()

	if this.ClientNum == this.Config.MaxCap {
		return nil, errors.New("Pool clientMysql error, client num to reach the max cap")
	}

	db, err := sql.Open("mysql", this.Config.Dsn)
	if err != nil {
		return nil, errors.New("Pool clientMysql error, sql.Open : " + err.Error())
	}

	return db, nil
}

func (this *Pool) CreateClient() (*Client, error) {
	db, err := this.clientMysql()
	if err != nil {
		return nil, err
	}

	this.ClientNum++
	time_unit, _ := time.ParseDuration("1h")

	return &Client{
		MysqlClient:	db,
		expire:			time.Now().Add(time.Duration(this.Config.KeepClientTime) * time_unit),
	}, nil
}

/*
 * 定期检查失效链接
 */
func (this *Pool) checkInvalidClient() {
	this.wait.Lock()
	defer this.wait.Unlock()

	tmp_client := make(chan *Client, len(this.Clients))
	for client := range this.Clients {
		if this.Len() >= this.Config.MaxCap {
			_ = this.Close(client)
			break
		}

		if time.Now().After(client.expire) {
			_ = this.Close(client)
			break
		}

		if err := client.MysqlClient.Ping(); err != nil {
			_ = this.Close(client)
			break
		}

		tmp_client <- client
	}

	for client := range tmp_client {
		this.Clients <- client
	}
}

func (this *Pool) Len() int {
	return len(this.Clients)
}

func (this *Pool) Close(client *Client) error {
	err := client.MysqlClient.Close()
	this.ClientNum--
	if err != nil {
		return errors.New("Pool Close errors, " + err.Error())
	}

	return nil
}