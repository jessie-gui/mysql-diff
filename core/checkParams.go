package core

import (
	"log"
	"strings"
)

type ParamsOption func(p *Params)

// Params /**
type Params struct {
	source string
	target string
	db     string
}

func Source(source string) ParamsOption {
	return func(p *Params) {
		p.source = source
	}
}

func Target(target string) ParamsOption {
	return func(p *Params) {
		p.target = target
	}
}

func Db(db string) ParamsOption {
	return func(p *Params) {
		p.db = db
	}
}

func NewParams(opts ...ParamsOption) *Params {
	params := &Params{}
	for _, opt := range opts {
		opt(params)
	}

	return params
}

func (p *Params) CheckParams() map[string]*Mysql {
	var sourceUser, sourcePwd, sourceHost, sourcePort, sourceDb, targetUser, targetPwd, targetHost, targetPort, targetDb string

	sourceArr := strings.Split(p.source, "@")
	if sourceArr == nil || len(sourceArr) != 2 {
		log.Fatal("源服务器格式错误！格式: <user>:<password>@<host>:<port>")
	}

	for k, v := range sourceArr {
		vArr := strings.Split(v, ":")
		if len(vArr) != 2 {
			log.Fatal("源服务器格式错误！格式: <user>:<password>@<host>:<port>")
		}

		if k == 0 {
			sourceUser = vArr[0]
			sourcePwd = vArr[1]
		} else {
			sourceHost = vArr[0]
			sourcePort = vArr[1]
		}
	}

	targetArr := strings.Split(p.target, "@")
	if targetArr == nil || len(sourceArr) != 2 {
		log.Fatal("目标服务器格式错误！格式: <user>:<password>@<host>:<port>")
	}

	for k, v := range targetArr {
		vArr := strings.Split(v, ":")
		if len(vArr) != 2 {
			log.Fatal("目标服务器格式错误！格式: <user>:<password>@<host>:<port>")
		}

		if k == 0 {
			targetUser = vArr[0]
			targetPwd = vArr[1]
		} else {
			targetHost = vArr[0]
			targetPort = vArr[1]
		}
	}

	dbArr := strings.Split(p.db, ":")
	if len(dbArr) != 2 {
		log.Fatal("指定数据库格式错误！格式: <source_db>:<target_db>")
	}

	sourceDb = dbArr[0]
	targetDb = dbArr[1]

	sourceMysql := NewMysql(User(sourceUser), Pwd(sourcePwd), Host(sourceHost), Port(sourcePort), DbName(sourceDb))
	targetMysql := NewMysql(User(targetUser), Pwd(targetPwd), Host(targetHost), Port(targetPort), DbName(targetDb))

	return map[string]*Mysql{
		"sourceMysql": sourceMysql,
		"targetMysql": targetMysql,
	}
}
