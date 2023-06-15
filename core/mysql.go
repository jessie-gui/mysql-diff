package core

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
)

type MysqlOption func(m *Mysql)

// Mysql /**
type Mysql struct {
	user string
	pwd  string
	host string
	port string
	db   string
}

func User(user string) MysqlOption {
	return func(m *Mysql) {
		m.user = user
	}
}

func Pwd(pwd string) MysqlOption {
	return func(m *Mysql) {
		m.pwd = pwd
	}
}

func Host(host string) MysqlOption {
	return func(m *Mysql) {
		m.host = host
	}
}

func Port(port string) MysqlOption {
	return func(m *Mysql) {
		m.port = port
	}
}

func DbName(db string) MysqlOption {
	return func(m *Mysql) {
		m.db = db
	}
}

func NewMysql(opts ...MysqlOption) *Mysql {
	m := &Mysql{}
	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Mysql) Conn() *gorm.DB {
	dsnBase := fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema?timeout=10s&charset=utf8mb4&parseTime=True&loc=Local", m.user, m.pwd, m.host, m.port)
	dbBase, err := gorm.Open(mysql.New(mysql.Config{DSN: dsnBase}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal("connect DbBase failed", err)
	}

	return dbBase
}

func (m *Mysql) GetUser() string {
	return m.user
}

func (m *Mysql) GetPwd() string {
	return m.pwd
}

func (m *Mysql) GetHost() string {
	return m.host
}

func (m *Mysql) GetPort() string {
	return m.port
}

func (m *Mysql) GetDb() string {
	return m.db
}
